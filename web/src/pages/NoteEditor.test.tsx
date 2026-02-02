import { render, screen, fireEvent, waitFor, act, within } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import NoteEditor from './NoteEditor';
import { createMemoryRouter, RouterProvider } from 'react-router-dom';
import * as api from '@/services/api';

// Mock sonner
vi.mock('sonner', () => ({
    toast: {
        success: vi.fn(),
        error: vi.fn(),
    }
}));

// Mock RichTextEditor
vi.mock('@/components/RichTextEditor', async () => {
    const React = await import('react');
    const { useEffect, useRef, useState } = React;

    type MockChain = {
        focus: () => MockChain;
        insertContent: (html: string) => MockChain;
        run: () => void;
    };

    const MockRichTextEditor = ({ initialContent, onChange, editable, onEditorReady }: { 
        initialContent: string, 
        onChange: (v: string) => void, 
        editable: boolean, 
        onEditorReady?: (e: { chain: () => MockChain }) => void 
    }) => {
        const [content, setContent] = useState(initialContent);
        const contentRef = useRef(content);

        useEffect(() => {
            contentRef.current = content;
        }, [content]);

        useEffect(() => {
            if (onEditorReady) {
                const chain = {
                    focus: () => chain,
                    insertContent: (html: string) => {
                        const newContent = contentRef.current + html;
                        setContent(newContent);
                        onChange(newContent);
                        return chain;
                    },
                    run: () => {}
                };

                const mockEditor = {
                    chain: () => chain
                };
                onEditorReady(mockEditor);
            }
            // eslint-disable-next-line react-hooks/exhaustive-deps
        }, []); // Run once on mount

        return (
            <textarea
                data-testid="rich-text-editor"
                value={content}
                onChange={(e) => {
                    setContent(e.target.value);
                    onChange(e.target.value);
                }}
                disabled={!editable}
            />
        );
    };

    return {
        default: MockRichTextEditor
    };
});

