import { render, screen, waitFor, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import ConnectionsPage from './ConnectionsPage';
import { MemoryRouter } from 'react-router-dom';
import * as api from '@/services/api';
import { User } from 'firebase/auth';

// Mock sonner
vi.mock('sonner', () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
    }
}));

// Mock the API
vi.mock('@/services/api', () => ({
    getConnections: vi.fn(),
    searchUsers: vi.fn(),
    sendConnectionRequest: vi.fn(),
    respondToConnectionRequest: vi.fn(),
}));

// Mock Firebase
vi.mock('@/lib/firebase', () => ({
    auth: {
        signOut: vi.fn(),
    }
}));

// Mock useAuthState
vi.mock('react-firebase-hooks/auth', () => ({
    useAuthState: () => [{ email: 'test@example.com', uid: '123' } as User, false, undefined],
}));


describe('ConnectionsPage', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        vi.useRealTimers();
    });

    afterEach(() => {
        vi.useRealTimers();
    });

    it('debounces search in Find People tab', async () => {
        const user = userEvent.setup();
        const mockSearchUsers = vi.mocked(api.searchUsers);
        mockSearchUsers.mockResolvedValue([
            { id: 'user-1', email: 'found@example.com', username: 'Found User' }
        ]);

        const mockGetConnections = vi.mocked(api.getConnections);
        mockGetConnections.mockResolvedValue([]);

        render(
            <MemoryRouter>
                <ConnectionsPage />
            </MemoryRouter>
        );

        // Wait for initial load
        await waitFor(() => {
            expect(mockGetConnections).toHaveBeenCalled();
        });

        // Switch to Find People tab
        const findTab = screen.getByRole('tab', { name: /find people/i });
        await user.click(findTab);

        // Wait for the tab content to be visible
        const searchInput = await screen.findByPlaceholderText('Search by email or username...');

        // Now switch to fake timers for debounce testing
        vi.useFakeTimers();

        // Type "te" (length < 3)
        // We use user.type, but user.type uses real timers by default.
        // With fake timers enabled, we need to advance time for typing to happen if we use user.type?
        // Actually, userEvent with fake timers can be tricky.
        // It's often safer to use fireEvent.change for just the input change when controlling timers manually,
        // OR configure userEvent to use the fake timers.
        // Let's use fireEvent.change for simplicity in this specific debounce test block.
        // import fireEvent from testing-library/react if needed, or just use userEvent but be careful.
        // Let's stick to userEvent but we need to advance timers if it has delays.
        // Actually, let's use userEvent.type but wrapped in act maybe?
        // Or simpler: just use fireEvent.change for the input which is synchronous and we only care about the effect hook debounce.

        const { fireEvent } = await import('@testing-library/react');

        fireEvent.change(searchInput, { target: { value: 'te' } });

        act(() => {
            vi.advanceTimersByTime(600);
        });

        // Should NOT call search
        expect(mockSearchUsers).not.toHaveBeenCalled();

        // Type "test" (length >= 3)
        fireEvent.change(searchInput, { target: { value: 'test' } });

        // Shouldn't call immediately
        expect(mockSearchUsers).not.toHaveBeenCalled();

        act(() => {
            vi.advanceTimersByTime(500);
        });

        vi.useRealTimers();

        await waitFor(() => {
            expect(mockSearchUsers).toHaveBeenCalledWith('test');
        });

        // Verify results are shown
        expect(await screen.findByText('Found User')).toBeInTheDocument();
        expect(screen.getByText('found@example.com')).toBeInTheDocument();
    });
});
