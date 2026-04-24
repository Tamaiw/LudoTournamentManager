import { Card } from '../components/ui/Card';
import { LoginForm } from '../components/auth/LoginForm';

export function LoginPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <Card className="w-full max-w-md">
        <h1 className="text-2xl font-bold mb-6 text-center">Sign in</h1>
        <LoginForm />
      </Card>
    </div>
  );
}