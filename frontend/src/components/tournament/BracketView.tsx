import { useQuery } from '@tanstack/react-query';
import { useParams } from 'react-router-dom';
import { api } from '../../services/api';
import { TableAssignment } from './TableAssignment';
import { MatchCard } from './MatchCard';

interface BracketData {
  rounds: { [round: number]: GamePairing[] };
}

export function BracketView() {
  const { id } = useParams<{ id: string }>();

  const { data: pairings, isLoading } = useQuery({
    queryKey: ['pairings', id],
    queryFn: () => api.getPairings(id!),
    enabled: !!id,
  });

  if (isLoading) return <div>Loading...</div>;
  if (!pairings) return <div>No pairings found</div>;

  // Group pairings by round
  const byRound: { [key: number]: GamePairing[] } = {};
  pairings.forEach(p => {
    if (!byRound[p.round]) byRound[p.round] = [];
    byRound[p.round].push(p);
  });

  const rounds = Object.keys(byRound).map(Number).sort();

  return (
    <div className="space-y-8">
      {rounds.map(round => (
        <TableAssignment
          key={round}
          round={round}
          pairings={byRound[round]}
        />
      ))}
    </div>
  );
}