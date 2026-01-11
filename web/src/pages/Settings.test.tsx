import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeAll } from 'vitest';
import userEvent from '@testing-library/user-event';
import Settings from './Settings';
import { ThemeProvider } from '@/components/theme/ThemeProvider';
import { BrowserRouter } from 'react-router-dom';
import * as api from '@/services/api';

// Mock API
vi.mock('@/services/api', () => ({
  syncUser: vi.fn().mockResolvedValue({ username: 'testuser', settings: { bible_version: 'ESV' } }),
  updateUser: vi.fn().mockResolvedValue({}),
  getBibleVersions: vi.fn().mockResolvedValue({
    data: [
        { id: "uuid-1", name: "English Standard Version", abbreviation: "ESV", language: "en" },
        { id: "uuid-2", name: "New International Version", abbreviation: "NIV", language: "en" }
    ]
  }),
}));

describe('Settings Page', () => {
  beforeAll(() => {
    // Polyfill ResizeObserver for cmdk
    global.ResizeObserver = class ResizeObserver {
      observe() {}
      unobserve() {}
      disconnect() {}
    };

    // Polyfill scrollIntoView for cmdk
    window.HTMLElement.prototype.scrollIntoView = function() {};

    Object.defineProperty(window, 'matchMedia', {
      writable: true,
      value: vi.fn().mockImplementation(query => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: vi.fn(), // deprecated
        removeListener: vi.fn(), // deprecated
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
        dispatchEvent: vi.fn(),
      })),
    });
  });

  it('renders appearance settings', async () => {
    render(
      <BrowserRouter>
        <ThemeProvider>
          <Settings />
        </ThemeProvider>
      </BrowserRouter>
    );

    // Wait for loading to finish and check for "Appearance" section
    await waitFor(() => {
        expect(screen.getByText('Appearance')).toBeInTheDocument();
    });

    // Check for theme buttons
    expect(screen.getByText('Default')).toBeInTheDocument();
    expect(screen.getByText('Serene')).toBeInTheDocument();
    expect(screen.getByText('Elegant')).toBeInTheDocument();

    // Check for mode buttons
    expect(screen.getByText('Light')).toBeInTheDocument();
    expect(screen.getByText('Dark')).toBeInTheDocument();
    expect(screen.getByText('System')).toBeInTheDocument();
  });

  it('updates bible version using abbreviation', async () => {
    const user = userEvent.setup();
    render(
      <BrowserRouter>
        <ThemeProvider>
          <Settings />
        </ThemeProvider>
      </BrowserRouter>
    );

    await waitFor(() => {
        expect(screen.getByText('Profile')).toBeInTheDocument();
    });

    // Find the selector (it shows current version "ESV" or placeholder)
    // Initially loaded with ESV
    // We search by role combobox. The accessible name should be the text content.
    const trigger = screen.getByRole('combobox');
    expect(trigger).toBeInTheDocument();
    expect(trigger).toHaveTextContent(/ESV/);

    // Open the dropdown
    await user.click(trigger);

    // Select NIV
    // Wait for options to appear
    const nivOption = await screen.findByText('NIV');
    await user.click(nivOption);

    // Click Save
    const saveButton = screen.getByText('Save Changes');
    await user.click(saveButton);

    // Expect updateUser to be called with 'NIV' (abbreviation), not 'uuid-2' (ID)
    expect(api.updateUser).toHaveBeenCalledWith(expect.objectContaining({
        bible_version: 'NIV'
    }));
  });
});
