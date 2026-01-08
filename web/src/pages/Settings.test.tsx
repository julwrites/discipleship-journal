import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeAll } from 'vitest';
import Settings from './Settings';
import { ThemeProvider } from '@/components/theme/ThemeProvider';
import { BrowserRouter } from 'react-router-dom';

// Mock API
vi.mock('@/services/api', () => ({
  syncUser: vi.fn().mockResolvedValue({ username: 'testuser', settings: { bible_version: 'ESV' } }),
  updateUser: vi.fn().mockResolvedValue({}),
}));

describe('Settings Page', () => {
  beforeAll(() => {
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
});
