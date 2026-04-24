import { useUsers } from '../hooks/useUsers';
import { UserTable } from '../components/admin/UserTable';
import { Card } from '../components/ui/Card';

export function AdminPage() {
  const { data, isLoading } = useUsers();

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-2xl font-bold">Admin Panel</h1>

      <div>
        <h2 className="text-lg font-semibold mb-3">Users</h2>
        {isLoading ? (
          <p className="text-gray-500">Loading...</p>
        ) : data?.users ? (
          <Card>
            <UserTable users={data.users} />
          </Card>
        ) : (
          <Card>
            <p className="text-gray-500 text-center py-4">No users found</p>
          </Card>
        )}
      </div>

      <div>
        <h2 className="text-lg font-semibold mb-3">Invitations</h2>
        <Card>
          <p className="text-gray-500 py-4 text-center">
            Invitation management not yet implemented (backend missing)
          </p>
        </Card>
      </div>
    </div>
  );
}