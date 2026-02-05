import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import MemoryVersesPage from './MemoryVersesPage';
import { MemoryRouter } from 'react-router-dom';
import * as api from '@/services/api';
import userEvent from '@testing-library/user-event';

// Mock sonner
vi.mock('sonner', () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
    }
}));

// Mock API
vi.mock('@/services/api', () => ({
    getVersePacks: vi.fn(),
    createVersePack: vi.fn(),
    // Add other methods to prevent "not a function" errors if they are called implicitly
    getPackDetails: vi.fn(),
    syncUser: vi.fn(),
}));

const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
    const actual = await vi.importActual('react-router-dom');
    return {
        ...actual,
        useNavigate: () => mockNavigate,
    };
});

describe('MemoryVersesPage - List View', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('renders loading state initially', async () => {
        const mockGetVersePacks = vi.mocked(api.getVersePacks);
        // Return a promise that doesn't resolve immediately to test loading state
        mockGetVersePacks.mockImplementation(() => new Promise(() => {}));

        render(
            <MemoryRouter>
                <MemoryVersesPage />
            </MemoryRouter>
        );

        // "My Packs" is default tab. PackList should be loading.
        expect(screen.getByText('Loading...')).toBeInTheDocument();
    });

    it('renders empty state when no packs', async () => {
        const mockGetVersePacks = vi.mocked(api.getVersePacks);
        mockGetVersePacks.mockResolvedValue({ data: [] });

        render(
            <MemoryRouter>
                <MemoryVersesPage />
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(screen.getByText('No packs found.')).toBeInTheDocument();
        });
    });

    it('renders a list of user packs', async () => {
        const mockGetVersePacks = vi.mocked(api.getVersePacks);
        mockGetVersePacks.mockResolvedValue({
            data: [
                { id: '1', title: 'My Favorites', verse_count: 5, user_id: 'u1' },
                { id: '2', title: 'Study Verses', verse_count: 2, user_id: 'u1' }
            ]
        });

        render(
            <MemoryRouter>
                <MemoryVersesPage />
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(screen.getByText('My Favorites')).toBeInTheDocument();
            expect(screen.getByText('Study Verses')).toBeInTheDocument();
            expect(screen.getByText('5 verses')).toBeInTheDocument();
            expect(screen.getByText('2 verses')).toBeInTheDocument();
        });

        // Click on a pack navigates
        fireEvent.click(screen.getByText('My Favorites'));
        expect(mockNavigate).toHaveBeenCalledWith('/memory-verses/1');
    });

    it('switches to System Packs tab', async () => {
        const user = userEvent.setup();
        const mockGetVersePacks = vi.mocked(api.getVersePacks);

        // Setup returns based on type argument
        mockGetVersePacks.mockImplementation(async (type) => {
            if (type === 'system') {
                return { data: [{ id: 's1', title: 'Top 100 Verses', verse_count: 100, is_public: true }] };
            }
            return { data: [] }; // User packs empty
        });

        render(
            <MemoryRouter>
                <MemoryVersesPage />
            </MemoryRouter>
        );

        // Default is My Packs (empty in this test)
        await waitFor(() => {
            expect(screen.getByText('No packs found.')).toBeInTheDocument();
        });

        // Click System Packs tab
        await user.click(screen.getByText('System Packs'));

        await waitFor(() => {
            expect(screen.getByText('Top 100 Verses')).toBeInTheDocument();
        });

        expect(mockGetVersePacks).toHaveBeenCalledWith('system');
    });

    it('creates a new pack', async () => {
        const user = userEvent.setup();
        const mockGetVersePacks = vi.mocked(api.getVersePacks);
        const mockCreateVersePack = vi.mocked(api.createVersePack);

        mockGetVersePacks.mockResolvedValue({ data: [] });
        mockCreateVersePack.mockResolvedValue({ id: 'new-1', title: 'New Pack' });

        render(
            <MemoryRouter>
                <MemoryVersesPage />
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(screen.getByText('Create Pack')).toBeInTheDocument();
        });

        // Open Dialog
        await user.click(screen.getByText('Create Pack'));
        expect(screen.getByText('Create New Pack')).toBeInTheDocument();

        // Fill Form
        const input = screen.getByPlaceholderText('e.g. My Favorites');
        await user.type(input, 'New Pack');

        // Submit
        await user.click(screen.getByRole('button', { name: 'Create' }));

        await waitFor(() => {
            expect(mockCreateVersePack).toHaveBeenCalledWith('New Pack');
        });

        // Should reload packs (mockGetVersePacks called again)
        expect(mockGetVersePacks).toHaveBeenCalledTimes(2); // Initial + After create
    });
});
