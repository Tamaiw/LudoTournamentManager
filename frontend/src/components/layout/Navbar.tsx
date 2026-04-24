import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import { Badge, BadgeVariant } from '../ui/Badge';

export function Navbar() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  return (
    <nav className="bg-white shadow-sm border-b">
      <div className="max-w-5xl mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          <div className="flex items-center gap-6">
            <Link to="/" className="text-xl font-bold text-blue-600">Ludo</Link>
            <div className="flex gap-4 text-sm">
              <Link to="/tournaments" className="text-gray-600 hover:text-gray-900">Tournaments</Link>
              <Link to="/leagues" className="text-gray-600 hover:text-gray-900">Leagues</Link>
              {user?.role === 'admin' && (
                <Link to="/admin" className="text-gray-600 hover:text-gray-900">Admin</Link>
              )}
            </div>
          </div>
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              <span className="text-sm text-gray-600">{user?.email}</span>
              <Badge variant={(user?.role ?? 'guest') as BadgeVariant}>{(user?.role ?? 'guest')}</Badge>
            </div>
            <Link to="/profile" className="text-sm text-blue-600 hover:underline">Profile</Link>
            <button onClick={handleLogout} className="text-sm text-gray-600 hover:text-gray-900">Logout</button>
          </div>
        </div>
      </div>
    </nav>
  );
}