import { useState, FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../services/api';
import { Card } from '../components/ui/Card';
import { Input } from '../components/ui/Input';
import { Button } from '../components/ui/Button';

export function CreateTournamentPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [name, setName] = useState('');
  const [tablesCount, setTablesCount] = useState(10);
  const [error, setError] = useState('');

  const mutation = useMutation({
    mutationFn: (data: { name: string; settings: { tablesCount: number } }) => api.createTournament(data),
    onSuccess: (tournament) => {
      queryClient.invalidateQueries({ queryKey: ['tournaments'] });
      navigate(`/tournaments/${tournament.id}`);
    },
    onError: (err: any) => setError(err.message),
  });

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    setError('');
    mutation.mutate({ name, settings: { tablesCount } });
  };

  return (
    <div className="max-w-lg">
      <h1 className="text-2xl font-bold mb-6">Create Tournament</h1>
      <Card>
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <Input
            label="Tournament Name"
            value={name}
            onChange={e => setName(e.target.value)}
            required
            placeholder="Spring Championship 2026"
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
            <Button type="button" variant="ghost" onClick={() => navigate('/tournaments')}>
              Cancel
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
}