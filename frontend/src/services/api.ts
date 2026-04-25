import { Role, User, League, PlayerStanding, Tournament, Match, GamePairing, MatchResult } from '../types';

const API_BASE = '';

let authToken: string | null = null;

export function setAuthToken(token: string | null) {
  authToken = token;
}

export function getAuthToken(): string | null {
  return authToken;
}

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const headers = new Headers(options?.headers);
  headers.set('Content-Type', 'application/json');

  if (authToken) {
    headers.set('Authorization', `Bearer ${authToken}`);
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
  } as RequestInit);

  if (!res.ok) {
    const error = await res.json().catch(() => ({}));
    throw new Error(error.error?.message || 'Request failed');
  }

  return res.json();
}

export const api = {
  login: (email: string, password: string) =>
    request<{ token: string }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),

  register: (email: string, password: string, inviteCode: string) =>
    request<{ token: string }>('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ email, password, invite_code: inviteCode }),
    }),

  logout: () => request<void>('/auth/logout', { method: 'POST' }),

  getMe: () => request<User>('/auth/me'),

  listTournaments: () => request<Tournament[]>('/tournaments'),
  getTournament: (id: string) => request<Tournament>(`/tournaments/${id}`),
  createTournament: (data: Partial<Tournament>) =>
    request<Tournament>('/tournaments', { method: 'POST', body: JSON.stringify(data) }),
  getTournamentMatches: (id: string) => request<Match[]>(`/tournaments/${id}/matches`),
  reportMatch: (tournamentId: string, matchId: string, results: MatchResult[]) =>
    request(`/tournaments/${tournamentId}/matches`, {
      method: 'POST',
      body: JSON.stringify({ match_id: matchId, results }),
    }),
  getPairings: (tournamentId: string) =>
    request<GamePairing[]>(`/tournaments/${tournamentId}/pairings`),

  listLeagues: () => request<League[]>('/leagues'),
  getLeague: (id: string) => request<League>(`/leagues/${id}`),
  getLeagueStandings: (id: string) => request<PlayerStanding[]>(`/leagues/${id}/standings`),
  generatePairings: (leagueId: string, playDate: string) =>
    request(`/leagues/${leagueId}/pairings/generate`, {
      method: 'POST',
      body: JSON.stringify({ play_date: playDate }),
    }),

  listUsers: () => request<{ users: User[] }>('/users'),
  updateUser: (id: string, data: { role: Role }) =>
    request<User>(`/users/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  deleteUser: (id: string) =>
    request<void>(`/users/${id}`, { method: 'DELETE' }),
  createLeague: (data: Partial<League>) =>
    request<League>('/leagues', { method: 'POST', body: JSON.stringify(data) }),
  generateLeaguePairings: (leagueId: string, playDate: string) =>
    request(`/leagues/${leagueId}/pairings/generate`, {
      method: 'POST',
      body: JSON.stringify({ play_date: playDate }),
    }),
};