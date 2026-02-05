import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { TagInput } from './TagInput';
import userEvent from '@testing-library/user-event';

describe('TagInput', () => {
    it('renders with initial tags', () => {
        const onChange = vi.fn();
        render(<TagInput value={['faith', 'love']} onChange={onChange} suggestions={[]} />);

        expect(screen.getByText('faith')).toBeInTheDocument();
        expect(screen.getByText('love')).toBeInTheDocument();
    });

    it('adds a new tag via keyboard (Enter key)', async () => {
        const user = userEvent.setup();
        const onChange = vi.fn();
        render(<TagInput value={[]} onChange={onChange} suggestions={[]} />);

        const trigger = screen.getByRole('combobox');
        await user.click(trigger);

        const input = screen.getByPlaceholderText('Search tag...');

        // Use fireEvent for change to ensure state updates immediately in JSDOM/CMDK context
        fireEvent.change(input, { target: { value: 'hope' } });

        // Then press Enter
        fireEvent.keyDown(input, { key: 'Enter', code: 'Enter' });

        expect(onChange).toHaveBeenCalledWith(['hope']);
    });

    it('adds a new tag via suggestion click', async () => {
        const user = userEvent.setup();
        const onChange = vi.fn();
        render(<TagInput value={[]} onChange={onChange} suggestions={['peace', 'joy']} />);

        const trigger = screen.getByRole('combobox');
        await user.click(trigger);

        const suggestion = await screen.findByText('peace');

        await user.click(suggestion);

        expect(onChange).toHaveBeenCalledWith(['peace']);
    });

    it('removes a tag', async () => {
        const user = userEvent.setup();
        const onChange = vi.fn();
        render(<TagInput value={['sin']} onChange={onChange} suggestions={[]} />);

        const removeButton = screen.getByRole('button', { name: /remove sin/i });
        await user.click(removeButton);

        expect(onChange).toHaveBeenCalledWith([]);
    });

    it('filters suggestions', async () => {
        const user = userEvent.setup();
        const onChange = vi.fn();
        render(<TagInput value={[]} onChange={onChange} suggestions={['apple', 'banana', 'cherry']} />);

        const trigger = screen.getByRole('combobox');
        await user.click(trigger);

        const input = screen.getByPlaceholderText('Search tag...');
        fireEvent.change(input, { target: { value: 'ban' } });

        expect(screen.getByText('banana')).toBeInTheDocument();
    });

    it('creates new tag via create option click', async () => {
        const user = userEvent.setup();
        const onChange = vi.fn();
        render(<TagInput value={[]} onChange={onChange} suggestions={['existing']} />);

        const trigger = screen.getByRole('combobox');
        await user.click(trigger);

        const input = screen.getByPlaceholderText('Search tag...');
        fireEvent.change(input, { target: { value: 'newtag' } });

        const createOption = await screen.findByText('Create "newtag"');
        await user.click(createOption);

        expect(onChange).toHaveBeenCalledWith(['newtag']);
    });
});
