package outbound

import "ludo-tournament/core/domain/models"

type MatchRepository interface {
	Create(match *models.Match) error
	GetByID(id string) (*models.Match, error)
	Update(match *models.Match) error
	ListByTournament(tournamentID string) ([]models.Match, error)
	ListByLeague(leagueID string) ([]models.Match, error)
	ListByRound(tournamentID string, round int) ([]models.Match, error)
	GetCompletedCountInRound(tournamentID string, round int) (int, error)
}
