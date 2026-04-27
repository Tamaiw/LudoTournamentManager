/// <reference types="vitest" />
import { render, screen, fireEvent } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { LoginForm } from './LoginForm';
import { useAuth } from '../../hooks/useAuth';
import { describe, it, expect, vi, beforeEach } from 'vitest';

vi.mock('../../hooks/useAuth', () => ({
  useAuth: vi.fn(),
}));

const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

describe('LoginForm navigation', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should redirect to home page after successful login', async () => {
    const loginMock = vi.fn().mockResolvedValue(undefined);

    vi.mocked(useAuth).mockImplementation(() => ({
      login: loginMock,
      user: null,
      isLoading: false,
      register: vi.fn(),
      logout: vi.fn(),
    }));

    render(
      <MemoryRouter>
        <LoginForm />
      </MemoryRouter>
    );

    const inputs = document.querySelectorAll('input');
    expect(inputs.length).toBe(2);

    fireEvent.change(inputs[0], { target: { value: 'test@example.com' } });
    fireEvent.change(inputs[1], { target: { value: 'password123' } });

    fireEvent.click(screen.getByRole('button', { name: /sign in/i }));

    await vi.waitFor(() => {
      expect(loginMock).toHaveBeenCalledWith('test@example.com', 'password123');
      expect(mockNavigate).toHaveBeenCalledWith('/');
    });
  });
});