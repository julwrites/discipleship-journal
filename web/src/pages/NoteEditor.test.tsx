import { render, screen, fireEvent, waitFor, act, within } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import NoteEditor from './NoteEditor';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import * as api from '@/services/api';
import { toast } from 'sonner';
import { useEffect } from 'react';

// Mock sonner
vi.mock('sonner', () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
    }
}));

// Mock RichTextEditor to avoid complex Tiptap interaction in integration tests
vi.mock('@/components/RichTextEditor', () => ({
    default: ({ initialContent, onChange, editable, onEditorReady }: { initialContent: string, onChange: (value: string) => void, editable: boolean, onEditorReady: (editor: any) => void }) => {
        // We need a local state to track content for the text area, initialized from initialContent
        // In the real app, Tiptap handles its own state. Here we need to mimic it.
        // But since we are mocking, we can just use the props passed from parent if parent updates it.
        // HOWEVER, NoteEditor is now uncontrolled-ish for the editor. It passes initialContent only once.
        // But it DOES listen to onChange.
        // For the tests to see "updates", we need the mock to call onChange.

        // Let's create a fake editor instance to pass back to parent
        useEffect(() => {
            if (onEditorReady) {
                const mockEditor = {
                    chain: () => ({
                        focus: () => ({
                            insertContent: (html: string) => ({
                                run: () => {
                                    // Simulate Tiptap inserting content by appending to the current content
                                    // and calling onChange
                                    // Note: We don't have easy access to the *current* content here inside the closure
                                    // if we just use initialContent.
                                    // But we can approximate by assuming the test setup.
                                    // A better way is to fire the onChange with the NEW combined content.
                                    // Since we don't track state inside this mock properly, let's just
                                    // call onChange with the inserted content appended to a placeholder or similar,
                                    // OR, simpler: just call onChange with the HTML so the test sees it.
                                    onChange(initialContent + html);
                                }
                            })
                        })
                    })
                };
                onEditorReady(mockEditor);
            }
        }, [onEditorReady, initialContent, onChange]);

        return (
            <textarea
                data-testid="rich-text-editor"
                value={initialContent}
                onChange={(e) => onChange(e.target.value)}
                disabled={!editable}
            />
        );
    }
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

        // Note: The textarea value will be 'Initial content' because we use initialContent prop.
        // Changing it fires onChange which updates parent 'markdown' state.
        // But parent 'markdown' state does NOT flow back into 'initialContent' prop of RichTextEditor
        // because we removed that sync.
        // However, the test uses fireEvent.change on the textarea. The mock textarea uses `value={initialContent}`.
        // If parent re-renders, does it pass new initialContent?
        // No, NoteEditor uses `markdown` state as `initialContent`.
        // So when `setMarkdown` is called, `NoteEditor` re-renders, passing new `markdown` as `initialContent`.
        // So the textarea SHOULD update.

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

        const deleteTrigger = screen.getByRole('button', { name: 'Delete' });
        fireEvent.click(deleteTrigger);

        expect(await screen.findByText(/Are you sure you want to delete/i)).toBeInTheDocument();

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

        await waitFor(() => {
            const textarea = screen.getByDisplayValue((content) => content.includes('For God so loved the world...'));
            expect(textarea).toBeInTheDocument();
        });
    });

    it('asks AI and inserts response into note', async () => {
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

        const insertBtn = within(dialog).getByRole('button', { name: 'Insert into Note' });
        fireEvent.click(insertBtn);

        await waitFor(() => {
            const textarea = screen.getByTestId('rich-text-editor');
            expect(textarea).toHaveValue(expect.stringMatching(/AI Answer/));
            expect(textarea).toHaveValue(expect.stringMatching(/Question: Explain this/));
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
