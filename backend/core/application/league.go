package application

import (
	"context"
	"errors"
	"sort"

	"ludo-tournament/core/domain"
	"ludo-tournament/core/domain/models"
	"ludo-tournament/core/ports/inbound"
	"ludo-tournament/core/ports/outbound"
)

// Pair represents a pair of players who have played together
type Pair = models.Pair

// LeagueMatchResult represents a player's result in a league match
type LeagueMatchResult = models.LeagueMatchResult

// LeagueService implements the inbound.LeagueService interface
type LeagueService struct {
	leagueRepo outbound.LeagueRepository
	matchRepo  outbound.MatchRepository
}

// NewLeagueService creates a new LeagueService instance
func NewLeagueService(leagueRepo outbound.LeagueRepository, matchRepo outbound.MatchRepository) *LeagueService {
	return &LeagueService{
		leagueRepo: leagueRepo,
		matchRepo:  matchRepo,
	}
}

// GenerateFairPairings creates pairings that minimize repeat matches using a greedy algorithm
// Algorithm:
// 1. Build conflict map from priorPairings (who has played with whom)
// 2. For each table, select 4 players minimizing conflicts with each other
// 3. Use greedy approach: find least-matched player first, then find best opponents
// 4. Return pairs for seating (A-B, C-D typical Ludo pairing)
func GenerateFairPairings(players []string, priorPairings []Pair, numTables int) ([]inbound.TablePairing, error) {
	if len(players) == 0 {
		return nil, errors.New("no players provided")
	}

	playersPerTable := 4
	totalSlots := numTables * playersPerTable
	if len(players) < 4 || len(players) != totalSlots {
		return nil, errors.New("insufficient players: need exactly 4 players per table")
	}

	// Build conflict map: player -> set of players they've played with
	conflictMap := make(map[string]map[string]bool)
	for _, p := range players {
		conflictMap[p] = make(map[string]bool)
	}
	for _, pair := range priorPairings {
		conflictMap[pair.Player1][pair.Player2] = true
		conflictMap[pair.Player2][pair.Player1] = true
	}

	// Count how many times each player has been paired
	pairCount := make(map[string]int)
	for _, p := range players {
		pairCount[p] = 0
	}
	for _, pair := range priorPairings {
		pairCount[pair.Player1]++
		pairCount[pair.Player2]++
	}

	// Track assigned players
	assigned := make(map[string]bool)
	for _, p := range players {
		assigned[p] = false
	}

	// Result pairings
	var pairings []inbound.TablePairing

	// Greedy assignment: for each table, pick 4 players with minimum conflicts
	for tableNum := 1; tableNum <= numTables; tableNum++ {
		// Find player with lowest pair count among unassigned
		var selectedPlayer string
		minCount := int(^uint(0) >> 1) // Max int

		for _, p := range players {
			if !assigned[p] && pairCount[p] < minCount {
				minCount = pairCount[p]
				selectedPlayer = p
			}
		}

		if selectedPlayer == "" {
			break
		}

		// Find 3 other players that minimize conflicts with selected player
		candidates := getUnassignedPlayers(players, assigned)
		// Remove selectedPlayer from candidates
		var opponents []string
		for _, c := range candidates {
			if c != selectedPlayer {
				opponents = append(opponents, c)
			}
		}

		if len(opponents) < 3 {
			return nil, errors.New("not enough unassigned players")
		}

		// Find best combination of 3 opponents
		var bestOpponents []string
		bestConflictScore := int(^uint(0) >> 1)

		for i := 0; i < len(opponents); i++ {
			for j := i + 1; j < len(opponents); j++ {
				for k := j + 1; k < len(opponents); k++ {
					potentialOpponents := []string{opponents[i], opponents[j], opponents[k]}
					score := countConflicts(selectedPlayer, potentialOpponents, conflictMap)
					if score < bestConflictScore {
						bestConflictScore = score
						bestOpponents = potentialOpponents
					}
				}
			}
		}

		if len(bestOpponents) != 3 {
			return nil, errors.New("could not find suitable opponents")
		}

		// Build table players list
		tablePlayers := []string{selectedPlayer}
		tablePlayers = append(tablePlayers, bestOpponents...)

		// Mark these players as assigned
		assigned[selectedPlayer] = true
		for _, op := range bestOpponents {
			assigned[op] = true
		}

		// Update pair counts for future iterations
		for _, p1 := range tablePlayers {
			for _, p2 := range tablePlayers {
				if p1 != p2 {
					pairCount[p1]++
				}
			}
		}

		// Create pairing for this table
		pairing := inbound.TablePairing{
			TableNumber: tableNum,
			PlayerIDs:   tablePlayers,
		}
		pairings = append(pairings, pairing)
	}

	// Verify all players assigned
	for _, p := range players {
		if !assigned[p] {
			return nil, errors.New("not all players were assigned")
		}
	}

	return pairings, nil
}

