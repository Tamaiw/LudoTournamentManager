import { Link } from 'react-router-dom';
import { Card } from '../ui/Card';
import { Badge } from '../ui/Badge';
import { Button } from '../ui/Button';

interface NextEventCardProps {
  event: { id: string; name: string; type: 'tournament' | 'league'; status: string } | null;
}

export function NextEventCard({ event }: NextEventCardProps) {
  if (!event) {
    return (
      <Card>
        <h2 className="text-lg font-semibold mb-2">Next Event</h2>
        <p className="text-gray-500 mb-4">No active events</p>
        <div className="flex gap-2">
          <Link to="/tournaments/new">
            <Button variant="secondary">Create Tournament</Button>
          </Link>
          <Link to="/leagues/new">
            <Button variant="secondary">Create League</Button>
          </Link>
        </div>
      </Card>
    );
  }

  return (
    <Card className="border-l-4 border-l-blue-500">
      <div className="flex items-center gap-2 mb-1">
        <Badge variant={event.status as any}>{event.status}</Badge>
        <span className="text-xs text-gray-500">{event.type}</span>
      </div>
      <h2 className="text-xl font-semibold mb-2">{event.name}</h2>
      <Link to={event.type === 'tournament' ? `/tournaments/${event.id}` : `/leagues/${event.id}`}>
        <Button>Enter</Button>
      </Link>
    </Card>
  );
}