import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../services/api';

export function useLeague(id: string) {
  return useQuery({
    queryKey: ['league', id],
    queryFn: () => api.getLeague(id),
  });
}

export function useLeagueStandings(id: string) {
  return useQuery({
    queryKey: ['league', id, 'standings'],
    queryFn: () => api.getLeagueStandings(id),
  });
}

export function useGenerateLeaguePairings() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ leagueId, playDate }: { leagueId: string; playDate: string }) =>
      api.generateLeaguePairings(leagueId, playDate),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['league', variables.leagueId] });
    },
  });
}