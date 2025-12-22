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
    getSharedNote: vi.fn(),
    searchUsers: vi.fn(),
    addGroupMember: vi.fn(),
    removeGroupMember: vi.fn(),
}));

// Mock useAuth
vi.mock('@/hooks/useAuth', () => ({
    useAuth: () => ({ user: { uid: '123', email: 'test@example.com' } }),
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

    it('debounces user search in Add Member dialog', async () => {
        const user = userEvent.setup();
        const mockGetGroups = vi.mocked(api.getGroups);
        mockGetGroups.mockResolvedValue([
            { id: 'group-1', name: 'My Group', description: 'My group', role: 'admin' }
        ]);
        const mockGetGroupMembers = vi.mocked(api.getGroupMembers);
        mockGetGroupMembers.mockResolvedValue([]);
        const mockGetGroupShares = vi.mocked(api.getGroupShares);
        mockGetGroupShares.mockResolvedValue([]);

        const mockSearchUsers = vi.mocked(api.searchUsers);
        mockSearchUsers.mockResolvedValue([
            { id: 'user-1', email: 'newmember@example.com', display_name: 'New Member' }
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

        // Wait for dialog
        const searchInput = await screen.findByPlaceholderText(/search by email, name, or username/i);

        const { fireEvent } = await import('@testing-library/react');
        vi.useFakeTimers();

        fireEvent.change(searchInput, { target: { value: 'user' } });

        expect(mockSearchUsers).not.toHaveBeenCalled();

        act(() => {
            vi.advanceTimersByTime(500);
        });

        vi.useRealTimers();

        await waitFor(() => {
            expect(mockSearchUsers).toHaveBeenCalledWith('user');
        });

        expect(await screen.findByText('New Member')).toBeInTheDocument();
    });
});
