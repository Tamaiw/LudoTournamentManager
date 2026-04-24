import { useState, FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../services/api';
import { Card } from '../components/ui/Card';
import { Input } from '../components/ui/Input';
import { Button } from '../components/ui/Button';

export function CreateLeaguePage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [name, setName] = useState('');
  const [gamesPerPlayer, setGamesPerPlayer] = useState(3);
  const [tablesCount, setTablesCount] = useState(10);
  const [error, setError] = useState('');

  const mutation = useMutation({
    mutationFn: (data: { name: string; settings: { scoringRules: { placement: number; points: number }[]; gamesPerPlayer: number; tablesCount: number } }) =>
      api.createLeague(data),
    onSuccess: (league) => {
      queryClient.invalidateQueries({ queryKey: ['leagues'] });
      navigate(`/leagues/${league.id}`);
    },
    onError: (err: any) => setError(err.message),
  });

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    setError('');
    mutation.mutate({
      name,
      settings: {
        scoringRules: [
          { placement: 1, points: 3 },
          { placement: 2, points: 2 },
          { placement: 3, points: 1 },
          { placement: 4, points: 0 },
        ],
        gamesPerPlayer,
        tablesCount,
      },
    });
  };

  return (
    <div className="max-w-lg">
      <h1 className="text-2xl font-bold mb-6">Create League</h1>
      <Card>
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <Input
            label="League Name"
            value={name}
            onChange={e => setName(e.target.value)}
            required
            placeholder="Spring League 2026"
          />
          <Input
            label="Games Per Player"
            type="number"
            min={1}
            value={gamesPerPlayer}
            onChange={e => setGamesPerPlayer(Number(e.target.value))}
            required
          />
          <Input
            label="Number of Tables"
            type="number"
            min={1}
            value={tablesCount}
            onChange={e => setTablesCount(Number(e.target.value))}
            required
          />
          {error && <p className="text-red-600 text-sm">{error}</p>}
          <div className="flex gap-2">
            <Button type="submit" disabled={mutation.isPending}>
              {mutation.isPending ? 'Creating...' : 'Create'}
            </Button>
            <Button type="button" variant="ghost" onClick={() => navigate('/leagues')}>
              Cancel
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
}