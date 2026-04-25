import { useParams } from 'react-router-dom';
import { useTournament, useTournamentMatches, useTournamentPairings } from '../hooks/useTournament';
import { Tabs } from '../components/ui/Tabs';
import { Card } from '../components/ui/Card';
import { Badge } from '../components/ui/Badge';
import { BracketView } from '../components/tournament/BracketView';
import { MatchCard } from '../components/tournament/MatchCard';

export function TournamentDetailPage() {
  const { id } = useParams<{ id: string }>();
  if (!id) return <div>Invalid tournament ID</div>;

  const { data: tournament, isLoading } = useTournament(id);
  const { data: matches } = useTournamentMatches(id);
  const { data: pairings } = useTournamentPairings(id);

  if (isLoading) return <div>Loading...</div>;
  if (!tournament) return <div>Tournament not found</div>;

  const tabs = [
    {
      id: 'bracket',
      label: 'Bracket',
      content: pairings ? (
        <BracketView pairings={pairings} />
      ) : (
        <p className="text-gray-500">No bracket data available</p>
      ),
    },
    {
      id: 'matches',
      label: 'Matches',
      content: matches && matches.length > 0 ? (
        <div className="flex flex-col gap-2">
          {matches.map(m => (
            <MatchCard key={m.id} match={m} />
          ))}
        </div>
      ) : (
        <p className="text-gray-500">No matches yet</p>
      ),
    },
  ];

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-3">
        <h1 className="text-2xl font-bold">{tournament.name}</h1>
        <Badge variant={tournament.status as any}>{tournament.status}</Badge>
      </div>
      <Tabs tabs={tabs} />
    </div>
  );
}