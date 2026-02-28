import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import App from './App';

// Mock Wails runtime calls
vi.mock('../wailsjs/go/main/App', () => ({
  GetPhotos: vi.fn(() => Promise.resolve([])),
  SelectFolder: vi.fn(),
  ImportFromFolder: vi.fn(),
  UpdatePhotoDate: vi.fn(),
}));

describe('App Component', () => {
  it('renders without crashing', async () => {
    render(<App />);
    expect(screen.getByText('Photoo')).toBeDefined();
  });
});
