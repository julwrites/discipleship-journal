import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { Calendar } from './calendar';

describe('Calendar', () => {
  it('renders correctly', () => {
    render(<Calendar />);
    // React Day Picker usually renders a table with "grid" role
    expect(screen.getByRole('grid')).toBeInTheDocument();
  });

  it('renders with custom class name', () => {
    const customClass = 'test-class';
    const { container } = render(<Calendar className={customClass} />);
    // The class is applied to the root element or internal div
    expect(container.getElementsByClassName(customClass).length).toBeGreaterThan(0);
  });
});
