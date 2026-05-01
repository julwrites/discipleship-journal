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

    it('navigates on Enter key and Space key', () => {
        render(
             <MemoryRouter>
                <NoteCard note={mockNote} {...mockHandlers} />
            </MemoryRouter>
        );

        const card = screen.getByRole('link');

        fireEvent.keyDown(card, { key: 'Enter', code: 'Enter' });
        expect(mockNavigate).toHaveBeenCalledWith('/notes/1');

        fireEvent.keyDown(card, { key: ' ', code: 'Space' });
        expect(mockNavigate).toHaveBeenCalledTimes(2);
    });

    it('ignores other key presses', () => {
        render(
             <MemoryRouter>
                <NoteCard note={mockNote} {...mockHandlers} />
            </MemoryRouter>
        );

        const card = screen.getByRole('link');
        fireEvent.keyDown(card, { key: 'A', code: 'KeyA' });

        expect(mockNavigate).not.toHaveBeenCalled();
    });

    it('renders untitled note if title is missing', () => {
        render(
            <MemoryRouter>
                <NoteCard note={{...mockNote, title: ''}} {...mockHandlers} />
            </MemoryRouter>
        );

        expect(screen.getByText('Untitled Note')).toBeInTheDocument();
    });

    it('renders correctly without tags', () => {
         render(
            <MemoryRouter>
                <NoteCard note={{...mockNote, tags: undefined}} {...mockHandlers} />
            </MemoryRouter>
        );
         expect(screen.queryByText('Faith')).not.toBeInTheDocument();
    });

    it('renders empty tag list correctly', () => {
         render(
            <MemoryRouter>
                <NoteCard note={{...mockNote, tags: []}} {...mockHandlers} />
            </MemoryRouter>
        );
         expect(screen.queryByText('Faith')).not.toBeInTheDocument();
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
        expect(mockNavigate).not.toHaveBeenCalled();

        fireEvent.click(shareButton);
        expect(mockHandlers.onShare).toHaveBeenCalledWith(mockNote);

        fireEvent.click(askAIButton);
        expect(mockHandlers.onAskAI).toHaveBeenCalledWith(mockNote);
    });
});