// getUnassignedPlayers returns list of unassigned players
func getUnassignedPlayers(allPlayers []string, assigned map[string]bool) []string {
	var result []string
	for _, p := range allPlayers {
		if !assigned[p] {
			result = append(result, p)
		}
	}
	return result
}

// countConflicts counts how many times player has played with any of the opponents
// and also counts conflicts between opponents themselves
func countConflicts(player string, opponents []string, conflictMap map[string]map[string]bool) int {
	count := 0
	// Count conflicts between player and opponents
	for _, op := range opponents {
		if conflictMap[player][op] {
			count++
		}
	}
	// Count conflicts between opponents themselves
	for i := 0; i < len(opponents); i++ {
		for j := i + 1; j < len(opponents); j++ {
			if conflictMap[opponents[i]][opponents[j]] {
				count++
			}
		}
	}
	return count
}

// hasPlayedTogether checks if two players have previously played together
func hasPlayedTogether(p1, p2 string, priorPairings []Pair) bool {
	for _, pair := range priorPairings {
		if (pair.Player1 == p1 && pair.Player2 == p2) ||
			(pair.Player1 == p2 && pair.Player2 == p1) {
			return true
		}
	}
	return false
}

// CalculateLeagueStandings calculates standings based on match results and scoring rules
func CalculateLeagueStandings(matchResults []LeagueMatchResult, scoringRules []models.ScoringRule) []inbound.PlayerStanding {
	if len(matchResults) == 0 {
		return []inbound.PlayerStanding{}
	}

	// Build scoring lookup
	pointsLookup := make(map[int]float64)
	for _, rule := range scoringRules {
		pointsLookup[rule.Placement] = rule.Points
	}

	// Aggregate player stats
	type playerStats struct {
		totalPoints float64
		wins        int
		gamesPlayed int
	}

	playerStatsMap := make(map[string]*playerStats)
	playerOrder := []string{}

	for _, result := range matchResults {
		stats, exists := playerStatsMap[result.PlayerID]
		if !exists {
			stats = &playerStats{}
			playerStatsMap[result.PlayerID] = stats
			playerOrder = append(playerOrder, result.PlayerID)
		}

		stats.gamesPlayed++
		stats.totalPoints += pointsLookup[result.Placement]

		if result.Placement == 1 {
			stats.wins++
		}
	}

	// Build standings
	standings := make([]inbound.PlayerStanding, 0, len(playerStatsMap))
	for _, playerID := range playerOrder {
		stats := playerStatsMap[playerID]
		standings = append(standings, inbound.PlayerStanding{
			PlayerID:    playerID,
			DisplayName: playerID, // Would be fetched from player repository in real impl
			GamesPlayed: stats.gamesPlayed,
			TotalPoints: stats.totalPoints,
			Wins:        stats.wins,
		})
	}

	// Sort by total points (desc), then by wins (desc)
	sort.Slice(standings, func(i, j int) bool {
		if standings[i].TotalPoints != standings[j].TotalPoints {
			return standings[i].TotalPoints > standings[j].TotalPoints
		}
		return standings[i].Wins > standings[j].Wins
	})

	// Assign ranks
	for i := range standings {
		standings[i].Rank = i + 1
	}

	return standings
}

// CreateLeague creates a new league
func (s *LeagueService) CreateLeague(ctx context.Context, name string, organizerID string, settings models.LeagueSettings) (*models.League, error) {
	if name == "" {
		return nil, domain.ErrInvalidInput
	}

	league := &models.League{
		Name:        name,
		OrganizerID: organizerID,
		Status:      models.LeagueStatusDraft,
		Settings:    settings,
	}

	if err := s.leagueRepo.Create(ctx, league); err != nil {
		return nil, err
	}

	return league, nil
}

