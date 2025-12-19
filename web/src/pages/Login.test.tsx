import { render, screen, fireEvent, waitFor, within } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import LoginPage from './Login';
import { auth } from '@/lib/firebase';
import { signInWithEmailAndPassword, signInWithPopup, signInWithRedirect } from 'firebase/auth';

// Mock Firebase Auth
vi.mock('@/lib/firebase', () => ({
    auth: {
        currentUser: null,
    },
}));

// Mock Firebase Auth methods
vi.mock('firebase/auth', () => ({
    getAuth: vi.fn(),
    signInWithPopup: vi.fn(),
    signInWithRedirect: vi.fn(),
    GoogleAuthProvider: vi.fn(),
    signInWithEmailAndPassword: vi.fn(),
    createUserWithEmailAndPassword: vi.fn(),
    getRedirectResult: vi.fn().mockResolvedValue(null),
    AuthError: vi.fn(),
}));

// Mock Sonner
vi.mock('sonner', () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
        info: vi.fn(),
    }
}));

describe('LoginPage', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('renders the login page correctly', () => {
        render(<LoginPage />);
        expect(screen.getByText('Discipleship Journal')).toBeInTheDocument();
        expect(screen.getByRole('button', { name: /Google/i })).toBeInTheDocument();
    });

    it('handles email login', async () => {
        render(<LoginPage />);

        const signInButton = screen.getByRole('button', { name: 'Sign In with Email' });
        const form = signInButton.closest('form');
        expect(form).toBeInTheDocument();

        const emailInput = within(form!).getByPlaceholderText('name@example.com');
        const passwordInput = within(form!).getByPlaceholderText('Password');

        fireEvent.change(emailInput, { target: { value: 'test@example.com' } });
        fireEvent.change(passwordInput, { target: { value: 'password123' } });
        fireEvent.click(signInButton);

        await waitFor(() => {
            expect(signInWithEmailAndPassword).toHaveBeenCalledWith(auth, 'test@example.com', 'password123');
        });
    });

    it('handles google login with popup', async () => {
        render(<LoginPage />);

        const googleButton = screen.getByRole('button', { name: /Google/i });
        fireEvent.click(googleButton);

        await waitFor(() => {
            expect(signInWithPopup).toHaveBeenCalled();
            expect(signInWithRedirect).not.toHaveBeenCalled();
        });
    });

    it('handles google login popup blocked fallback', async () => {
        // Mock popup failure
        vi.mocked(signInWithPopup).mockRejectedValueOnce({
            code: 'auth/popup-blocked',
            message: 'Popup blocked'
        });

        render(<LoginPage />);

        const googleButton = screen.getByRole('button', { name: /Google/i });
        fireEvent.click(googleButton);

        await waitFor(() => {
            expect(signInWithPopup).toHaveBeenCalled();
            // Should fallback to redirect
            expect(signInWithRedirect).toHaveBeenCalled();
        });
    });
});
