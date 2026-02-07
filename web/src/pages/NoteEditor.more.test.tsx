import { render, screen, fireEvent, waitFor, act, within } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import NoteEditor from './NoteEditor';
import { createMemoryRouter, RouterProvider } from 'react-router-dom';
import userEvent from '@testing-library/user-event';
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
    shareItem: vi.fn(),
    syncUser: vi.fn().mockResolvedValue({ settings: { bible_version: 'ESV' }, email: 'me@example.com' }),
    getBibleVersions: vi.fn().mockResolvedValue({ data: [] }),
    searchMemoryVerses: vi.fn().mockResolvedValue({ data: [] }),
    getTags: vi.fn().mockResolvedValue([]),
    getOrCreateDirectGroup: vi.fn(),
    getConnections: vi.fn().mockResolvedValue([]),
}));

// Mock scrollIntoView
window.HTMLElement.prototype.scrollIntoView = vi.fn();

describe('NoteEditor Additional Coverage', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        vi.useRealTimers();
    });

    afterEach(() => {
        vi.useRealTimers();
    });

    const renderEditor = (route = '/notes/123', state: any = null) => {
        const router = createMemoryRouter(
            [
                { path: '/notes/:id', element: <NoteEditor /> },
                { path: '/', element: <div>Dashboard</div> }
            ],
            {
                initialEntries: [
                    { pathname: route, state: state }
                ],
                initialIndex: 0
            }
        );
        return render(<RouterProvider router={router} />);
    };

    it('handles memory verse insertion', async () => {
        vi.mocked(api.getNote).mockResolvedValue({
            id: '123',
            title: 'Test Note',
            content: 'Content'
        });
        vi.mocked(api.searchMemoryVerses).mockResolvedValue({
            data: [
                { id: 'v1', reference: 'John 3:16', text: 'For God so loved...', pack_title: 'Pack 1', tags: [] }
            ]
        });
        vi.mocked(api.getBiblePassage).mockResolvedValue({ verse: 'For God so loved...', reference: 'John 3:16' });

        renderEditor();
        await screen.findByDisplayValue('Test Note');

        const verseBtn = screen.getByRole('button', { name: 'Insert Memory Verse' });
        fireEvent.click(verseBtn);

        const dialog = await screen.findByRole('dialog');
        const input = within(dialog).getByPlaceholderText('Search by reference, pack or tags...');
        fireEvent.change(input, { target: { value: 'John' } });

        // Wait for debounce
        await act(async () => {
            await new Promise(r => setTimeout(r, 400));
        });

        await waitFor(() => {
            expect(api.searchMemoryVerses).toHaveBeenCalledWith('John');
            expect(within(dialog).getByText(/John 3:16/)).toBeInTheDocument();
        });

        const verseItem = within(dialog).getByText(/John 3:16/);
        fireEvent.click(verseItem);

        await waitFor(() => {
            const textarea = screen.getByTestId('rich-text-editor') as HTMLTextAreaElement;
            expect(textarea.value).toContain('John 3:16');
            expect(textarea.value).toContain('For God so loved...');
        });
    });

    it('handles sharing to a group', async () => {
        vi.mocked(api.getNote).mockResolvedValue({
            id: '123',
            title: 'Test Note',
            content: 'Content'
        });
        vi.mocked(api.getGroups).mockResolvedValue([
            { id: 'g1', name: 'Study Group' }
        ]);
        vi.mocked(api.shareItem).mockResolvedValue({ success: true });

        renderEditor();
        await screen.findByDisplayValue('Test Note');

        const shareBtn = screen.getByRole('button', { name: 'Share' });
        fireEvent.click(shareBtn);

        const dialog = await screen.findByRole('dialog');
        await waitFor(() => {
            expect(within(dialog).getByText('Study Group')).toBeInTheDocument();
        });

        // Select Group
        const select = within(dialog).getByRole('combobox');
        fireEvent.change(select, { target: { value: 'g1' } });

        const commentInput = within(dialog).getByPlaceholderText('Add a comment (optional)...');
        fireEvent.change(commentInput, { target: { value: 'Check this out' } });

        const confirmShare = within(dialog).getByRole('button', { name: 'Share Note' });
        fireEvent.click(confirmShare);

        await waitFor(() => {
            expect(api.shareItem).toHaveBeenCalledWith('g1', { note_id: '123', comment: 'Check this out' });
            expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
        });
    });

    it('handles sharing to a connection', async () => {
        const user = userEvent.setup();
        vi.mocked(api.getNote).mockResolvedValue({
            id: '123',
            title: 'Test Note',
            content: 'Content'
        });
        vi.mocked(api.getConnections).mockResolvedValue([
            { id: 'c1', requester_email: 'me@example.com', receiver_email: 'friend@example.com', requester_id: 'me', receiver_id: 'friend', status: 'accepted' }
        ]);
        vi.mocked(api.getOrCreateDirectGroup).mockResolvedValue({ id: 'dg1' });
        vi.mocked(api.shareItem).mockResolvedValue({ success: true });

        renderEditor();
        await screen.findByDisplayValue('Test Note');

        const shareBtn = screen.getByRole('button', { name: 'Share' });
        await user.click(shareBtn);

        const dialog = await screen.findByRole('dialog');

        await waitFor(() => {
             expect(api.getConnections).toHaveBeenCalled();
        });

        const connectionsTab = within(dialog).getByRole('tab', { name: /Direct/ });
        await user.click(connectionsTab);

        await waitFor(() => {
            expect(within(dialog).getByText('Select a Connection...')).toBeInTheDocument();
        });

        await waitFor(() => {
             expect(within(dialog).getByText(/@example.com/)).toBeInTheDocument();
        });

        const selects = within(dialog).getAllByRole('combobox');
        const select = selects.find(s => s.innerHTML.includes('example.com'));

        if (select) {
             fireEvent.change(select, { target: { value: 'friend' } }); // friend is receiver_id
        } else {
            throw new Error("Connection select not found");
        }

        const confirmShare = within(dialog).getByRole('button', { name: 'Share Note' });
        await user.click(confirmShare);

        await waitFor(() => {
            expect(api.getOrCreateDirectGroup).toHaveBeenCalledWith('friend');
            expect(api.shareItem).toHaveBeenCalledWith('dg1', { note_id: '123', comment: '' });
        });
    });

    it('toggles preview mode', async () => {
        vi.mocked(api.getNote).mockResolvedValue({
            id: '123',
            title: 'Test Note',
            content: '<b>Bold Content</b>'
        });

        renderEditor();
        await screen.findByDisplayValue('Test Note');

        const previewBtn = screen.getByRole('button', { name: 'Preview' });
        fireEvent.click(previewBtn);

        expect(screen.queryByTestId('rich-text-editor')).not.toBeInTheDocument();
        expect(screen.getByText('Bold Content')).toBeInTheDocument();

        const editBtn = screen.getByRole('button', { name: 'Edit' });
        fireEvent.click(editBtn);

        expect(screen.getByTestId('rich-text-editor')).toBeInTheDocument();
    });

    it('initializes from location state (Create Note)', async () => {
        vi.mocked(api.getBiblePassage).mockResolvedValue({ verse: 'The Verse Text', reference: 'Gen 1:1' });

        renderEditor('/notes/new', {
            title: 'New Note',
            tags: ['tag1'],
            content: 'Start',
            passageRef: 'Gen 1:1',
            passageVersion: 'KJV'
        });

        await waitFor(() => {
            expect(screen.getByDisplayValue('New Note')).toBeInTheDocument();
        });

        await waitFor(() => {
            expect(api.getBiblePassage).toHaveBeenCalledWith('Gen 1:1', 'KJV');
            const textarea = screen.getByTestId('rich-text-editor') as HTMLTextAreaElement;
            expect(textarea.value).toContain('The Verse Text');
            expect(textarea.value).toContain('Start');
        });
    });

    it('handles save error', async () => {
        vi.mocked(api.getNote).mockResolvedValue({
            id: '123',
            title: 'Test Note',
            content: 'Content'
        });
        vi.mocked(api.updateNote).mockRejectedValue(new Error('Network Error'));

        renderEditor();
        await screen.findByDisplayValue('Test Note');

        const textarea = screen.getByTestId('rich-text-editor');
        fireEvent.change(textarea, { target: { value: 'New Content' } });

        const saveBtn = screen.getByRole('button', { name: 'Save' });
        fireEvent.click(saveBtn);

        await waitFor(() => {
            expect(screen.getByText('Error saving')).toBeInTheDocument();
        });
    });
});
