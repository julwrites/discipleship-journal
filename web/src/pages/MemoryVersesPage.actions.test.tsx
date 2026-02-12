import { render, screen, waitFor, within } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { VersePackDetail } from './MemoryVersesPage';
import { MemoryRouter, Routes, Route } from 'react-router-dom';
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
    getVersePacks: vi.fn(),
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

// Mock ResizeObserver
window.ResizeObserver = class {
  observe() {}
  unobserve() {}
  disconnect() {}
} as unknown as typeof ResizeObserver;

// Mock window.confirm
const mockConfirm = vi.spyOn(window, 'confirm');

describe('VersePackDetail Actions', () => {
    const mockUser = { id: 'user-1', settings: { bible_version: 'ESV' } };
    const mockPack = {
        id: 'pack-123',
        title: 'My Pack',
        description: 'Description',
        is_public: false,
        verse_count: 1,
        user_id: 'user-1'
    };
    const mockVerses = [
        {
            id: 'v1',
            verse_pack_id: 'pack-123',
            reference: 'John 3:16',
            version: 'ESV',
            version_source: 'original',
            tags: ['love']
        }
    ];

    beforeEach(() => {
        vi.clearAllMocks();
        mockConfirm.mockImplementation(() => true);
    });

    const setup = (packOverride = {}, versesOverride = mockVerses) => {
        const mockGetPackDetails = vi.mocked(api.getPackDetails);
        const mockSyncUser = vi.mocked(api.syncUser);

        mockSyncUser.mockResolvedValue(mockUser);
        mockGetPackDetails.mockResolvedValue({
            pack: { ...mockPack, ...packOverride },
            verses: versesOverride
        });

        const utils = render(
            <MemoryRouter initialEntries={['/memory-verses/pack-123']}>
                <Routes>
                    <Route path="/memory-verses/:id" element={<VersePackDetail />} />
                </Routes>
            </MemoryRouter>
        );
        return { ...utils };
    };

    it('adds a new verse', async () => {
        const user = userEvent.setup();
        const mockCreateVerseInPack = vi.mocked(api.createVerseInPack);
        mockCreateVerseInPack.mockResolvedValue({ id: 'v2' });

        setup();

        await waitFor(() => expect(screen.getByText('My Pack')).toBeInTheDocument());

        const addButton = screen.getByRole('button', { name: /Add Verse/i });
        await user.click(addButton);

        const dialog = await screen.findByRole('dialog');
        const refInput = within(dialog).getByPlaceholderText('e.g. John 3:16');

        await user.click(refInput);
        await user.paste('Rom 8:28'); // Try paste instead of type
        // await user.type(refInput, 'Rom 8:28');

        // await user.tab();

        expect(refInput).toHaveValue('Rom 8:28');

        const titleInput = within(dialog).getByPlaceholderText('e.g. God\'s Love');
        await user.type(titleInput, 'Good for those');

        const submitButton = within(dialog).getByRole('button', { name: 'Add Verse' });
        await user.click(submitButton);

        await waitFor(() => {
            expect(mockCreateVerseInPack).toHaveBeenCalledWith('pack-123', expect.objectContaining({
                reference: 'Romans 8:28',
                title: 'Good for those'
            }));
        });
    });

    it('deletes a verse', async () => {
        const user = userEvent.setup();
        const mockDeleteMemoryVerse = vi.mocked(api.deleteMemoryVerse);

        setup();

        await waitFor(() => expect(screen.getByText('John 3:16')).toBeInTheDocument());

        const verseText = screen.getByText('John 3:16');
        const card = verseText.closest('.border') as HTMLElement;

        const trashIcon = within(card).getByText((_, element) => {
            return element?.classList.contains('lucide-trash-2') || false;
        }, { selector: 'svg' });

        const deleteVerseBtn = trashIcon.closest('button');
        if (!deleteVerseBtn) throw new Error("Delete verse button not found");

        await user.click(deleteVerseBtn);

        expect(mockConfirm).toHaveBeenCalled();
        await waitFor(() => {
            expect(mockDeleteMemoryVerse).toHaveBeenCalledWith('v1');
        });
    });

    it('updates a verse in user pack', async () => {
        const user = userEvent.setup();
        const mockUpdateMemoryVerse = vi.mocked(api.updateMemoryVerse);

        setup();

        await waitFor(() => expect(screen.getByText('John 3:16')).toBeInTheDocument());

        const verseText = screen.getByText('John 3:16');
        const card = verseText.closest('.border') as HTMLElement;

        const pencilIcon = within(card).getByText((_, element) => {
            return element?.classList.contains('lucide-pencil') || false;
        }, { selector: 'svg' });

        const editVerseBtn = pencilIcon.closest('button');
        if (!editVerseBtn) throw new Error("Edit verse button not found");

        await user.click(editVerseBtn);

        const dialog = await screen.findByRole('dialog');
        const titleInput = within(dialog).getByPlaceholderText('e.g. God\'s Love');

        await user.clear(titleInput);
        await user.type(titleInput, 'New Title');

        const saveButton = within(dialog).getByRole('button', { name: 'Save Changes' });
        await user.click(saveButton);

        await waitFor(() => {
            expect(mockUpdateMemoryVerse).toHaveBeenCalledWith('v1', expect.objectContaining({
                title: 'New Title'
            }));
        });
    });

    it('sets verse preference in system pack', async () => {
        const user = userEvent.setup();
        const mockSetVersePreference = vi.mocked(api.setVersePreference);

        setup({ user_id: undefined }, [{ ...mockVerses[0], verse_pack_id: 'sys-1' }]);

        await waitFor(() => expect(screen.getByText('My Pack')).toBeInTheDocument());

        const verseText = screen.getByText('John 3:16');
        const card = verseText.closest('.border') as HTMLElement;

        const pencilIcon = within(card).getByText((_, element) => {
            return element?.classList.contains('lucide-pencil') || false;
        }, { selector: 'svg' });

        const editVerseBtn = pencilIcon.closest('button');
        if (!editVerseBtn) throw new Error("Edit verse button not found");

        await user.click(editVerseBtn);

        const dialog = await screen.findByRole('dialog');
        expect(within(dialog).getByText('Verse Preferences')).toBeInTheDocument();

        const saveButton = within(dialog).getByRole('button', { name: 'Save Changes' });
        await user.click(saveButton);

        await waitFor(() => {
            expect(mockSetVersePreference).toHaveBeenCalledWith('v1', 'ESV');
        });
    });

    it('deletes a pack', async () => {
        const user = userEvent.setup();
        const mockDeletePack = vi.mocked(api.deletePack);

        const { container } = setup();

        await waitFor(() => expect(screen.getByText('My Pack')).toBeInTheDocument());

        const allButtons = Array.from(container.querySelectorAll('button'));
        const packDeleteBtn = allButtons.find(b =>
            b.className.includes('bg-destructive') &&
            !b.closest('.border')
        );

        if (!packDeleteBtn) throw new Error("Could not find delete pack button");

        await user.click(packDeleteBtn);

        expect(mockConfirm).toHaveBeenCalled();
        await waitFor(() => {
            expect(mockDeletePack).toHaveBeenCalledWith('pack-123');
        });
    });

    it('views verse text', async () => {
        const user = userEvent.setup();
        const mockGetBiblePassage = vi.mocked(api.getBiblePassage);
        mockGetBiblePassage.mockResolvedValue({
            verse: 'For God so loved the world...',
            ref: 'John 3:16'
        });

        setup();

        await waitFor(() => expect(screen.getByText('John 3:16')).toBeInTheDocument());

        const verseText = screen.getByText('John 3:16');
        await user.click(verseText);

        await waitFor(() => {
            expect(mockGetBiblePassage).toHaveBeenCalledWith('John 3:16', 'ESV');
        });

        expect(await screen.findByText('For God so loved the world...')).toBeInTheDocument();
    });

    it('manages tags', async () => {
        const user = userEvent.setup();

        setup();

        await waitFor(() => expect(screen.getByText('My Pack')).toBeInTheDocument());

        const addButton = screen.getByRole('button', { name: /Add Verse/i });
        await user.click(addButton);

        const dialog = await screen.findByRole('dialog');

        const tagInput = within(dialog).getByPlaceholderText('Add tag...');
        const addTagBtn = within(dialog).getByRole('button', { name: 'Add' });

        await user.type(tagInput, 'hope');
        await user.click(addTagBtn);

        expect(within(dialog).getByText('hope')).toBeInTheDocument();

        const removeTagBtn = within(dialog).getByText('×');
        await user.click(removeTagBtn);

        expect(within(dialog).queryByText('hope')).not.toBeInTheDocument();
    });
});