// GetLeague retrieves a league by ID
func (s *LeagueService) GetLeague(ctx context.Context, id string) (*models.League, error) {
	league, err := s.leagueRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return league, nil
}

// UpdateLeague updates league settings
func (s *LeagueService) UpdateLeague(ctx context.Context, id string, settings models.LeagueSettings) error {
	league, err := s.leagueRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	league.Settings = settings
	return s.leagueRepo.Update(ctx, league)
}

// DeleteLeague soft-deletes a league
func (s *LeagueService) DeleteLeague(ctx context.Context, id string) error {
	return s.leagueRepo.SoftDelete(ctx, id)
}

// GenerateSchedule stores play dates for the league
func (s *LeagueService) GenerateSchedule(ctx context.Context, leagueID string, playDates []string) error {
	league, err := s.leagueRepo.GetByID(ctx, leagueID)
	if err != nil {
		return err
	}

	// Store play dates in league settings or a separate structure
	// For now, just validate the league exists and mark it as live
	league.Status = models.LeagueStatusLive

	// In a full implementation, we would store play dates somewhere
	// This is a simplified version

	return s.leagueRepo.Update(ctx, league)
}

// GeneratePairings generates pairings for a specific play date
func (s *LeagueService) GeneratePairings(ctx context.Context, leagueID string, playDate string) ([]inbound.TablePairing, error) {
	league, err := s.leagueRepo.GetByID(ctx, leagueID)
	if err != nil {
		return nil, err
	}

	// Get all matches for this league to build prior pairings
	matches, err := s.matchRepo.ListByLeague(ctx, leagueID)
	if err != nil {
		return nil, err
	}

	// Build prior pairings from completed matches
	var priorPairings []Pair
	for _, match := range matches {
		if match.Status == models.MatchStatusCompleted {
			// In a real implementation, we would get player assignments from match
			// For now, we assume the match has player assignments stored elsewhere
		}
	}

	// Get league players (in real impl, would come from league player registration)
	// For this implementation, we'll use a placeholder
	players := []string{"player1", "player2", "player3", "player4",
		"player5", "player6", "player7", "player8"}

	numTables := league.Settings.TablesCount
	if numTables == 0 {
		numTables = 2 // Default
	}

	return GenerateFairPairings(players, priorPairings, numTables)
}

// ReportLeagueMatch reports results for a league match
func (s *LeagueService) ReportLeagueMatch(ctx context.Context, matchID string, results []inbound.LeagueMatchResultInput, reportedBy string) error {
	match, err := s.matchRepo.GetByID(ctx, matchID)
	if err != nil {
		return err
	}

	if match.LeagueID == nil {
		return errors.New("match is not a league match")
	}

	// Update match status and results
	match.Status = models.MatchStatusCompleted

	return s.matchRepo.Update(ctx, match)
}

// GetStandings returns current league standings
func (s *LeagueService) GetStandings(ctx context.Context, leagueID string) ([]inbound.PlayerStanding, error) {
	league, err := s.leagueRepo.GetByID(ctx, leagueID)
	if err != nil {
		return nil, err
	}

	matches, err := s.matchRepo.ListByLeague(ctx, leagueID)
	if err != nil {
		return nil, err
	}

	// Build match results from completed matches
	var matchResults []LeagueMatchResult
	for _, match := range matches {
		if match.Status == models.MatchStatusCompleted && match.PlacementPoints != nil {
			// In real impl, would iterate over match assignments
			// Here we assume placement points are stored in match
		}
	}

	return CalculateLeagueStandings(matchResults, league.Settings.ScoringRules), nil
}

// AddTiebreaker handles tiebreaker logic for players with equal standings
func (s *LeagueService) AddTiebreaker(ctx context.Context, leagueID string, playerIDs []string) error {
	// In a full implementation, this would store tiebreaker rules for the league
	// For now, this is a placeholder
	return nil
}

// Ensure LeagueService implements inbound.LeagueService
var _ inbound.LeagueService = (*LeagueService)(nil)
