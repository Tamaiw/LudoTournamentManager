import { useState, FormEvent, ChangeEvent } from 'react';
import { useAuth } from '../../hooks/useAuth';
import { Input } from '../ui/Input';
import { Button } from '../ui/Button';
import { Link } from 'react-router-dom';

export function RegisterForm() {
  const { register } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [inviteCode, setInviteCode] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    try {
      setError('');
      await register(email, password, inviteCode);
    } catch (err: any) {
      setError(err.message || 'Registration failed');
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      <Input label="Email" type="email" value={email} onChange={(e: ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)} required />
      <Input label="Password" type="password" value={password} onChange={(e: ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)} required />
      <Input label="Invite Code" value={inviteCode} onChange={(e: ChangeEvent<HTMLInputElement>) => setInviteCode(e.target.value)} required />
      {error && <p className="text-red-600 text-sm">{error}</p>}
      <Button type="submit">Create account</Button>
      <p className="text-sm text-center text-gray-600">
        Have an account? <Link to="/login" className="text-blue-600 hover:underline">Sign in</Link>
      </p>
    </form>
  );
}