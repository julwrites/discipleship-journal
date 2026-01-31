import { render, screen, fireEvent, waitFor, within } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { MemoryRouter } from 'react-router-dom';
import LoginPage from './Login';
import { auth } from '@/lib/firebase';
import { signInWithEmailAndPassword, signInWithPopup, signInWithRedirect, sendSignInLinkToEmail, isSignInWithEmailLink, signInWithEmailLink } from 'firebase/auth';
import userEvent from '@testing-library/user-event';

// Mock Firebase Auth
vi.mock('@/lib/firebase', () => ({
    auth: {
        currentUser: null,
    },
    authReadyPromise: null,
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
    sendSignInLinkToEmail: vi.fn(),
    isSignInWithEmailLink: vi.fn().mockReturnValue(false),
    signInWithEmailLink: vi.fn(),
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

// Mock window.prompt
window.prompt = vi.fn();

describe('LoginPage', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        window.localStorage.clear();
    });

    it('renders the login page correctly', () => {
        render(
            <MemoryRouter>
                <LoginPage />
            </MemoryRouter>
        );
        expect(screen.getByText('Discipleship Journal')).toBeInTheDocument();
        expect(screen.getByRole('button', { name: /Google/i })).toBeInTheDocument();
    });

    it('handles email login', async () => {
        render(
            <MemoryRouter>
                <LoginPage />
            </MemoryRouter>
        );

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
        render(
            <MemoryRouter>
                <LoginPage />
            </MemoryRouter>
        );

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

        render(
            <MemoryRouter>
                <LoginPage />
            </MemoryRouter>
        );

        const googleButton = screen.getByRole('button', { name: /Google/i });
        fireEvent.click(googleButton);

        await waitFor(() => {
            expect(signInWithPopup).toHaveBeenCalled();
            // Should fallback to redirect
            expect(signInWithRedirect).toHaveBeenCalled();
        });
    });

    it('sends magic link', async () => {
        render(
            <MemoryRouter>
                <LoginPage />
            </MemoryRouter>
        );

        // Click Magic Link tab
        const magicLinkTab = screen.getByRole('tab', { name: 'Magic Link' });
        await userEvent.click(magicLinkTab);

        const sendLinkButton = screen.getByRole('button', { name: 'Send Magic Link' });
        const form = sendLinkButton.closest('form');
        expect(form).toBeInTheDocument();

        const emailInput = within(form!).getByPlaceholderText('name@example.com');
        fireEvent.change(emailInput, { target: { value: 'magic@example.com' } });
        fireEvent.click(sendLinkButton);

        await waitFor(() => {
            expect(sendSignInLinkToEmail).toHaveBeenCalledWith(auth, 'magic@example.com', expect.objectContaining({
                handleCodeInApp: true,
                url: expect.stringContaining('/login')
            }));
            expect(window.localStorage.getItem('emailForSignIn')).toBe('magic@example.com');
        });
    });

    it('completes sign in with email link', async () => {
        // Mock isSignInWithEmailLink to return true
        vi.mocked(isSignInWithEmailLink).mockReturnValue(true);
        // Mock stored email
        window.localStorage.setItem('emailForSignIn', 'stored@example.com');
        // Mock successful sign in
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        vi.mocked(signInWithEmailLink).mockResolvedValue({ user: {} } as any);

        render(
            <MemoryRouter>
                <LoginPage />
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(signInWithEmailLink).toHaveBeenCalledWith(auth, 'stored@example.com', expect.any(String));
            expect(window.localStorage.getItem('emailForSignIn')).toBeNull(); // Should be cleared
        });
    });

    it('prompts for email if not in storage during link sign in', async () => {
        vi.mocked(isSignInWithEmailLink).mockReturnValue(true);
        window.localStorage.removeItem('emailForSignIn');
        vi.mocked(window.prompt).mockReturnValue('prompted@example.com');
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        vi.mocked(signInWithEmailLink).mockResolvedValue({ user: {} } as any);

        render(
            <MemoryRouter>
                <LoginPage />
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(window.prompt).toHaveBeenCalled();
            expect(signInWithEmailLink).toHaveBeenCalledWith(auth, 'prompted@example.com', expect.any(String));
        });
    });
});
