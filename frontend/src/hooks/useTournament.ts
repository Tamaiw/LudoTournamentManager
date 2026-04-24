import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../services/api';

export function useTournament(id: string) {
  return useQuery({
    queryKey: ['tournament', id],
    queryFn: () => api.getTournament(id),
  });
}

export function useTournamentMatches(id: string) {
  return useQuery({
    queryKey: ['tournament', id, 'matches'],
    queryFn: () => api.getTournamentMatches(id),
  });
}

export function useTournamentPairings(id: string) {
  return useQuery({
    queryKey: ['tournament', id, 'pairings'],
    queryFn: () => api.getPairings(id),
  });
}

export function useReportMatch() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ tournamentId, matchId, results }: { tournamentId: string; matchId: string; results: any[] }) =>
      api.reportMatch(tournamentId, matchId, results),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['tournament', variables.tournamentId] });
    },
  });
}