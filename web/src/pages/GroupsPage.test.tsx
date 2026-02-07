import { render, screen, waitFor, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import GroupsPage from './GroupsPage';
import { MemoryRouter } from 'react-router-dom';
import * as api from '@/services/api';

// Mock sonner
vi.mock('sonner', () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
    }
}));

// Mock the API
vi.mock('@/services/api', () => ({
    getGroups: vi.fn(),
    searchGroups: vi.fn(),
    createGroup: vi.fn(),
    joinGroup: vi.fn(),
    leaveGroup: vi.fn(),
    getGroupMembers: vi.fn(),
    getGroupShares: vi.fn(),
    getSharedItem: vi.fn(),
    getConnections: vi.fn(),
    searchUsers: vi.fn(),
    addGroupMember: vi.fn(),
    removeGroupMember: vi.fn(),
}));

// Mock useAuth
const mockUser = { uid: '123', email: 'test@example.com' };
vi.mock('@/hooks/useAuth', () => ({
    useAuth: () => ({ user: mockUser }),
}));

describe('GroupsPage', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        vi.useRealTimers();
    });

    afterEach(() => {
        vi.useRealTimers();
    });

    it('debounces search in Find Groups tab', async () => {
        const user = userEvent.setup();
        const mockSearchGroups = vi.mocked(api.searchGroups);
        mockSearchGroups.mockResolvedValue([
            { id: 'group-1', name: 'Found Group', description: 'A found group' }
        ]);

        const mockGetGroups = vi.mocked(api.getGroups);
        mockGetGroups.mockResolvedValue([]);

        render(
            <MemoryRouter>
                <GroupsPage />
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(mockGetGroups).toHaveBeenCalled();
        });

        // Switch to Find Groups tab
        const findTab = screen.getByRole('tab', { name: /find groups/i });
        await user.click(findTab);

        const searchInput = await screen.findByPlaceholderText('Search groups...');

        const { fireEvent } = await import('@testing-library/react');
        vi.useFakeTimers();

        // Type "te" (length < 3)
        fireEvent.change(searchInput, { target: { value: 'te' } });

        act(() => {
            vi.advanceTimersByTime(600);
        });

        expect(mockSearchGroups).not.toHaveBeenCalled();

        // Type "test" (length >= 3)
        fireEvent.change(searchInput, { target: { value: 'test' } });

        expect(mockSearchGroups).not.toHaveBeenCalled();

        act(() => {
            vi.advanceTimersByTime(500);
        });

        vi.useRealTimers();

        await waitFor(() => {
            expect(mockSearchGroups).toHaveBeenCalledWith('test');
        });

        expect(await screen.findByText('Found Group')).toBeInTheDocument();
    });

    it('lists connections in Add Member dialog', async () => {
        const user = userEvent.setup();
        const mockGetGroups = vi.mocked(api.getGroups);
        mockGetGroups.mockResolvedValue([
            { id: 'group-1', name: 'My Group', description: 'My group', role: 'admin' }
        ]);
        const mockGetGroupMembers = vi.mocked(api.getGroupMembers);
        mockGetGroupMembers.mockResolvedValue([]);
        const mockGetGroupShares = vi.mocked(api.getGroupShares);
        mockGetGroupShares.mockResolvedValue([]);

        const mockGetConnections = vi.mocked(api.getConnections);
        mockGetConnections.mockResolvedValue([
            { id: 'conn-1', requester_id: '123', receiver_id: 'user-2', status: 'accepted', receiver_email: 'friend@example.com', requester_email: 'test@example.com' }
        ]);

        render(
            <MemoryRouter>
                <GroupsPage />
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(screen.getByText('My Group')).toBeInTheDocument();
        });

        // Expand group details
        await user.click(screen.getByText('My Group'));

        await waitFor(() => {
            expect(mockGetGroupMembers).toHaveBeenCalledWith('group-1');
        });

        // Switch to Members tab
        const membersTab = screen.getByRole('tab', { name: /members/i });
        await user.click(membersTab);

        // Click Add Member
        const addMemberBtn = await screen.findByText(/add member/i);
        await user.click(addMemberBtn);

        // Wait for connections to load
        await waitFor(() => {
            expect(mockGetConnections).toHaveBeenCalled();
        });

        expect(await screen.findByText('friend@example.com')).toBeInTheDocument();
    });

    it('creates a new group', async () => {
        const user = userEvent.setup();
        const mockGetGroups = vi.mocked(api.getGroups);
        mockGetGroups.mockResolvedValue([]);
        const mockCreateGroup = vi.mocked(api.createGroup);
        mockCreateGroup.mockResolvedValue({ id: 'new-group-1' });

        render(
            <MemoryRouter>
                <GroupsPage />
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(mockGetGroups).toHaveBeenCalled();
        });

        // Click Create Group button
        const createBtn = screen.getByRole('button', { name: /create group/i });
        await user.click(createBtn);

        // Fill form
        const nameInput = await screen.findByPlaceholderText(/group name/i);
        const descInput = await screen.findByPlaceholderText(/description/i);

        await user.type(nameInput, 'New Group Name');
        await user.type(descInput, 'New Group Description');

        // Submit
        const submitBtn = screen.getByRole('button', { name: /create/i });
        await user.click(submitBtn);

        await waitFor(() => {
            expect(mockCreateGroup).toHaveBeenCalledWith({
                name: 'New Group Name',
                description: 'New Group Description',
            });
        });

        // Should refresh groups
        expect(mockGetGroups).toHaveBeenCalledTimes(2);
    });
});
