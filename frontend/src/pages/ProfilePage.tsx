import { useAuth } from '../hooks/useAuth';
import { Card } from '../components/ui/Card';
import { Badge } from '../components/ui/Badge';
import { Button } from '../components/ui/Button';

export function ProfilePage() {
  const { user, logout } = useAuth();

  if (!user) return <div>Not logged in</div>;

  return (
    <div className="max-w-lg flex flex-col gap-6">
      <h1 className="text-2xl font-bold">Profile</h1>
      <Card>
        <div className="flex flex-col gap-3">
          <div>
            <p className="text-sm text-gray-500">Email</p>
            <p className="font-medium">{user.email}</p>
          </div>
          <div>
            <p className="text-sm text-gray-500">Role</p>
            <Badge variant={user.role}>{user.role}</Badge>
          </div>
          <div>
            <p className="text-sm text-gray-500">Member since</p>
            <p className="font-medium">{new Date(user.createdAt).toLocaleDateString()}</p>
          </div>
        </div>
      </Card>
      <Button variant="danger" onClick={logout}>Logout</Button>
    </div>
  );
}