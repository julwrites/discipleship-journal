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

// Mock the API
vi.mock('@/services/api', () => ({
    fetchNotes: vi.fn(),
    syncUser: vi.fn(),
}));

// Mock Firebase
vi.mock('@/lib/firebase', () => ({
    auth: {
        signOut: vi.fn(),
    }
}));

describe('Dashboard', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        vi.useRealTimers();
    });

    afterEach(() => {
        vi.useRealTimers();
    });

    it('reproduces double fetch issue on search from page 2', async () => {
        // Setup mocks
        const mockFetchNotes = vi.mocked(api.fetchNotes);

        // Initial load: returns enough notes to have a second page
        mockFetchNotes.mockResolvedValue({
            data: Array(20).fill(null).map((_, i) => ({ id: `note-${i}`, title: `Note ${i}`, updated_at: '2023-01-01' })),
            meta: { total_pages: 2 }
        });

        render(
            <MemoryRouter>
                <Dashboard />
            </MemoryRouter>
        );

        // Wait for initial load
        await waitFor(() => {
            expect(mockFetchNotes).toHaveBeenCalledWith(1, 20, "");
        });

        // Click Load More to go to page 2
        // We need to make sure we render the Load More button.
        // The component checks: hasMore && notes.length > 0
        // hasMore is true initially.
        // After first fetch, hasMore = 1 < 2 = true.

        const loadMoreBtn = await screen.findByText('Load More');
        fireEvent.click(loadMoreBtn);

        // Wait for page 2 fetch
        await waitFor(() => {
            expect(mockFetchNotes).toHaveBeenCalledWith(2, 20, "");
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
            // Filter calls that have "test" as 3rd arg
            const searchCalls = calls.filter(call => call[2] === 'test');

            // If searchCalls has length > 1, and one of them is page 2, bug confirmed.
            if (searchCalls.length > 0) {
                 // Verify if we have the bad call
                 const hasBadCall = searchCalls.some(call => call[0] === 2);
                 expect(hasBadCall).toBe(false);
            }
            // We expect at least one call with page 1
            expect(mockFetchNotes).toHaveBeenCalledWith(1, 20, "test");
        });
    });
});
