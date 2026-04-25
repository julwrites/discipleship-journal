import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { NoteCard, Note } from './NoteCard';
import { MemoryRouter } from 'react-router-dom';

// Mock useNavigate
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async (importOriginal) => {
  const actual = await importOriginal<typeof import('react-router-dom')>();
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

describe('NoteCard', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    const mockNote: Note = {
        id: '1',
        title: 'Test Note',
        updated_at: '2023-01-01T00:00:00Z',
        tags: [
            { id: 't1', name: 'Faith' },
            { id: 't2', name: 'Prayer' }
        ]
    };

    const mockHandlers = {
        onDelete: vi.fn(),
        onShare: vi.fn(),
        onAskAI: vi.fn()
    };

    it('renders note details correctly', () => {
        render(
            <MemoryRouter>
                <NoteCard note={mockNote} {...mockHandlers} />
            </MemoryRouter>
        );

        expect(screen.getByText('Test Note')).toBeInTheDocument();
        expect(screen.getByText('Faith')).toBeInTheDocument();
        expect(screen.getByText('Prayer')).toBeInTheDocument();
        // Date formatting might depend on locale, but let's check basic presence
        // 1/1/2023 or similar
        // expect(screen.getByText(/2023/)).toBeInTheDocument();
    });

    it('navigates to note detail on card click', () => {
        render(
            <MemoryRouter>
                <NoteCard note={mockNote} {...mockHandlers} />
            </MemoryRouter>
        );

        const card = screen.getByRole('link');
        fireEvent.click(card);

        expect(mockNavigate).toHaveBeenCalledWith('/notes/1');
    });

    it('navigates on Space key', () => {
        render(
             <MemoryRouter>
                <NoteCard note={mockNote} {...mockHandlers} />
            </MemoryRouter>
        );

        const card = screen.getByRole('link');
        fireEvent.keyDown(card, { key: ' ', code: 'Space' });

        expect(mockNavigate).toHaveBeenCalledWith('/notes/1');
    });

    it('navigates on Enter key', () => {
        render(
             <MemoryRouter>
                <NoteCard note={mockNote} {...mockHandlers} />
            </MemoryRouter>
        );

        const card = screen.getByRole('link');
        fireEvent.keyDown(card, { key: 'Enter', code: 'Enter' });

        expect(mockNavigate).toHaveBeenCalledWith('/notes/1');
    });

    it('triggers actions without navigation', () => {
        render(
             <MemoryRouter>
                <NoteCard note={mockNote} {...mockHandlers} />
            </MemoryRouter>
        );

        const deleteButton = screen.getByTitle('Delete');
        const shareButton = screen.getByTitle('Share');
        const askAIButton = screen.getByTitle('Ask AI');

        fireEvent.click(deleteButton);
        expect(mockHandlers.onDelete).toHaveBeenCalledWith(mockNote);
        expect(mockNavigate).not.toHaveBeenCalled(); // Should assume test isolation or clear mocks

        fireEvent.click(shareButton);
        expect(mockHandlers.onShare).toHaveBeenCalledWith(mockNote);

        fireEvent.click(askAIButton);
        expect(mockHandlers.onAskAI).toHaveBeenCalledWith(mockNote);
    });

    it('has correct aria-labels for action buttons', () => {
        render(
             <MemoryRouter>
                <NoteCard note={mockNote} {...mockHandlers} />
            </MemoryRouter>
        );

        const deleteButton = screen.getByTitle('Delete');
        const shareButton = screen.getByTitle('Share');
        const askAIButton = screen.getByTitle('Ask AI');

        expect(deleteButton).toHaveAttribute('aria-label', `Delete ${mockNote.title}`);
        expect(shareButton).toHaveAttribute('aria-label', `Share ${mockNote.title}`);
        expect(askAIButton).toHaveAttribute('aria-label', `Ask AI about ${mockNote.title}`);
    });

    it('uses fallback text for aria-labels when title is missing', () => {
        const noteWithoutTitle = { ...mockNote, title: '' } as unknown as Note;
        render(
             <MemoryRouter>
                <NoteCard note={noteWithoutTitle} {...mockHandlers} />
            </MemoryRouter>
        );

        const deleteButton = screen.getByTitle('Delete');
        const shareButton = screen.getByTitle('Share');
        const askAIButton = screen.getByTitle('Ask AI');

        expect(deleteButton).toHaveAttribute('aria-label', 'Delete Untitled Note');
        expect(shareButton).toHaveAttribute('aria-label', 'Share Untitled Note');
        expect(askAIButton).toHaveAttribute('aria-label', 'Ask AI about Untitled Note');
    });

    it('renders tags when available', () => {
        render(
             <MemoryRouter>
                <NoteCard note={mockNote} {...mockHandlers} />
            </MemoryRouter>
        );

        expect(screen.getByText('Faith')).toBeInTheDocument();
        expect(screen.getByText('Prayer')).toBeInTheDocument();
    });

    it('does not render tags if empty or undefined', () => {
        const noteWithoutTags = { ...mockNote, tags: [] };
        const { rerender, unmount } = render(
             <MemoryRouter>
                <NoteCard note={noteWithoutTags} {...mockHandlers} />
            </MemoryRouter>
        );

        expect(screen.queryByText('Faith')).not.toBeInTheDocument();
        unmount();

        const noteUndefinedTags = { ...mockNote, tags: undefined };
        render(
             <MemoryRouter>
                <NoteCard note={noteUndefinedTags} {...mockHandlers} />
            </MemoryRouter>
        );
        expect(screen.queryByText('Faith')).not.toBeInTheDocument();
    });

    it('renders with fallback title when title is empty', () => {
        const noteWithoutTitle = { ...mockNote, title: '' };
        render(
             <MemoryRouter>
                <NoteCard note={noteWithoutTitle} {...mockHandlers} />
            </MemoryRouter>
        );

        expect(screen.getByText('Untitled Note')).toBeInTheDocument();
    });

    it('ignores non-Enter/Space keys', () => {
        render(
             <MemoryRouter>
                <NoteCard note={mockNote} {...mockHandlers} />
            </MemoryRouter>
        );

        const card = screen.getByRole('link');
        fireEvent.keyDown(card, { key: 'a', code: 'KeyA' });

        expect(mockNavigate).not.toHaveBeenCalled();
    });
});
