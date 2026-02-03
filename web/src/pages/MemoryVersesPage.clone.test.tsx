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

describe('VersePackDetail - Clone Functionality', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('prompts for version preference when cloning and respects user choice', async () => {
        const mockGetPackDetails = vi.mocked(api.getPackDetails);
        const mockSyncUser = vi.mocked(api.syncUser);
        const mockClonePack = vi.mocked(api.clonePack);

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
            ]
        });

        render(
            <MemoryRouter initialEntries={['/memory-verses/pack-123']}>
                <Routes>
                    <Route path="/memory-verses/:id" element={<VersePackDetail />} />
                </Routes>
            </MemoryRouter>
        );

        // Wait for "Save to My Packs" button
        const saveButton = await screen.findByText('Save to My Packs');
        expect(saveButton).toBeInTheDocument();

        // Click "Save to My Packs"
        fireEvent.click(saveButton);

        // Wait for dialog to open
        const dialogTitle = await screen.findByText('Save Pack');
        expect(dialogTitle).toBeInTheDocument();

        // Check if checkbox exists and is checked by default
        const checkboxLabel = screen.getByText(/Apply my default version \(KJV\) to all verses/);
        expect(checkboxLabel).toBeInTheDocument();

        // Find the checkbox input associated with the label
        const checkbox = screen.getByRole('checkbox');
        expect(checkbox).toBeChecked(); // Default is true

        // Click Save (defaults to true)
        const saveConfirmButton = screen.getByRole('button', { name: 'Save' });
        fireEvent.click(saveConfirmButton);

        await waitFor(() => {
            expect(mockClonePack).toHaveBeenCalledWith('pack-123', undefined, true);
        });

        // Test unchecking
        // Re-open dialog (mock implementation might need reset or we rely on component state)
        // Since we clicked save, the dialog closed. Let's open it again.
        fireEvent.click(saveButton);
        await screen.findByText('Save Pack');

        const checkbox2 = screen.getByRole('checkbox');
        fireEvent.click(checkbox2); // Uncheck
        expect(checkbox2).not.toBeChecked();

        const saveConfirmButton2 = screen.getByRole('button', { name: 'Save' });
        fireEvent.click(saveConfirmButton2);

        await waitFor(() => {
            expect(mockClonePack).toHaveBeenCalledWith('pack-123', undefined, false);
        });
    });

    it('does not show checkbox if user has no default version', async () => {
         const mockGetPackDetails = vi.mocked(api.getPackDetails);
        const mockSyncUser = vi.mocked(api.syncUser);
        const mockClonePack = vi.mocked(api.clonePack);

        // User has NO default version
        mockSyncUser.mockResolvedValue({ id: 'user-1', settings: {} });

        mockGetPackDetails.mockResolvedValue({
            pack: {
                id: 'pack-123',
                title: 'System Pack',
                is_public: true,
                verse_count: 2,
            },
            verses: []
        });

        render(
            <MemoryRouter initialEntries={['/memory-verses/pack-123']}>
                <Routes>
                    <Route path="/memory-verses/:id" element={<VersePackDetail />} />
                </Routes>
            </MemoryRouter>
        );

         // Wait for "Save to My Packs" button
        const saveButton = await screen.findByText('Save to My Packs');
        fireEvent.click(saveButton);

        // Wait for dialog
        await screen.findByText('Save Pack');

        // Checkbox should NOT be there
        const checkbox = screen.queryByRole('checkbox');
        expect(checkbox).not.toBeInTheDocument();

        const saveConfirmButton = screen.getByRole('button', { name: 'Save' });
        fireEvent.click(saveConfirmButton);

        await waitFor(() => {
             // Should default to true anyway in the call, but UI option is hidden
            expect(mockClonePack).toHaveBeenCalledWith('pack-123', undefined, true);
        });
    });
});
