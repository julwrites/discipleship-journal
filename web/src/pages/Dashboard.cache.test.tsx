import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import Dashboard from './Dashboard';
import { MemoryRouter } from 'react-router-dom';
import * as api from '@/services/api';

// Mocks (similar to Dashboard.test.tsx)
vi.mock('sonner', () => ({ toast: { success: vi.fn(), error: vi.fn() } }));
window.ResizeObserver = vi.fn().mockImplementation(() => ({
    observe: vi.fn(), unobserve: vi.fn(), disconnect: vi.fn()
}));
Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: vi.fn().mockImplementation(query => ({
        matches: false, media: query, onchange: null, addListener: vi.fn(), removeListener: vi.fn(), addEventListener: vi.fn(), removeEventListener: vi.fn(), dispatchEvent: vi.fn(),
    })),
});
vi.mock('@/services/api', () => ({
    fetchNotes: vi.fn(),
    syncUser: vi.fn(),
    askAIStream: vi.fn(),
    getTags: vi.fn().mockResolvedValue([]),
    getGroups: vi.fn().mockResolvedValue([]),
    getConnections: vi.fn().mockResolvedValue([]),
    getNote: vi.fn(),
    createNote: vi.fn(),
    deleteNote: vi.fn(),
    shareNote: vi.fn(),
    getOrCreateDirectGroup: vi.fn(),
}));
vi.mock('@/lib/firebase', () => ({ auth: { signOut: vi.fn() } }));
vi.mock('@/hooks/useAuth', () => ({
    useAuth: vi.fn().mockReturnValue({
        user: { uid: 'test-user-id', email: 'test@example.com' },
        loading: false
    })
}));

describe('Dashboard Caching', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        window.localStorage.clear();
    });

    it('shows cached notes immediately and updates from API', async () => {
        const userId = 'test-user-id';
        const cacheKey = `dashboard_notes_cache_${userId}`;
        const cachedNotes = [{ id: 'cached-1', title: 'Cached Note', updated_at: '2023-01-01', created_at: '2023-01-01', tags: [] }];

        // Seed cache
        window.localStorage.setItem(cacheKey, JSON.stringify({
            version: 'v1',
            timestamp: Date.now(),
            notes: cachedNotes,
            hasMore: false
        }));

        const freshNotes = [{ id: 'fresh-1', title: 'Fresh Note', updated_at: '2023-01-02', created_at: '2023-01-02', tags: [] }];
        const mockFetchNotes = vi.mocked(api.fetchNotes);

        // Delay API response
        mockFetchNotes.mockImplementation(async () => {
            await new Promise(resolve => setTimeout(resolve, 100));
            return { data: freshNotes, meta: { total_pages: 1 } };
        });

        render(
            <MemoryRouter>
                <Dashboard />
            </MemoryRouter>
        );

        // Verify cached note is shown immediately
        expect(screen.getByText('Cached Note')).toBeInTheDocument();
        // Fresh note should not be there yet
        expect(screen.queryByText('Fresh Note')).not.toBeInTheDocument();

        // Wait for API update
        await waitFor(() => {
            expect(screen.getByText('Fresh Note')).toBeInTheDocument();
        });

        // Verify cached note is gone (replaced)
        expect(screen.queryByText('Cached Note')).not.toBeInTheDocument();

        // Verify cache updated
        const updatedCache = JSON.parse(window.localStorage.getItem(cacheKey) || '{}');
        expect(updatedCache.notes[0].title).toBe('Fresh Note');
    });
});
