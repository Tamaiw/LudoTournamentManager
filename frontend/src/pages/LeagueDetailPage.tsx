import { useState } from 'react';
import { useParams } from 'react-router-dom';
import { useLeague, useLeagueStandings, useGenerateLeaguePairings } from '../hooks/useLeague';
import { Tabs } from '../components/ui/Tabs';
import { Badge } from '../components/ui/Badge';
import { StandingsTable } from '../components/league/StandingsTable';
import { Button } from '../components/ui/Button';

export function LeagueDetailPage() {
  const { id } = useParams<{ id: string }>();
  const [playDate] = useState(() => new Date().toISOString().split('T')[0]);

  if (!id) return <div>Invalid league ID</div>;

  const { data: league, isLoading } = useLeague(id);
  const { data: standingsData } = useLeagueStandings(id);
  const generatePairings = useGenerateLeaguePairings();

  if (isLoading) return <div>Loading...</div>;
  if (!league) return <div>League not found</div>;

  const tabs = [
    {
      id: 'standings',
      label: 'Standings',
      content: standingsData && standingsData.length > 0 ? (
        <StandingsTable standings={standingsData} />
      ) : (
        <p className="text-gray-500">No standings available</p>
      ),
    },
    {
      id: 'generate',
      label: 'Generate Pairings',
      content: (
        <div className="flex flex-col gap-4">
          <p className="text-gray-600">Generate fair pairings for the next play date.</p>
          <div className="flex items-center gap-4">
            <Button
              onClick={() => generatePairings.mutate({ leagueId: id, playDate })}
              disabled={generatePairings.isPending}
            >
              {generatePairings.isPending ? 'Generating...' : 'Generate Pairings'}
            </Button>
            <span className="text-sm text-gray-500">Play date: {playDate}</span>
          </div>
          {generatePairings.isError && (
            <p className="text-red-600 text-sm">Failed to generate pairings</p>
          )}
          <p className="text-xs text-gray-400">Note: Play dates management not yet implemented</p>
        </div>
      ),
    },
  ];

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-3">
        <h1 className="text-2xl font-bold">{league.name}</h1>
        <Badge variant={league.status as any}>{league.status}</Badge>
      </div>
      <p className="text-gray-600">
        {league.settings.gamesPerPlayer} games/player · {league.settings.tablesCount} tables
      </p>
      <Tabs tabs={tabs} />
    </div>
  );
}