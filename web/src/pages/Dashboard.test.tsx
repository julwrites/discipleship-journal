import { render, screen, waitFor, act, within, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import Dashboard from './Dashboard';
import { MemoryRouter } from 'react-router-dom';
import * as api from '@/services/api';

// Mock sonner
vi.mock('sonner', () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
    }
}));

// Mock ResizeObserver
window.ResizeObserver = class ResizeObserver {
    observe = vi.fn();
    unobserve = vi.fn();
    disconnect = vi.fn();
} as unknown as typeof ResizeObserver;

// Mock matchMedia
Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: vi.fn().mockImplementation(query => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: vi.fn(),
        removeListener: vi.fn(),
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
        dispatchEvent: vi.fn(),
    })),
});

// Mock pointer capture methods for Radix UI
window.HTMLElement.prototype.hasPointerCapture = vi.fn();
window.HTMLElement.prototype.setPointerCapture = vi.fn();
window.HTMLElement.prototype.releasePointerCapture = vi.fn();
window.HTMLElement.prototype.scrollIntoView = vi.fn();

// Mock the API
vi.mock('@/services/api', () => ({
    fetchNotes: vi.fn(),
    syncUser: vi.fn(),
    askAIStream: vi.fn(),
    getTags: vi.fn().mockResolvedValue([]),
    getGroups: vi.fn().mockResolvedValue([]),
    getConnections: vi.fn().mockResolvedValue([]),
    deleteNote: vi.fn(),
    shareNote: vi.fn(),
    getOrCreateDirectGroup: vi.fn(),
    getNote: vi.fn(),
}));

// Mock cache
vi.mock('@/services/cache', () => ({
    getCachedNotes: vi.fn().mockReturnValue(null),
    setCachedNotes: vi.fn(),
}));

// Mock Firebase
vi.mock('@/lib/firebase', () => ({
    auth: {
        signOut: vi.fn(),
    }
}));

// Mock useAuth
vi.mock('@/hooks/useAuth', () => ({
    useAuth: vi.fn().mockReturnValue({
        user: { uid: 'test-user-id', email: 'test@example.com' },
        loading: false
    })
}));

