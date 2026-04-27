/// <reference types="vitest" />
import { setAuthToken } from './api';

describe('api token persistence', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('should store token in localStorage when setAuthToken is called', () => {
    setAuthToken('test-token-123');
    expect(localStorage.getItem('auth_token')).toBe('test-token-123');
  });

  it('should remove token from localStorage when setAuthToken(null) is called', () => {
    setAuthToken('some-token');
    setAuthToken(null);
    expect(localStorage.getItem('auth_token')).toBe(null);
  });

  it('should update localStorage when token is changed', () => {
    setAuthToken('first-token');
    expect(localStorage.getItem('auth_token')).toBe('first-token');
    setAuthToken('second-token');
    expect(localStorage.getItem('auth_token')).toBe('second-token');
  });
});