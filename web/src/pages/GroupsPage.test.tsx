import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import GroupsPage from './GroupsPage';
import { useAuth } from '@/hooks/useAuth';

// Mock the useAuth hook
vi.mock('@/hooks/useAuth', () => ({
    useAuth: vi.fn(),
}));

// Mock UI components if necessary (Shadcn components sometimes need mocking if they use complex browser APIs, but usually they are fine)
// We will rely on JSDOM.

// Mock Fetch
global.fetch = vi.fn();

describe('GroupsPage', () => {
    const mockUser = {
        getIdToken: vi.fn().mockResolvedValue('mock-token'),
        uid: 'user-123',
    };

    beforeEach(() => {
        vi.clearAllMocks();
        (useAuth as any).mockReturnValue({ user: mockUser });
    });

    it('renders and fetches groups', async () => {
        (global.fetch as any).mockResolvedValueOnce({
            ok: true,
            json: async () => ([
                { id: '1', name: 'Bible Study', description: 'Weekly study', role: 'member' }
            ]),
        });

        render(<GroupsPage />);

        expect(screen.getByText('Groups')).toBeInTheDocument();
        expect(screen.getByText('Create Group')).toBeInTheDocument();

        // Wait for fetch to be called
        await waitFor(() => {
            expect(global.fetch).toHaveBeenCalledWith(
                expect.stringContaining('/api/groups'),
                expect.objectContaining({
                    headers: expect.objectContaining({ Authorization: 'Bearer mock-token' })
                })
            );
        });

        // Check if group is displayed
        expect(await screen.findByText('Bible Study')).toBeInTheDocument();
        expect(screen.getByText('Weekly study')).toBeInTheDocument();
    });

    it('opens create dialog and submits form', async () => {
        // Initial fetch returns empty
        (global.fetch as any).mockResolvedValueOnce({
            ok: true,
            json: async () => ([]),
        });

        render(<GroupsPage />);

        // Open Dialog
        fireEvent.click(screen.getByText('Create Group'));
        expect(screen.getByText('Create a New Group')).toBeInTheDocument();

        // Fill Form
        fireEvent.change(screen.getByPlaceholderText('Group Name'), { target: { value: 'New Group' } });
        fireEvent.change(screen.getByPlaceholderText('Description'), { target: { value: 'New Description' } });

        // Mock Create Response
        (global.fetch as any).mockResolvedValueOnce({
            ok: true,
            json: async () => ({ id: '2' }),
        });
        // Mock Refresh Fetch
        (global.fetch as any).mockResolvedValueOnce({
            ok: true,
            json: async () => ([
                 { id: '2', name: 'New Group', description: 'New Description', role: 'admin' }
            ]),
        });

        // Submit
        const createBtns = screen.getAllByText('Create');
        // The first one is the trigger, the second one is inside the dialog (usually last in DOM)
        fireEvent.click(createBtns[createBtns.length - 1]);

        await waitFor(() => {
            expect(global.fetch).toHaveBeenCalledWith(
                expect.stringContaining('/api/groups'),
                expect.objectContaining({
                    method: 'POST',
                    body: JSON.stringify({ name: 'New Group', description: 'New Description' })
                })
            );
        });
    });

    it('searches for groups', async () => {
        const user = userEvent.setup();
         // Initial fetch returns empty
         (global.fetch as any).mockResolvedValueOnce({
            ok: true,
            json: async () => ([]),
        });

        render(<GroupsPage />);

        // Switch to Find Groups tab
        const findGroupsTrigger = screen.getByText('Find Groups');
        await user.click(findGroupsTrigger);

        // Wait for tab content to be visible
        // We look for the input which should be visible when the tab is active
        const searchInput = await screen.findByPlaceholderText('Search groups...');

        // Type search
        await user.type(searchInput, 'Prayer');

        // Mock Search Response
        (global.fetch as any).mockResolvedValueOnce({
            ok: true,
            json: async () => ([
                { id: '3', name: 'Prayer Warriors', description: 'Praying together', role: '' } // No role = not joined
            ]),
        });

        // Click Search
        // Since there are multiple "Search" buttons (one in My Groups add member, one in Find Groups), we need to be careful
        // Actually, the "Search" button in My Groups is inside a Dialog that is closed, so maybe it's not visible?
        // But the dialog content might be rendered? No, DialogContent is usually not rendered when closed.
        // But in `GroupsPage`, `searchUsers` button has text "Search", and `handleSearch` button has text "Search".
        // They are in different tabs.
        // My Groups tab is now hidden, Find Groups tab is visible.
        // screen.getByText('Search') might fail if multiple exist.
        // Let's scope it or use getAllByText.
        // The visible one should be the one we want.

        const searchBtn = screen.getByRole('button', { name: 'Search' });
        await user.click(searchBtn);

        await waitFor(() => {
            expect(global.fetch).toHaveBeenCalledWith(
                expect.stringContaining('/api/groups/search?q=Prayer'),
                expect.any(Object)
            );
        });

        expect(await screen.findByText('Prayer Warriors')).toBeInTheDocument();
        expect(screen.getByText('Join Group')).toBeInTheDocument();
    });
});
