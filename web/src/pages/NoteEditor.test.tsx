import { render, screen, fireEvent, waitFor, act, within } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import NoteEditor from './NoteEditor';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import * as api from '@/services/api';
import { toast } from 'sonner';

// Mock sonner
vi.mock('sonner', () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
    }
}));

// Mock RichTextEditor to avoid complex Tiptap interaction in integration tests
vi.mock('@/components/RichTextEditor', () => ({
    default: ({ content, onChange, editable }: { content: string, onChange: (value: string) => void, editable: boolean }) => (
        <textarea
            data-testid="rich-text-editor"
            value={content}
            onChange={(e) => onChange(e.target.value)}
            disabled={!editable}
        />
    )
}));

// Mock the API
vi.mock('@/services/api', () => ({
    getNote: vi.fn(),
    updateNote: vi.fn(),
    createNote: vi.fn(),
    deleteNote: vi.fn(),
    getBiblePassage: vi.fn(),
    askAI: vi.fn(),
    getGroups: vi.fn(),
    shareNote: vi.fn(),
}));

// Mock scrollIntoView
window.HTMLElement.prototype.scrollIntoView = vi.fn();

describe('NoteEditor', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        vi.useRealTimers();
    });

    afterEach(() => {
        vi.useRealTimers();
    });

    it('auto-saves changes after delay', async () => {
        const mockGetNote = vi.mocked(api.getNote).mockResolvedValue({
            id: '123',
            title: 'Test Note',
            content: { markdown: 'Initial content' }
        });
        const mockUpdateNote = vi.mocked(api.updateNote).mockResolvedValue(undefined);

        render(
            <MemoryRouter initialEntries={['/notes/123']}>
                <Routes>
                    <Route path="/notes/:id" element={<NoteEditor />} />
                </Routes>
            </MemoryRouter>
        );

        await screen.findByDisplayValue('Test Note');
        expect(mockGetNote).toHaveBeenCalledWith('123');

        // Switch to fake timers
        vi.useFakeTimers();

        const textarea = screen.getByDisplayValue('Initial content');
        fireEvent.change(textarea, { target: { value: 'Updated content' } });

        act(() => {
            vi.advanceTimersByTime(2000);
        });

        // Switch back to real timers to allow waitFor to function correctly
        vi.useRealTimers();

        await waitFor(() => {
            expect(mockUpdateNote).toHaveBeenCalledWith('123', 'Test Note', { markdown: 'Updated content' });
        });
    });

    it('shows error when auto-save fails', async () => {
        const mockGetNote = vi.mocked(api.getNote).mockResolvedValue({
            id: '123',
            title: 'Test Note',
            content: { markdown: 'Initial content' }
        });
        vi.mocked(api.updateNote).mockRejectedValue(new Error('Network error'));

        const consoleErrorMock = vi.spyOn(console, 'error').mockImplementation(() => {});

        render(
            <MemoryRouter initialEntries={['/notes/123']}>
                <Routes>
                    <Route path="/notes/:id" element={<NoteEditor />} />
                </Routes>
            </MemoryRouter>
        );

        await screen.findByDisplayValue('Test Note');
        expect(mockGetNote).toHaveBeenCalledWith('123');

        vi.useFakeTimers();

        const textarea = screen.getByDisplayValue('Initial content');
        fireEvent.change(textarea, { target: { value: 'Fail content' } });

        act(() => {
            vi.advanceTimersByTime(2000);
        });

        vi.useRealTimers();

        await waitFor(() => {
             const errorIndicator = screen.queryByText(/Error saving/i) || screen.queryByText(/Failed to save/i);
             expect(errorIndicator).toBeInTheDocument();
             // In auto-save, we usually don't show a toast/alert unless manual save,
             // but checking the code: if (!manual) return; in catch block?
             // The code is: if (manual) toast.error("Failed to save");
             // So no toast should be called for auto-save failure.
             expect(toast.error).not.toHaveBeenCalled();
        });

        consoleErrorMock.mockRestore();
    });

    it('deletes note after confirmation', async () => {
        vi.mocked(api.getNote).mockResolvedValue({
            id: '123',
            title: 'Test Note',
            content: { markdown: 'Initial content' }
        });
        const mockDeleteNote = vi.mocked(api.deleteNote).mockResolvedValue(undefined);

        render(
            <MemoryRouter initialEntries={['/notes/123']}>
                <Routes>
                    <Route path="/notes/:id" element={<NoteEditor />} />
                    <Route path="/" element={<div>Dashboard</div>} />
                </Routes>
            </MemoryRouter>
        );

        await screen.findByDisplayValue('Test Note');

        // Click Delete trigger button
        // We use getAllByText because the confirm button also says "Delete" but it's not visible yet?
        // Actually, trigger is visible.
        const deleteTrigger = screen.getByRole('button', { name: 'Delete' });
        fireEvent.click(deleteTrigger);

        // Check for confirmation dialog text
        expect(await screen.findByText(/Are you sure you want to delete/i)).toBeInTheDocument();

        // Use within to find the button inside the dialog
        const dialog = await screen.findByRole('dialog');
        const confirmBtn = within(dialog).getByRole('button', { name: 'Delete' });
        fireEvent.click(confirmBtn);

        await waitFor(() => {
            expect(mockDeleteNote).toHaveBeenCalledWith('123');
            expect(screen.getByText('Dashboard')).toBeInTheDocument();
        });
    });

    it('adds bible passage', async () => {
        vi.mocked(api.getNote).mockResolvedValue({
            id: '123',
            title: 'Test Note',
            content: { markdown: 'Initial content' }
        });
        vi.mocked(api.getBiblePassage).mockResolvedValue({ verse: 'For God so loved the world...', text: 'For God so loved the world...' });

        render(
            <MemoryRouter initialEntries={['/notes/123']}>
                <Routes>
                    <Route path="/notes/:id" element={<NoteEditor />} />
                </Routes>
            </MemoryRouter>
        );

        await screen.findByDisplayValue('Test Note');

        const addScriptureBtn = screen.getByRole('button', { name: 'Add Scripture' });
        fireEvent.click(addScriptureBtn);

        const dialog = await screen.findByRole('dialog');
        const input = within(dialog).getByPlaceholderText('e.g. John 3:16');
        fireEvent.change(input, { target: { value: 'John 3:16' } });

        const searchBtn = within(dialog).getByRole('button', { name: 'Search' });
        fireEvent.click(searchBtn);

        await waitFor(() => {
            expect(api.getBiblePassage).toHaveBeenCalledWith('John 3:16');
            expect(within(dialog).getByText('For God so loved the world...')).toBeInTheDocument();
        });

        const insertBtn = within(dialog).getByRole('button', { name: 'Insert into Note' });
        fireEvent.click(insertBtn);

        const textarea = screen.getByDisplayValue(/For God so loved the world.../);
        expect(textarea).toBeInTheDocument();
    });

    it('asks AI', async () => {
        vi.mocked(api.getNote).mockResolvedValue({
            id: '123',
            title: 'Test Note',
            content: { markdown: 'Content' }
        });
        vi.mocked(api.askAI).mockResolvedValue({ response: 'AI Answer' });

        render(
            <MemoryRouter initialEntries={['/notes/123']}>
                <Routes>
                    <Route path="/notes/:id" element={<NoteEditor />} />
                </Routes>
            </MemoryRouter>
        );

        await screen.findByDisplayValue('Test Note');

        const askAiBtn = screen.getByRole('button', { name: 'Ask AI' });
        fireEvent.click(askAiBtn);

        const dialog = await screen.findByRole('dialog');
        const input = within(dialog).getByPlaceholderText('Ask a question...');
        fireEvent.change(input, { target: { value: 'Explain this' } });

        const askBtn = within(dialog).getByRole('button', { name: 'Ask' });
        fireEvent.click(askBtn);

        await waitFor(() => {
            expect(api.askAI).toHaveBeenCalledWith('Content', 'Explain this');
            expect(within(dialog).getByText('AI Answer')).toBeInTheDocument();
        });
    });

    it('auto-searches bible passage after delay', async () => {
        vi.mocked(api.getNote).mockResolvedValue({
            id: '123',
            title: 'Test Note',
            content: { markdown: 'Initial content' }
        });
        vi.mocked(api.getBiblePassage).mockResolvedValue({ verse: 'In the beginning...', text: 'In the beginning...' });

        render(
            <MemoryRouter initialEntries={['/notes/123']}>
                <Routes>
                    <Route path="/notes/:id" element={<NoteEditor />} />
                </Routes>
            </MemoryRouter>
        );

        await screen.findByDisplayValue('Test Note');

        const addScriptureBtn = screen.getByRole('button', { name: 'Add Scripture' });
        fireEvent.click(addScriptureBtn);

        const dialog = await screen.findByRole('dialog');
        const input = within(dialog).getByPlaceholderText('e.g. John 3:16');

        vi.useFakeTimers();
        fireEvent.change(input, { target: { value: 'Gen 1:1' } });

        // Fast forward
        act(() => {
            vi.advanceTimersByTime(600);
        });

        vi.useRealTimers();

        await waitFor(() => {
            expect(api.getBiblePassage).toHaveBeenCalledWith('Gen 1:1');
            expect(within(dialog).getByText('In the beginning...')).toBeInTheDocument();
        });
    });
});
