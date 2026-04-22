import { useQuery } from '@tanstack/react-query';
import { useParams, Link } from 'react-router-dom';
import { api } from '../services/api';

export function TournamentDetail() {
  const { id } = useParams<{ id: string }>();

  const { data: tournament, isLoading: loadingTournament } = useQuery({
    queryKey: ['tournament', id],
    queryFn: () => api.getTournament(id!),
    enabled: !!id,
  });

  const { data: matches } = useQuery({
    queryKey: ['tournament-matches', id],
    queryFn: () => api.getTournamentMatches(id!),
    enabled: !!id,
  });

  if (loadingTournament) return <div>Loading...</div>;
  if (!tournament) return <div>Tournament not found</div>;

  return (
    <div className="container mx-auto p-4">
      <div className="mb-6">
        <Link to="/tournaments" className="text-blue-600 hover:underline">
          ← Back to Tournaments
        </Link>
      </div>

      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <h1 className="text-2xl font-bold mb-4">{tournament.name}</h1>
        <div className="grid grid-cols-2 gap-4 text-sm">
          <div>
            <span className="font-semibold">Status:</span>{' '}
            <span className={`px-2 py-1 rounded text-sm ${
              tournament.status === 'live' ? 'bg-green-100' : 'bg-gray-100'
            }`}>
              {tournament.status}
            </span>
          </div>
          <div>
            <span className="font-semibold">Tables:</span> {tournament.settings.tablesCount}
          </div>
        </div>
      </div>

      <div className="mb-6">
        <h2 className="text-xl font-bold mb-4">Current Round</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          {matches?.map(match => (
            <MatchCard
              key={match.id}
              pairing={{
                gameId: match.id,
                round: match.round,
                tableNumber: match.tableNumber,
                playerIds: [], // Will be populated from match assignments
                seatColors: [],
                status: match.status,
              }}
            />
          ))}
        </div>
      </div>
    </div>
  );
}