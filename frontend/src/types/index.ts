export type Role = 'admin' | 'member' | 'guest';

export interface User {
  id: string;
  email: string;
  role: Role;
  lastActive?: string;
  createdAt: string;
}

export interface TournamentSettings {
  tablesCount: number;
  advancement?: AdvancementConfig[];
  defaultReporter?: string;
}

export interface Tournament {
  id: string;
  name: string;
  type: 'knockout';
  organizerId: string;
  status: 'draft' | 'live' | 'completed';
  settings: TournamentSettings;
  createdAt: string;
}

export interface LeagueSettings {
  scoringRules: ScoringRule[];
  gamesPerPlayer: number;
  tablesCount: number;
}

export interface ScoringRule {
  placement: number;
  points: number;
}

export interface League {
  id: string;
  name: string;
  organizerId: string;
  status: 'draft' | 'live' | 'completed';
  settings: LeagueSettings;
  createdAt: string;
}

export interface Match {
  id: string;
  tournamentId?: string;
  leagueId?: string;
  round: number;
  tableNumber: number;
  status: 'pending' | 'completed';
  placementPoints?: number[];
  completedAt?: string;
}

export interface GamePairing {
  gameId: string;
  round: number;
  tableNumber: number;
  playerIds: string[];
  seatColors: ('yellow' | 'green' | 'blue' | 'red')[];
  status: 'pending' | 'completed';
}

export interface PlayerStanding {
  playerId: string;
  displayName: string;
  gamesPlayed: number;
  totalPoints: number;
  wins: number;
  rank: number;
}