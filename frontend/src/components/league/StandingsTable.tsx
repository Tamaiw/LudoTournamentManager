import { PlayerStanding } from '../../types';

interface StandingsTableProps {
  standings: PlayerStanding[];
}

export function StandingsTable({ standings }: StandingsTableProps) {
  if (!standings || standings.length === 0) {
    return <p className="text-gray-500">No standings available</p>;
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b text-left">
            <th className="py-2 pr-4 font-medium text-gray-600">Rank</th>
            <th className="py-2 pr-4 font-medium text-gray-600">Player</th>
            <th className="py-2 pr-4 font-medium text-gray-600">Games</th>
            <th className="py-2 pr-4 font-medium text-gray-600">Points</th>
            <th className="py-2 font-medium text-gray-600">Wins</th>
          </tr>
        </thead>
        <tbody>
          {standings.map(s => (
            <tr key={s.playerId} className="border-b last:border-0">
              <td className="py-2 pr-4 font-medium">#{s.rank}</td>
              <td className="py-2 pr-4">{s.displayName}</td>
              <td className="py-2 pr-4">{s.gamesPlayed}</td>
              <td className="py-2 pr-4">{s.totalPoints}</td>
              <td className="py-2">{s.wins}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}