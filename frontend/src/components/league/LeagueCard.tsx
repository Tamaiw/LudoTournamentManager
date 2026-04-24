import { Link } from 'react-router-dom';
import { Card } from '../ui/Card';
import { Badge } from '../ui/Badge';
import { League } from '../../types';

interface LeagueCardProps {
  league: League;
}

export function LeagueCard({ league }: LeagueCardProps) {
  return (
    <Link to={`/leagues/${league.id}`}>
      <Card className="hover:shadow-md transition-shadow cursor-pointer">
        <div className="flex items-center justify-between">
          <div>
            <p className="font-medium text-lg">{league.name}</p>
            <p className="text-sm text-gray-500">
              {league.settings.gamesPerPlayer} games/player · {league.settings.tablesCount} tables
            </p>
          </div>
          <Badge variant={league.status as any}>{league.status}</Badge>
        </div>
      </Card>
    </Link>
  );
}