import { render, screen, waitFor, act, fireEvent } from '@testing-library/react';
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
    getOrCreateDirectGroup: vi.fn(),
}));

// Mock Firebase
vi.mock('@/lib/firebase', () => ({
    auth: {
        signOut: vi.fn(),
    }
}));

// Mock useAuth
vi.mock('@/hooks/useAuth', () => ({
    useAuth: () => ({ user: { email: 'test@example.com', uid: '123' } as User, loading: false }),
}));


describe('ConnectionsPage', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        vi.useRealTimers();
    });

    afterEach(() => {
        vi.useRealTimers();
    });

    it('renders existing connections', async () => {
        const mockGetConnections = vi.mocked(api.getConnections);
        mockGetConnections.mockResolvedValue([
            {
                id: 'conn-1',
                requester_id: '123',
                receiver_id: '456',
                status: 'accepted',
                receiver_username: 'Friend User',
                receiver_email: 'friend@example.com',
                requester_username: 'Me',
                requester_email: 'test@example.com'
            }
        ]);

        render(
            <MemoryRouter>
                <ConnectionsPage />
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(screen.getByText('Friend User')).toBeInTheDocument();
        });
        expect(screen.getByText('friend@example.com')).toBeInTheDocument();
    });

    it('sends connection request', async () => {
        const user = userEvent.setup();
        const mockSearchUsers = vi.mocked(api.searchUsers);
        mockSearchUsers.mockResolvedValue([
            { id: 'user-2', email: 'stranger@example.com', username: 'Stranger' }
        ]);
        const mockSendRequest = vi.mocked(api.sendConnectionRequest);
        mockSendRequest.mockResolvedValue({ success: true });

        // Mock getConnections to return empty initially
        vi.mocked(api.getConnections).mockResolvedValue([]);

        render(
            <MemoryRouter>
                <ConnectionsPage />
            </MemoryRouter>
        );

        // Switch to Find People tab
        const findTab = screen.getByRole('tab', { name: /find people/i });
        await user.click(findTab);

        const searchInput = await screen.findByPlaceholderText('Search by email or username...');

        // Search
        fireEvent.change(searchInput, { target: { value: 'stranger' } });

        // Wait for debounce
        await act(async () => {
            await new Promise(r => setTimeout(r, 600));
        });

        await waitFor(() => {
            expect(screen.getByText('Stranger')).toBeInTheDocument();
        });

        // Click Connect
        const connectBtn = screen.getByRole('button', { name: /connect/i });
        await user.click(connectBtn);

        expect(mockSendRequest).toHaveBeenCalledWith('user-2', true);
    });

    it('accepts connection request', async () => {
        const user = userEvent.setup();
        const mockGetConnections = vi.mocked(api.getConnections);
        mockGetConnections.mockResolvedValue([
            {
                id: 'req-1',
                requester_id: '456',
                receiver_id: '123',
                status: 'pending',
                requester_username: 'Requester',
                requester_email: 'requester@example.com',
                receiver_username: 'Me',
                receiver_email: 'test@example.com'
            }
        ]);
        const mockRespond = vi.mocked(api.respondToConnectionRequest);
        mockRespond.mockResolvedValue({ success: true });

        render(
            <MemoryRouter>
                <ConnectionsPage />
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(screen.getByText('Requester')).toBeInTheDocument();
        });

        // Verify it's a request (has Accept/Decline)
        const acceptBtn = screen.getByRole('button', { name: /accept/i });
        await user.click(acceptBtn);

        expect(mockRespond).toHaveBeenCalledWith('req-1', 'accept');
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

        vi.useFakeTimers();
        fireEvent.change(searchInput, { target: { value: 'te' } });

        act(() => {
            vi.advanceTimersByTime(600);
        });

        expect(mockSearchUsers).not.toHaveBeenCalled();

        fireEvent.change(searchInput, { target: { value: 'test' } });
        expect(mockSearchUsers).not.toHaveBeenCalled();

        act(() => {
            vi.advanceTimersByTime(500);
        });

        vi.useRealTimers();

        await waitFor(() => {
            expect(mockSearchUsers).toHaveBeenCalledWith('test');
        });

        expect(await screen.findByText('Found User')).toBeInTheDocument();
    });
});
