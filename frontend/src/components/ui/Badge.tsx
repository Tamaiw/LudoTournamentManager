export type BadgeVariant = 'live' | 'draft' | 'completed' | 'admin' | 'member' | 'guest';

interface BadgeProps {
  variant: BadgeVariant;
  children: string;
}

export function Badge({ variant, children }: BadgeProps) {
  const styles: Record<BadgeVariant, string> = {
    live: 'bg-green-100 text-green-800',
    draft: 'bg-gray-100 text-gray-800',
    completed: 'bg-blue-100 text-blue-800',
    admin: 'bg-purple-100 text-purple-800',
    member: 'bg-blue-100 text-blue-800',
    guest: 'bg-gray-100 text-gray-800',
  };
  return (
    <span className={`inline-block px-2 py-0.5 rounded text-xs font-medium ${styles[variant]}`}>
      {children}
    </span>
  );
}