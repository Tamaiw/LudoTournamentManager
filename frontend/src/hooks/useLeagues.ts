import { useQuery } from '@tanstack/react-query';
import { api } from '../services/api';

export function useLeagues() {
  return useQuery({
    queryKey: ['leagues'],
    queryFn: () => api.listLeagues(),
  });
}