# Bug Report: Login Flow Issues in Frontend

**Date:** 2026-04-25
**Reported by:** Claude Code
**Severity:** Medium
**Status:** Open

---

## Summary

The login flow on the frontend has two related issues that prevent users from successfully authenticating and accessing protected routes.

---

## Issue 1: Missing Redirect After Successful Login

**Component:** `frontend/src/components/auth/LoginForm.tsx`

**Description:**
After a user successfully submits the login form and the API returns a valid token, the application remains on the `/login` page instead of redirecting to the dashboard (`/`).

**Root Cause:**
The `LoginForm` component calls `login()` from `useAuth()` but does not navigate to the home page after a successful login. The `AuthContext` correctly stores the user state, but the UI never transitions to a protected route.

**Expected Behavior:**
After login success, user should be redirected to `/`.

**Actual Behavior:**
User stays on `/login` page, seeing the sign-in form with no indication of success.

**Affected Files:**
- [frontend/src/components/auth/LoginForm.tsx](frontend/src/components/auth/LoginForm.tsx)

---

## Issue 2: Auth Token Not Persisted to localStorage

**Component:** `frontend/src/services/api.ts`

**Description:**
The authentication token is stored in a module-level JavaScript variable (`let authToken = null`) which is lost on page reload or navigation. This causes users to appear logged out immediately after any navigation.

**Root Cause:**
Tokens are stored in memory only. The code does not write tokens to `localStorage` for persistence across page reloads.

**Expected Behavior:**
Token should be persisted to `localStorage` and restored on page load.

**Actual Behavior:**
On page reload, `getAuthToken()` returns `null`, causing all API calls to fail with 401 errors and the user to appear logged out.

**Affected Files:**
- [frontend/src/services/api.ts](frontend/src/services/api.ts)

---

## Test Results

| Test | Result |
|------|--------|
| Login API accepts valid credentials | ✅ Pass |
| Login API returns JWT token | ✅ Pass |
| Token stored after login | ✅ Pass (in-memory) |
| Redirect to dashboard after login | ❌ Fail |
| Token persists on page reload | ❌ Fail |
| Protected routes accessible after reload | ❌ Fail |

### Test Credentials Used
- **Email:** `admin@ludo.local`
- **Password:** `changeme-in-production`

### API Response
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

## Recommended Fixes

### Fix 1: Add Navigation After Login

In `LoginForm.tsx`, add `useNavigate` to redirect after successful login:

```tsx
import { useNavigate } from 'react-router-dom';

export function LoginForm() {
  const { login } = useAuth();
  const navigate = useNavigate();
  // ...

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
  // ...
}
```

### Fix 2: Persist Token to localStorage

In `api.ts`, modify token handling to persist across sessions:

```tsx
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
```

---

## Impact

- Users cannot log in through the UI (stuck on login page after form submission)
- Even if login were fixed, users would be logged out on every page reload
- Affects all user-facing functionality for authenticated users
