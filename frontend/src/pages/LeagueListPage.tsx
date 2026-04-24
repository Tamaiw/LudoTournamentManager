import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useLeagues } from '../hooks/useLeagues';
import { useAuth } from '../hooks/useAuth';
import { LeagueCard } from '../components/league/LeagueCard';
import { Button } from '../components/ui/Button';
import { Card } from '../components/ui/Card';

type Filter = 'all' | 'live' | 'completed' | 'draft';

export function LeagueListPage() {
  const { user } = useAuth();
  const [filter, setFilter] = useState<Filter>('all');
  const { data: leagues, isLoading } = useLeagues();

  const filtered = leagues?.filter(l => filter === 'all' ? true : l.status === filter) || [];

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Leagues</h1>
        {user?.role !== 'guest' && (
          <Link to="/leagues/new">
            <Button>Create League</Button>
          </Link>
        )}
      </div>

      <div className="flex gap-2">
        {(['all', 'live', 'completed', 'draft'] as Filter[]).map(f => (
          <button
            key={f}
            onClick={() => setFilter(f)}
            className={`px-3 py-1 rounded text-sm font-medium transition-colors ${filter === f ? 'bg-blue-100 text-blue-700' : 'text-gray-600 hover:bg-gray-100'}`}
          >
            {f.charAt(0).toUpperCase() + f.slice(1)}
          </button>
        ))}
      </div>

      {isLoading ? (
        <p className="text-gray-500">Loading...</p>
      ) : filtered.length === 0 ? (
        <Card>
          <p className="text-gray-500 text-center py-8">No leagues found</p>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {filtered.map(l => <LeagueCard key={l.id} league={l} />)}
        </div>
      )}
    </div>
  );
}