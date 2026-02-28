import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import App from './App';
import { GetPhotos } from '../wailsjs/go/main/App';

// Mock Wails runtime calls
vi.mock('../wailsjs/go/main/App', () => ({
  GetPhotos: vi.fn(),
  SelectFolder: vi.fn(),
  ImportFromFolder: vi.fn(),
  UpdatePhotoDate: vi.fn(),
}));

describe('App Component', () => {
  it('renders without crashing', async () => {
    vi.mocked(GetPhotos).mockResolvedValue([]);
    render(<App />);
    expect(screen.getByText('Photoo')).toBeDefined();
  });

  it('displays empty state when no photos are returned', async () => {
    vi.mocked(GetPhotos).mockResolvedValue([]);
    render(<App />);
    const emptyState = await screen.findByText('No photos imported yet.');
    expect(emptyState).toBeDefined();
  });

  it('renders a grid of photos when data is received', async () => {
    const mockPhotos = [
      {
        id: 1,
        filename: 'test.jpg',
        date_taken: new Date('2023-01-01T12:00:00Z').toISOString(),
        original_path: '/path/to/test.jpg'
      }
    ];
    vi.mocked(GetPhotos).mockResolvedValue(mockPhotos as any);
    
    render(<App />);
    
    const dateText = await screen.findByText('1.1.2023');
    expect(dateText).toBeDefined();
  });
});
