import { render, screen, fireEvent } from '@testing-library/react';
import { BibleReferenceInput } from './BibleReferenceInput';
import { vi, describe, it, expect, beforeEach } from 'vitest';
import * as bibleRefUtils from '@/utils/bibleReference';

// Mock the utility
vi.mock('@/utils/bibleReference', () => ({
    parseBibleReference: vi.fn(),
}));

describe('BibleReferenceInput', () => {
    const mockOnChange = vi.fn();
    const mockOnBlur = vi.fn();

    beforeEach(() => {
        vi.clearAllMocks();
        mockOnChange.mockClear();
        mockOnBlur.mockClear();
    });

    it('renders an input by default', () => {
        render(<BibleReferenceInput value="" onChange={mockOnChange} />);
        const input = screen.getByRole('textbox');
        expect(input.tagName).toBe('INPUT');
    });

    it('renders a textarea when multiline is true', () => {
        render(<BibleReferenceInput value="" onChange={mockOnChange} multiline />);
        const input = screen.getByRole('textbox');
        expect(input.tagName).toBe('TEXTAREA');
    });

    it('calls onChange when typing', () => {
        render(<BibleReferenceInput value="" onChange={mockOnChange} />);
        const input = screen.getByRole('textbox');
        fireEvent.change(input, { target: { value: 'jn 3:16' } });
        expect(mockOnChange).toHaveBeenCalledWith('jn 3:16');
    });

    it('normalizes valid reference on blur', () => {
        vi.mocked(bibleRefUtils.parseBibleReference).mockReturnValue('John 3:16');

        render(<BibleReferenceInput value="jn 3:16" onChange={mockOnChange} />);
        const input = screen.getByRole('textbox');

        fireEvent.blur(input, { target: { value: 'jn 3:16' } });

        expect(bibleRefUtils.parseBibleReference).toHaveBeenCalledWith('jn 3:16');
        expect(mockOnChange).toHaveBeenCalledWith('John 3:16');
    });

    it('does not change invalid reference on blur', () => {
        vi.mocked(bibleRefUtils.parseBibleReference).mockReturnValue(null);

        render(<BibleReferenceInput value="invalid" onChange={mockOnChange} />);
        const input = screen.getByRole('textbox');

        fireEvent.blur(input, { target: { value: 'invalid' } });

        expect(bibleRefUtils.parseBibleReference).toHaveBeenCalledWith('invalid');
        expect(mockOnChange).not.toHaveBeenCalled();
    });

    it('handles allowMultiple with comma separation', () => {
        // Mock responses for individual parts
        vi.mocked(bibleRefUtils.parseBibleReference).mockImplementation((ref: string) => {
            if (ref === 'jn 3:16') return 'John 3:16';
            if (ref === 'rom 1:1') return 'Romans 1:1';
            return null;
        });

        render(<BibleReferenceInput value="jn 3:16, rom 1:1" onChange={mockOnChange} allowMultiple />);
        const input = screen.getByRole('textbox');

        fireEvent.blur(input, { target: { value: 'jn 3:16, rom 1:1' } });

        expect(mockOnChange).toHaveBeenCalledWith('John 3:16, Romans 1:1');
    });

    it('calls onBlur prop if provided', () => {
        render(<BibleReferenceInput value="" onChange={mockOnChange} onBlur={mockOnBlur} />);
        const input = screen.getByRole('textbox');
        fireEvent.blur(input);
        expect(mockOnBlur).toHaveBeenCalled();
    });

    it('does nothing on blur if value is empty', () => {
        render(<BibleReferenceInput value="" onChange={mockOnChange} />);
        const input = screen.getByRole('textbox');
        fireEvent.blur(input, { target: { value: '   ' } });
        expect(bibleRefUtils.parseBibleReference).not.toHaveBeenCalled();
        expect(mockOnChange).not.toHaveBeenCalled();
    });
});
