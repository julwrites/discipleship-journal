import { render, screen, fireEvent, waitFor, act } from '@testing-library/react';
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
window.ResizeObserver = vi.fn().mockImplementation(() => ({
    observe: vi.fn(),
    unobserve: vi.fn(),
    disconnect: vi.fn(),
}));

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: vi.fn().mockImplementation(query => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: vi.fn(), // deprecated
        removeListener: vi.fn(), // deprecated
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
        dispatchEvent: vi.fn(),
    })),
});

// Mock the API
vi.mock('@/services/api', () => ({
    fetchNotes: vi.fn(),
    syncUser: vi.fn(),
    askAIStream: vi.fn(),
    getTags: vi.fn().mockResolvedValue([]),
    getGroups: vi.fn().mockResolvedValue([]),
    getConnections: vi.fn().mockResolvedValue([]),
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
    });

    afterEach(() => {
        vi.useRealTimers();
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

        // Initial load
        await waitFor(() => {
            // NoteFilter object comparison
            expect(mockFetchNotes).toHaveBeenCalledWith(1, 20, expect.objectContaining({
                search: "",
                sortBy: "updated_at",
                sortOrder: "desc"
            }));
        });
    });

    it('reproduces double fetch issue on search from page 2', async () => {
        // Setup mocks
        const mockFetchNotes = vi.mocked(api.fetchNotes);

        // Initial load: returns enough notes to have a second page
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

        // Wait for initial load
        await waitFor(() => {
            expect(mockFetchNotes).toHaveBeenCalledWith(1, 20, expect.objectContaining({ search: "" }));
        });

        // Click Load More to go to page 2
        const loadMoreBtn = await screen.findByText('Load More');
        fireEvent.click(loadMoreBtn);

        // Wait for page 2 fetch
        await waitFor(() => {
            expect(mockFetchNotes).toHaveBeenCalledWith(2, 20, expect.objectContaining({ search: "" }));
        });

        // Clear mock calls to focus on search behavior
        mockFetchNotes.mockClear();

        // Now search
        vi.useFakeTimers();
        const searchInput = screen.getByPlaceholderText('Search notes...');
        fireEvent.change(searchInput, { target: { value: 'test' } });

        // Advance timer to trigger debounce
        act(() => {
            vi.advanceTimersByTime(500);
        });

        vi.useRealTimers();

        // Check calls
        await waitFor(() => {
            // We expect to potentially see a call with page 2 (bug) and then page 1 (correct)
            // Or just page 1 if lucky/fixed.
            // If the bug exists, we might see fetchNotes(2, 20, "test")
            const calls = mockFetchNotes.mock.calls;

            // Filter calls that have search "test"
            const searchCalls = calls.filter(call => {
                const filter = call[2] as api.NoteFilter;
                return filter && filter.search === 'test';
            });

            // If searchCalls has length > 1, and one of them is page 2, bug confirmed.
            if (searchCalls.length > 0) {
                 // Verify if we have the bad call
                 const hasBadCall = searchCalls.some(call => call[0] === 2);
                 expect(hasBadCall).toBe(false);
            }
            // We expect at least one call with page 1
            expect(mockFetchNotes).toHaveBeenCalledWith(1, 20, expect.objectContaining({ search: "test" }));
        });
    });
});
