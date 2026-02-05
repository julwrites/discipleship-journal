import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { NoteCard, Note } from './NoteCard';
import { MemoryRouter, useNavigate } from 'react-router-dom';

// Mock useNavigate
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async (importOriginal) => {
  const actual = await importOriginal();
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
});
