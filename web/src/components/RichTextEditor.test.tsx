import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import RichTextEditor from './RichTextEditor';

describe('RichTextEditor', () => {
    it('allows adding an image with alt text', async () => {
        const handleChange = vi.fn();
        render(<RichTextEditor initialContent="" onChange={handleChange} />);

        // Find Image button
        const imageBtn = screen.getByTitle('Image');
        fireEvent.click(imageBtn);

        // Popover should open
        // Wait for popover content
        const urlInput = await screen.findByLabelText('URL');
        fireEvent.change(urlInput, { target: { value: 'https://example.com/image.jpg' } });

        // Alt Text input
        // This is expected to fail initially until we implement the feature
        const altInput = screen.getByLabelText('Alt Text');
        fireEvent.change(altInput, { target: { value: 'A beautiful sunset' } });

        const addBtn = screen.getByRole('button', { name: 'Add Image' });
        fireEvent.click(addBtn);

        // Verify content
        await waitFor(() => {
            // Check the last call to onChange
            expect(handleChange).toHaveBeenCalled();
            const lastCall = handleChange.mock.calls[handleChange.mock.calls.length - 1];
            const html = lastCall[0];

            // Check for image tag with src and alt
            // Note: Tiptap might add other attributes or order them differently, so we check existence
            expect(html).toContain('src="https://example.com/image.jpg"');
            expect(html).toContain('alt="A beautiful sunset"');
        });
    });
});
