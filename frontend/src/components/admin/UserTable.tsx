import { User, Role } from '../../types';
import { Button } from '../ui/Button';
import { useUpdateUser, useDeleteUser } from '../../hooks/useUsers';

interface UserTableProps {
  users: User[];
}

export function UserTable({ users }: UserTableProps) {
  const updateUser = useUpdateUser();
  const deleteUser = useDeleteUser();

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b text-left">
            <th className="py-2 pr-4 font-medium text-gray-600">Email</th>
            <th className="py-2 pr-4 font-medium text-gray-600">Role</th>
            <th className="py-2 pr-4 font-medium text-gray-600">Last Active</th>
            <th className="py-2 font-medium text-gray-600">Actions</th>
          </tr>
        </thead>
        <tbody>
          {users.map(u => (
            <tr key={u.id} className="border-b last:border-0">
              <td className="py-2 pr-4">{u.email}</td>
              <td className="py-2 pr-4">
                <select
                  value={u.role}
                  onChange={e => updateUser.mutate({ id: u.id, role: e.target.value as Role })}
                  disabled={updateUser.isPending}
                  className="border rounded px-2 py-1 text-sm"
                >
                  <option value="guest">Guest</option>
                  <option value="member">Member</option>
                  <option value="admin">Admin</option>
                </select>
              </td>
              <td className="py-2 pr-4 text-gray-500">
                {u.lastActive ? new Date(u.lastActive).toLocaleDateString() : 'Never'}
              </td>
              <td className="py-2">
                <Button
                  variant="danger"
                  onClick={() => {
                    if (confirm(`Delete user ${u.email}?`)) {
                      deleteUser.mutate(u.id);
                    }
                  }}
                >
                  Delete
                </Button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}