describe('Dashboard', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        vi.useRealTimers();
        document.body.style.pointerEvents = 'auto';
        document.body.innerHTML = '';
    });

    afterEach(() => {
        vi.useRealTimers();
        document.body.style.pointerEvents = 'auto';
    });

    it('fetches notes with correct filter params', async () => {
        const mockFetchNotes = vi.mocked(api.fetchNotes);

        mockFetchNotes.mockResolvedValue({
            data: [],
            meta: { total_pages: 1 }
        });

        render(
            <MemoryRouter>
                <Dashboard />
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(mockFetchNotes).toHaveBeenCalledWith(1, 20, expect.objectContaining({
                search: "",
                sortBy: "updated_at",
                sortOrder: "desc"
            }));
        });
    });

    it('reproduces double fetch issue on search from page 2', async () => {
        // Use fireEvent as in original test for simplicity
        const mockFetchNotes = vi.mocked(api.fetchNotes);

        mockFetchNotes.mockImplementation(async (page = 1, limit = 20) => {
            const offset = (page - 1) * limit;
            return {
                data: Array(limit).fill(null).map((_, i) => ({
                    id: `note-${offset + i}`,
                    title: `Note ${offset + i}`,
                    updated_at: '2023-01-01'
                })),
                meta: { total_pages: 2 }
            };
        });

        render(
            <MemoryRouter>
                <Dashboard />
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(mockFetchNotes).toHaveBeenCalledWith(1, 20, expect.objectContaining({ search: "" }));
        });

        const loadMoreBtn = await screen.findByText('Load More');
        fireEvent.click(loadMoreBtn);

        await waitFor(() => {
            expect(mockFetchNotes).toHaveBeenCalledWith(2, 20, expect.objectContaining({ search: "" }));
        });

        mockFetchNotes.mockClear();

        vi.useFakeTimers();
        const searchInput = screen.getByPlaceholderText('Search notes...');
        fireEvent.change(searchInput, { target: { value: 'test' } });

        act(() => {
            vi.advanceTimersByTime(500);
        });

        vi.useRealTimers();

        await waitFor(() => {
            expect(mockFetchNotes).toHaveBeenCalledWith(1, 20, expect.objectContaining({ search: "test" }));
        });
    });

    it('deletes a note', async () => {
        const user = userEvent.setup();
        const mockFetchNotes = vi.mocked(api.fetchNotes);
        const mockDeleteNote = vi.mocked(api.deleteNote);
        const note = { id: '1', title: 'Note to Delete', content: 'Content', updated_at: '2023-01-01' };

        mockFetchNotes.mockResolvedValue({ data: [note], meta: { total_pages: 1 } });
        mockDeleteNote.mockResolvedValue(undefined);

        render(
            <MemoryRouter>
                <Dashboard />
            </MemoryRouter>
        );

        const noteTitle = await screen.findByText('Note to Delete');
        const card = noteTitle.closest('div[role="link"]') as HTMLElement;
        expect(card).toBeInTheDocument();
        if (!card) throw new Error("Card not found");

        const deleteBtn = within(card).getByTitle('Delete');
        await user.click(deleteBtn);

        await screen.findByText('Delete Note');
        const dialog = screen.getByRole('dialog');
        const confirmBtn = within(dialog).getByText('Delete', { selector: 'button' });

        await user.click(confirmBtn);

        await waitFor(() => {
            expect(mockDeleteNote).toHaveBeenCalledWith('1');
        });

        await waitFor(() => {
             expect(screen.queryByText('Note to Delete')).not.toBeInTheDocument();
        });
    });

    it('shares a note', async () => {
        const user = userEvent.setup();
        const mockFetchNotes = vi.mocked(api.fetchNotes);
        const mockShareNote = vi.mocked(api.shareNote);
        const mockGetGroups = vi.mocked(api.getGroups);
        const mockGetConnections = vi.mocked(api.getConnections);

        const note = { id: '1', title: 'Note to Share', content: 'Content', updated_at: '2023-01-01' };
        mockFetchNotes.mockResolvedValue({ data: [note], meta: { total_pages: 1 } });
        mockGetGroups.mockResolvedValue([{ id: 'g1', name: 'Study Group', type: 'normal' }]);
        mockGetConnections.mockResolvedValue([]);
        mockShareNote.mockResolvedValue(undefined);

        render(
            <MemoryRouter>
                <Dashboard />
            </MemoryRouter>
        );

        const noteTitle = await screen.findByText('Note to Share');
        const card = noteTitle.closest('div[role="link"]') as HTMLElement;
        if (!card) throw new Error("Card not found");

        const shareBtn = within(card).getByTitle('Share');
        // Use fireEvent to avoid potential userEvent issues with overlapping async ops
        fireEvent.click(shareBtn);

        const dialog = await screen.findByRole('dialog');
        // Wait for dialog content. We check for the submit button or title.
        // Use regex for case insensitivity or check exact string.
        await within(dialog).findByRole('heading', { name: /share note/i });

        const trigger = screen.getByText('Select a Group...');
        // Use fireEvent for interactions inside dialog to avoid pointer-events issues
        fireEvent.click(trigger);

        const groupOption = await screen.findByText('Study Group');
        fireEvent.click(groupOption);

        const commentInput = screen.getByPlaceholderText('Add a comment (optional)...');
        await user.type(commentInput, 'Check this out');

        // Dialog is already found
        const confirmBtn = within(dialog).getByText('Share Note', { selector: 'button' });
        await user.click(confirmBtn);

        await waitFor(() => {
             expect(mockShareNote).toHaveBeenCalledWith('g1', '1', 'Check this out');
        });
    });

    it('asks AI about a note', async () => {
        const user = userEvent.setup();
        const mockFetchNotes = vi.mocked(api.fetchNotes);
        const mockGetNote = vi.mocked(api.getNote);
        const mockAskAIStream = vi.mocked(api.askAIStream);

        const note = { id: '1', title: 'Bible Note', content: 'Genesis 1:1', updated_at: '2023-01-01' };
        mockFetchNotes.mockResolvedValue({ data: [note], meta: { total_pages: 1 } });
        mockGetNote.mockResolvedValue(note);

        mockAskAIStream.mockImplementation(async (_content, _prompt, _history, callbacks) => {
            if (callbacks?.onStart) callbacks.onStart("note-id");
            if (callbacks?.onChunk) callbacks.onChunk('God created');
            if (callbacks?.onDone) callbacks.onDone();
            return Promise.resolve();
        });

        render(
            <MemoryRouter>
                <Dashboard />
            </MemoryRouter>
        );

        const noteTitle = await screen.findByText('Bible Note');
        const card = noteTitle.closest('div[role="link"]') as HTMLElement;
        if (!card) throw new Error("Card not found");

        const aiBtn = within(card).getByTitle('Ask AI');
        await user.click(aiBtn);

        await screen.findByText('Ask AI');

        await waitFor(() => {
             expect(screen.queryByText('Loading note content...')).not.toBeInTheDocument();
        });

        const promptInput = screen.getByPlaceholderText('Ask a question...');
        await user.type(promptInput, 'What happened?');

        const dialog = screen.getByRole('dialog');
        const askBtn = within(dialog).getByText('Ask', { selector: 'button' });
        await user.click(askBtn);

        await screen.findByText('God created');
    });
});
