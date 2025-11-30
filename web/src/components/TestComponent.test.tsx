import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import React from 'react';

const TestComponent = () => {
  return <div>Hello Test World</div>;
};

describe('TestComponent', () => {
  it('renders correctly', () => {
    render(<TestComponent />);
    expect(screen.getByText('Hello Test World')).toBeInTheDocument();
  });
});
