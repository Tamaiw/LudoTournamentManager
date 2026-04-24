import { Link } from 'react-router-dom';
import { Card } from '../ui/Card';
import { Badge } from '../ui/Badge';
import { Tournament } from '../../types';

interface TournamentCardProps {
  tournament: Tournament;
}

export function TournamentCard({ tournament }: TournamentCardProps) {
  return (
    <Link to={`/tournaments/${tournament.id}`}>
      <Card className="hover:shadow-md transition-shadow cursor-pointer">
        <div className="flex items-center justify-between">
          <div>
            <p className="font-medium text-lg">{tournament.name}</p>
            <p className="text-sm text-gray-500">
              {tournament.settings.tablesCount} tables
            </p>
          </div>
          <Badge variant={tournament.status as any}>{tournament.status}</Badge>
        </div>
      </Card>
    </Link>
  );
}