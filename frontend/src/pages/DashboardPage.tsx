import { useAuth } from '../hooks/useAuth';
import { api } from '../services/api';
import { useQuery } from '@tanstack/react-query';
import { NextEventCard } from '../components/dashboard/NextEventCard';
import { QuickActions } from '../components/dashboard/QuickActions';
import { Card } from '../components/ui/Card';
import { Badge, BadgeVariant } from '../components/ui/Badge';
import { Link } from 'react-router-dom';

export function DashboardPage() {
  const { user } = useAuth();

  const { data: tournaments } = useQuery({
    queryKey: ['tournaments'],
    queryFn: () => api.listTournaments(),
  });

  const { data: leagues } = useQuery({
    queryKey: ['leagues'],
    queryFn: () => api.listLeagues(),
  });

  const liveTournament = tournaments?.find(t => t.status === 'live');
  const liveLeague = leagues?.find(l => l.status === 'live');

  const nextEvent = liveTournament
    ? { id: liveTournament.id, name: liveTournament.name, type: 'tournament' as const, status: liveTournament.status }
    : liveLeague
    ? { id: liveLeague.id, name: liveLeague.name, type: 'league' as const, status: liveLeague.status }
    : null;

  const recentTournaments = tournaments?.slice(0, 3) || [];
  const recentLeagues = leagues?.slice(0, 3) || [];

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-bold">Welcome back, {user?.email}</h1>
        <p className="text-gray-500">{new Date().toLocaleDateString('en-US', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}</p>
      </div>

      <NextEventCard event={nextEvent} />

      <div>
        <h2 className="text-lg font-semibold mb-3">Quick Actions</h2>
        <QuickActions role={user?.role || 'guest'} />
      </div>

      {user?.role === 'admin' && (
        <Card>
          <div className="flex items-center justify-between">
            <div>
              <h2 className="font-semibold">Admin Panel</h2>
              <p className="text-sm text-gray-500">Manage users and invitations</p>
            </div>
            <Link to="/admin" className="text-blue-600 hover:underline text-sm">Open →</Link>
          </div>
        </Card>
      )}

      {(recentTournaments.length > 0 || recentLeagues.length > 0) && (
        <div>
          <h2 className="text-lg font-semibold mb-3">Recent Activity</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {recentTournaments.map(t => (
              <Card key={t.id}>
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium">{t.name}</p>
                    <Badge variant={t.status as BadgeVariant}>{t.status}</Badge>
                  </div>
                  <Link to={`/tournaments/${t.id}`} className="text-blue-600 hover:underline text-sm">View</Link>
                </div>
              </Card>
            ))}
            {recentLeagues.map(l => (
              <Card key={l.id}>
                <div className="flex items-center justify-between">
                  <div>
                    <p className="font-medium">{l.name}</p>
                    <Badge variant={l.status as BadgeVariant}>{l.status}</Badge>
                  </div>
                  <Link to={`/leagues/${l.id}`} className="text-blue-600 hover:underline text-sm">View</Link>
                </div>
              </Card>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}