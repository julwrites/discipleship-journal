import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { VersePackDetail } from './MemoryVersesPage';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
import * as api from '@/services/api';

// Mock sonner
vi.mock('sonner', () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
    }
}));

// Mock API
vi.mock('@/services/api', () => ({
    getPackDetails: vi.fn(),
    syncUser: vi.fn(),
    setVersePreference: vi.fn(),
    updateMemoryVerse: vi.fn(),
    deleteMemoryVerse: vi.fn(),
    createVerseInPack: vi.fn(),
    getBiblePassage: vi.fn(),
    removeVersePreference: vi.fn(),
    clonePack: vi.fn(),
    deletePack: vi.fn(),
}));

const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
    const actual = await vi.importActual('react-router-dom');
    return {
        ...actual,
        useNavigate: () => mockNavigate,
        useParams: () => ({ id: 'pack-123' }),
    };
});

describe('VersePackDetail', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('renders verses with correct version source indicators', async () => {
        const mockGetPackDetails = vi.mocked(api.getPackDetails);
        const mockSyncUser = vi.mocked(api.syncUser);

        mockSyncUser.mockResolvedValue({ id: 'user-1', settings: { bible_version: 'ESV' } });

        mockGetPackDetails.mockResolvedValue({
            pack: {
                id: 'pack-123',
                title: 'System Pack',
                description: 'A system pack',
                is_public: true,
                verse_count: 3,
                // user_id is undefined, so it's a system pack
            },
            verses: [
                {
                    id: 'v1',
                    verse_pack_id: 'pack-123',
                    reference: 'John 3:16',
                    version: 'NIV',
                    version_source: 'override',
                    tags: []
                },
                {
                    id: 'v2',
                    verse_pack_id: 'pack-123',
                    reference: 'Gen 1:1',
                    version: 'ESV',
                    version_source: 'user_default',
                    tags: []
                },
                {
                    id: 'v3',
                    verse_pack_id: 'pack-123',
                    reference: 'Ps 23:1',
                    version: 'KJV',
                    version_source: 'original',
                    tags: []
                }
            ]
        });

        render(
            <MemoryRouter initialEntries={['/memory-verses/pack-123']}>
                <Routes>
                    <Route path="/memory-verses/:id" element={<VersePackDetail />} />
                </Routes>
            </MemoryRouter>
        );

        await waitFor(() => {
            expect(screen.getByText('John 3:16')).toBeInTheDocument();
        });

        // Check for indicators (by text/version)
        const nivBadge = screen.getByText('NIV');
        expect(nivBadge).toHaveClass('text-blue-800'); // Check part of the class for override
        expect(nivBadge).toHaveAttribute('title', 'Your custom preference');

        const esvBadge = screen.getByText('ESV');
        expect(esvBadge).toHaveClass('text-green-800'); // Check part of the class for user_default
        expect(esvBadge).toHaveAttribute('title', 'Your default setting');

        const kjvBadge = screen.getByText('KJV');
        expect(kjvBadge).toHaveClass('text-muted-foreground'); // Check part of the class for original
    });

    it('handles bulk apply default version', async () => {
        const mockGetPackDetails = vi.mocked(api.getPackDetails);
        const mockSyncUser = vi.mocked(api.syncUser);
        const mockSetVersePreference = vi.mocked(api.setVersePreference);

        // User default is KJV
        mockSyncUser.mockResolvedValue({ id: 'user-1', settings: { bible_version: 'KJV' } });

        mockGetPackDetails.mockResolvedValue({
            pack: {
                id: 'pack-123',
                title: 'System Pack',
                is_public: true,
                verse_count: 2,
            },
            verses: [
                { id: 'v1', verse_pack_id: 'pack-123', reference: 'John 3:16', version: 'NIV', version_source: 'override', tags: [] },
                { id: 'v2', verse_pack_id: 'pack-123', reference: 'Gen 1:1', version: 'ESV', version_source: 'original', tags: [] },
            ]
        });

        // Mock confirm
        vi.spyOn(window, 'confirm').mockImplementation(() => true);

        render(
            <MemoryRouter initialEntries={['/memory-verses/pack-123']}>
                <Routes>
                    <Route path="/memory-verses/:id" element={<VersePackDetail />} />
                </Routes>
            </MemoryRouter>
        );

        // Wait for button to appear
        const bulkButton = await screen.findByText('Set all to KJV');
        expect(bulkButton).toBeInTheDocument();

        // Click it
        fireEvent.click(bulkButton);

        await waitFor(() => {
            expect(mockSetVersePreference).toHaveBeenCalledTimes(2);
            expect(mockSetVersePreference).toHaveBeenCalledWith('v1', 'KJV');
            expect(mockSetVersePreference).toHaveBeenCalledWith('v2', 'KJV');
        });
    });
});
