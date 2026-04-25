# Login Flow Bug Fix Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Fix two login flow bugs: missing redirect after successful login and auth token not persisted to localStorage.

**Architecture:** Two focused changes: (1) Add `useNavigate` to LoginForm to redirect after login, (2) Persist auth token to localStorage in api.ts with restore on init.

**Tech Stack:** React Router (useNavigate), localStorage, TypeScript

---

## Task 1: Persist Auth Token to localStorage in api.ts

**Files:**
- Modify: `frontend/src/services/api.ts:1-13`

- [ ] **Step 1: Write the failing test**

Create test file to verify token persistence:

```typescript
// frontend/src/services/api.test.ts
import { setAuthToken, getAuthToken } from './api';

describe('api token persistence', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('should store token in localStorage when setAuthToken is called', () => {
    setAuthToken('test-token-123');
    expect(localStorage.getItem('auth_token')).toBe('test-token-123');
  });

  it('should return token from localStorage on getAuthToken', () => {
    localStorage.setItem('auth_token', 'stored-token-456');
    expect(getAuthToken()).toBe('stored-token-456');
  });

  it('should return null when no token exists', () => {
    localStorage.clear();
    expect(getAuthToken()).toBe(null);
  });

  it('should remove token from localStorage when setAuthToken(null) is called', () => {
    localStorage.setItem('auth_token', 'to-be-removed');
    setAuthToken(null);
    expect(localStorage.getItem('auth_token')).toBe(null);
  });
});
```

Run: `cd frontend && npm test -- --testPathPattern="api.test.ts" --watchAll=false`
Expected: FAIL - setAuthToken/getAuthToken don't read/write localStorage yet

- [ ] **Step 2: Implement localStorage persistence**

Modify `frontend/src/services/api.ts` lines 1-13:

```typescript
const TOKEN_KEY = 'auth_token';

let authToken: string | null = localStorage.getItem(TOKEN_KEY);

export function setAuthToken(token: string | null) {
  authToken = token;
  if (token) {
    localStorage.setItem(TOKEN_KEY, token);
  } else {
    localStorage.removeItem(TOKEN_KEY);
  }
}

export function getAuthToken(): string | null {
  return authToken;
}
```

- [ ] **Step 3: Run test to verify it passes**

Run: `cd frontend && npm test -- --testPathPattern="api.test.ts" --watchAll=false`
Expected: PASS - all 4 tests green

- [ ] **Step 4: Commit**

```bash
git add frontend/src/services/api.ts frontend/src/services/api.test.ts
git commit -m "fix(api): persist auth token to localStorage

- Store token in localStorage on setAuthToken()
- Remove token from localStorage when set to null
- Restore token from localStorage on module init"
```

---

## Task 2: Add Navigation After Login in LoginForm.tsx

**Files:**
- Modify: `frontend/src/components/auth/LoginForm.tsx:1-33`

- [ ] **Step 1: Write the failing test**

Create test file:

```typescript
// frontend/src/components/auth/LoginForm.test.tsx
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { LoginForm } from './LoginForm';
import { useAuth } from '../../hooks/useAuth';

// Mock useAuth
vi.mock('../../hooks/useAuth', () => ({
  useAuth: () => ({
    login: vi.fn().mockResolvedValue(undefined),
  }),
}));

describe('LoginForm navigation', () => {
  it('should redirect to home page after successful login', async () => {
    const user = userEvent.setup();
    const loginMock = vi.fn().mockResolvedValue(undefined);

    vi.mocked(useAuth).mockImplementation(() => ({
      login: loginMock,
    }));

    render(
      <MemoryRouter>
        <LoginForm />
      </MemoryRouter>
    );

    await user.type(screen.getByLabelText(/email/i), 'test@example.com');
    await user.type(screen.getByLabelText(/password/i), 'password123');
    await user.click(screen.getByRole('button', { name: /sign in/i }));

    // After successful login, should have navigated to home
    expect(window.location.pathname).toBe('/');
  });
});
```

Run: `cd frontend && npm test -- --testPathPattern="LoginForm.test.tsx" --watchAll=false`
Expected: FAIL - navigate is not called

- [ ] **Step 2: Add useNavigate to LoginForm**

Modify `frontend/src/components/auth/LoginForm.tsx`:

```typescript
import { useState, FormEvent, ChangeEvent } from 'react';
import { useAuth } from '../../hooks/useAuth';
import { Input } from '../ui/Input';
import { Button } from '../ui/Button';
import { Link, useNavigate } from 'react-router-dom';

export function LoginForm() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    try {
      setError('');
      await login(email, password);
      navigate('/');
    } catch (err: any) {
      setError(err.message || 'Invalid credentials');
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4">
      <Input label="Email" type="email" value={email} onChange={(e: ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)} required />
      <Input label="Password" type="password" value={password} onChange={(e: ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)} required />
      {error && <p className="text-red-600 text-sm">{error}</p>}
      <Button type="submit">Sign in</Button>
      <p className="text-sm text-center text-gray-600">
        No account? <Link to="/register" className="text-blue-600 hover:underline">Register</Link>
      </p>
    </form>
  );
}
```

- [ ] **Step 3: Run test to verify it passes**

Run: `cd frontend && npm test -- --testPathPattern="LoginForm.test.tsx" --watchAll=false`
Expected: PASS - login form redirects after success

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/auth/LoginForm.tsx
git commit -m "fix(LoginForm): redirect to home page after successful login

- Import useNavigate from react-router-dom
- Add navigate('/') call after successful login"
```

---

## Verification Steps

After completing both tasks:

1. Run: `cd frontend && npm test -- --watchAll=false`
2. Start dev server: `cd frontend && npm run dev`
3. Open browser to http://localhost:5173/login
4. Login with credentials: `admin@ludo.local` / `changeme-in-production`
5. Verify: redirected to `/` after login
6. Verify: refresh page, still logged in (token persisted)