import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { api, setAuthToken } from '../services/api';
import { User } from '../types';

interface AuthContextValue {
  user: User | null;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string, inviteCode: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    api.getMe().then(setUser).catch(() => setUser(null)).finally(() => setIsLoading(false));
  }, []);

  const login = async (email: string, password: string) => {
    const { token } = await api.login(email, password);
    setAuthToken(token);
    const user = await api.getMe();
    setUser(user);
  };

  const register = async (email: string, password: string, inviteCode: string) => {
    const { token } = await api.register(email, password, inviteCode);
    setAuthToken(token);
    const user = await api.getMe();
    setUser(user);
  };

  const logout = async () => {
    try {
      await api.logout();
    } catch {
      // Ignore logout errors
    }
    setAuthToken(null);
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, isLoading, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}