// Mock the API
vi.mock('@/services/api', () => ({
    getNote: vi.fn(),
    updateNote: vi.fn(),
    createNote: vi.fn(),
    deleteNote: vi.fn(),
    getBiblePassage: vi.fn(),
    askAI: vi.fn(),
    askAIStream: vi.fn(),
    getGroups: vi.fn(),
    shareNote: vi.fn(),
    syncUser: vi.fn().mockResolvedValue({ settings: { bible_version: 'ESV' } }),
    getBibleVersions: vi.fn().mockResolvedValue({ data: [] }),
    searchMemoryVerses: vi.fn().mockResolvedValue({ data: [] }),
    getTags: vi.fn().mockResolvedValue([]),
    getOrCreateDirectGroup: vi.fn(),
    shareItem: vi.fn(),
    getConnections: vi.fn().mockResolvedValue([]),
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

    const renderEditor = (route = '/notes/123') => {
        const router = createMemoryRouter(
            [
                { path: '/notes/:id', element: <NoteEditor /> },
                { path: '/', element: <div>Dashboard</div> }
            ],
            {
                // Ensure there is a history stack so Back button works
                initialEntries: ['/', route],
                initialIndex: 1
            }
        );
        return render(<RouterProvider router={router} />);
    };

    it('does NOT auto-save changes', async () => {
        const mockGetNote = vi.mocked(api.getNote).mockResolvedValue({
            id: '123',
            title: 'Test Note',
            content: 'Initial content'
        });
        const mockUpdateNote = vi.mocked(api.updateNote).mockResolvedValue(undefined);

        renderEditor();

        await screen.findByDisplayValue('Test Note');
        expect(mockGetNote).toHaveBeenCalledWith('123');

        vi.useFakeTimers();

        const textarea = screen.getByDisplayValue('Initial content');
        fireEvent.change(textarea, { target: { value: 'Updated content' } });

        act(() => {
            vi.advanceTimersByTime(3000);
        });

        vi.useRealTimers();

        expect(mockUpdateNote).not.toHaveBeenCalled();
    });

    it('saves manually when save button clicked', async () => {
        vi.mocked(api.getNote).mockResolvedValue({
            id: '123',
            title: 'Test Note',
            content: 'Initial content'
        });
        const mockUpdateNote = vi.mocked(api.updateNote).mockResolvedValue(undefined);

        renderEditor();

        await screen.findByDisplayValue('Test Note');

        const textarea = screen.getByDisplayValue('Initial content');
        fireEvent.change(textarea, { target: { value: 'Updated content' } });

        const saveBtn = screen.getByRole('button', { name: 'Save' });
        // It should be enabled now that it is dirty
        expect(saveBtn).not.toBeDisabled();

        fireEvent.click(saveBtn);

        await waitFor(() => {
            expect(mockUpdateNote).toHaveBeenCalledWith('123', 'Test Note', 'Updated content', []);
            expect(screen.getByText(/Saved at/i)).toBeInTheDocument();
        });
    });

    it('shows unsaved changes indicator when dirty', async () => {
        vi.mocked(api.getNote).mockResolvedValue({
            id: '123',
            title: 'Test Note',
            content: 'Initial content'
        });

        renderEditor();

        await screen.findByDisplayValue('Test Note');

        expect(screen.queryByText('Unsaved changes')).not.toBeInTheDocument();

        const textarea = screen.getByDisplayValue('Initial content');
        fireEvent.change(textarea, { target: { value: 'Updated content' } });

        // useBlocker/isDirty might take a render cycle
        await waitFor(() => {
            expect(screen.getByText('Unsaved changes')).toBeInTheDocument();
        });
    });

    it('warns on navigation when dirty', async () => {
        vi.mocked(api.getNote).mockResolvedValue({
            id: '123',
            title: 'Test Note',
            content: 'Initial content'
        });

        renderEditor();

        await screen.findByDisplayValue('Test Note');

        const textarea = screen.getByDisplayValue('Initial content');
        fireEvent.change(textarea, { target: { value: 'Updated content' } });

        // Try to navigate away via Back button
        const backBtn = screen.getByText('← Back');
        fireEvent.click(backBtn);

        // Expect AlertDialog
        expect(await screen.findByText('Unsaved Changes')).toBeInTheDocument();
        expect(screen.getByText('You have unsaved changes. Do you want to save them before leaving?')).toBeInTheDocument();
        expect(screen.getByText('Save & Leave')).toBeInTheDocument();
        expect(screen.getByText('Discard Changes')).toBeInTheDocument();

        // Click Cancel
        const cancelBtn = screen.getByText('Cancel');
        fireEvent.click(cancelBtn);

        // Dialog should close and we stay
        await waitFor(() => {
             expect(screen.queryByText('Unsaved Changes')).not.toBeInTheDocument();
             expect(screen.getByDisplayValue('Updated content')).toBeInTheDocument();
        });

        // Try again and leave
        fireEvent.click(backBtn);
        const discardBtn = await screen.findByText('Discard Changes');
        fireEvent.click(discardBtn);

        // Should navigate to dashboard
        await waitFor(() => {
             expect(screen.getByText('Dashboard')).toBeInTheDocument();
        });
    });

    it('deletes note after confirmation', async () => {
        vi.mocked(api.getNote).mockResolvedValue({
            id: '123',
            title: 'Test Note',
            content: 'Initial content'
        });
        const mockDeleteNote = vi.mocked(api.deleteNote).mockResolvedValue(undefined);

        renderEditor();

        await screen.findByDisplayValue('Test Note');

        const deleteTrigger = screen.getByRole('button', { name: 'Delete' });
        fireEvent.click(deleteTrigger);

        expect(await screen.findByText(/Are you sure you want to delete/i)).toBeInTheDocument();

        const dialog = await screen.findByRole('dialog');
        // In the dialog, we have "Delete Note" (title) and "Delete" (button)
        // Or we can find by the button that is destructive
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
            content: 'Initial content'
        });
        vi.mocked(api.getBiblePassage).mockResolvedValue({ verse: 'For God so loved the world...', text: 'For God so loved the world...' });

        renderEditor();

        await screen.findByDisplayValue('Test Note');

        const addScriptureBtn = screen.getByRole('button', { name: 'Add Scripture' });
        fireEvent.click(addScriptureBtn);

        const dialog = await screen.findByRole('dialog');
        const input = within(dialog).getByPlaceholderText('e.g. John 3:16');
        fireEvent.change(input, { target: { value: 'John 3:16' } });

        const searchBtn = within(dialog).getByRole('button', { name: 'Search' });
        fireEvent.click(searchBtn);

        await waitFor(() => {
            expect(api.getBiblePassage).toHaveBeenCalledWith('John 3:16', 'ESV');
            expect(within(dialog).getByText('For God so loved the world...')).toBeInTheDocument();
        });

        const insertBtn = within(dialog).getByRole('button', { name: 'Insert into Note' });
        fireEvent.click(insertBtn);

        await waitFor(() => {
            const textarea = screen.getByTestId('rich-text-editor') as HTMLTextAreaElement;
            expect(textarea.value).toContain('For God so loved the world...');
        });
    });

    it('asks AI and inserts response into note', async () => {
        vi.mocked(api.getNote).mockResolvedValue({
            id: '123',
            title: 'Test Note',
            content: 'Content'
        });
        vi.mocked(api.askAIStream).mockImplementation(async (_context, _prompt, _version, callbacks) => {
            if (callbacks.onStart) callbacks.onStart("note-id");
            if (callbacks.onChunk) callbacks.onChunk('AI Answer');
            if (callbacks.onDone) callbacks.onDone();
        });

        renderEditor();

        await screen.findByDisplayValue('Test Note');

        const askAiBtn = screen.getByRole('button', { name: 'Ask AI' });
        fireEvent.click(askAiBtn);

        const dialog = await screen.findByRole('dialog');
        const input = within(dialog).getByPlaceholderText('Ask a question...');
        fireEvent.change(input, { target: { value: 'Explain this' } });

        const askBtn = within(dialog).getByRole('button', { name: 'Ask' });
        fireEvent.click(askBtn);

        await waitFor(() => {
            expect(api.askAIStream).toHaveBeenCalledWith('Content', 'Explain this', 'ESV', expect.any(Object));
            expect(within(dialog).getByText('AI Answer')).toBeInTheDocument();
        });

        const insertBtn = within(dialog).getByRole('button', { name: 'Insert into Note' });
        fireEvent.click(insertBtn);

        await waitFor(() => {
            const textarea = screen.getByTestId('rich-text-editor') as HTMLTextAreaElement;
            expect(textarea.value).toMatch(/AI Answer/);
            expect(textarea.value).toMatch(/Question: Explain this/);
            expect(textarea.value).toMatch(/Content/);
        });
    });
});
