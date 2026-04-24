import { Link } from 'react-router-dom';
import { Button } from '../ui/Button';
import { Role } from '../../types';

interface QuickActionsProps {
  role: Role;
}

export function QuickActions({ role }: QuickActionsProps) {
  return (
    <div className="flex flex-wrap gap-2">
      {role !== 'guest' && (
        <>
          <Link to="/tournaments/new">
            <Button variant="secondary">Create Tournament</Button>
          </Link>
          <Link to="/leagues/new">
            <Button variant="secondary">Create League</Button>
          </Link>
        </>
      )}
      <Link to="/tournaments">
        <Button variant="ghost">Browse Tournaments</Button>
      </Link>
      <Link to="/leagues">
        <Button variant="ghost">Browse Leagues</Button>
      </Link>
    </div>
  );
}