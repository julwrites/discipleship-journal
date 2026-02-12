import { renderHook, waitFor, act } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, Mock, afterEach } from 'vitest';
import { useAuth } from './useAuth';
import { onAuthStateChanged, User } from 'firebase/auth';

// Mock firebase/auth
vi.mock('firebase/auth', () => ({
  getAuth: vi.fn(),
  onAuthStateChanged: vi.fn(),
}));

// Mock the firebase lib which initializes app
vi.mock('@/lib/firebase', () => ({
  auth: {},
}));

describe('useAuth', () => {
  const originalEnv = process.env;

  beforeEach(() => {
    vi.clearAllMocks();
    vi.resetModules();
    process.env = { ...originalEnv };
  });

  afterEach(() => {
    process.env = originalEnv;
    vi.unstubAllEnvs();
  });

  it('should start with loading true and null user', () => {
    (onAuthStateChanged as Mock).mockImplementation(() => () => {});
    const { result } = renderHook(() => useAuth());
    expect(result.current.loading).toBe(true);
    expect(result.current.user).toBeNull();
  });

  it('should update user and loading when auth state changes', async () => {
    const mockUser = { uid: '123', email: 'test@example.com' } as User;
    (onAuthStateChanged as Mock).mockImplementation((_auth: unknown, callback: (user: User | null) => void) => {
      callback(mockUser);
      return () => {};
    });

    const { result } = renderHook(() => useAuth());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
      expect(result.current.user).toEqual(mockUser);
    });
  });

  it('should handle auth state change error', async () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    (onAuthStateChanged as Mock).mockImplementation((_auth: unknown, _cb: unknown, errorCallback: (error: Error) => void) => {
      errorCallback(new Error('Auth error'));
      return () => {};
    });

    const { result } = renderHook(() => useAuth());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
      expect(result.current.user).toBeNull();
    });

    expect(consoleSpy).toHaveBeenCalledWith('Auth state change error:', expect.any(Error));
    consoleSpy.mockRestore();
  });

  it('should bypass auth if mock key is present', async () => {
     vi.stubEnv('FIREBASE_API_KEY', 'mock-key');
     vi.useFakeTimers();

     const { result } = renderHook(() => useAuth());

     expect(result.current.loading).toBe(true);

     await act(async () => {
         await vi.advanceTimersByTimeAsync(1100);
     });

     expect(result.current.loading).toBe(false);

     vi.useRealTimers();
  });
});
