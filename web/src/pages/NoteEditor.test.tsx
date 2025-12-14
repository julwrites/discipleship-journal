import { render, screen, fireEvent, waitFor, act } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import NoteEditor from './NoteEditor';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import * as api from '@/services/api';

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

        const alertMock = vi.spyOn(window, 'alert').mockImplementation(() => {});
        const consoleErrorMock = vi.spyOn(console, 'error').mockImplementation(() => {});

        render(
            <MemoryRouter initialEntries={['/notes/123']}>
                <Routes>
                    <Route path="/notes/:id" element={<NoteEditor />} />
                </Routes>
            </MemoryRouter>
        );

        await screen.findByDisplayValue('Test Note');

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
        });

        alertMock.mockRestore();
        consoleErrorMock.mockRestore();
    });
});
