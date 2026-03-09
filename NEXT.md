# Photoo - Next Steps

This document tracks the prioritized roadmap for the Photoo project.

## High Priority (Immediate Focus)

- [ ] **1. Date-Based Library Organization**
    - Implement logic to move or symlink imported photos into a structured hierarchy: `library/YYYY/MM/DD/`.
    - Ensure the database paths are updated automatically during this movement.
- [ ] **2. Import Progress & Status UI**
    - Add a global progress bar or modal for active import operations.
    - Show real-time statistics: "Processed X of Y", "Skipped Z duplicates".
    - Log detailed import events to the automation bridge.
- [ ] **3. Recursive Folder Scanning**
    - Extend `ImportFromFolder` to traverse subdirectories recursively.
    - Implement a job queue or worker pool to handle massive imports without blocking the UI.
    - Ensure self-test coverage for deep directory structures.

## Mid Priority

- [ ] **4. EXIF Metadata Write-back**
    - Implement `go-exif` logic to write edited timestamps back to the original file's EXIF header.
    - Automatically rename the file in the library if the date changes (maintaining the `YYYY-MM-DD_HH-MM-SS` convention).
- [ ] **5. Intelligent Duplicate Detection**
    - Calculate SHA-256 hashes during import to skip exact file duplicates.
    - Research/Implement Perceptual Hashing (PHash) to detect visually similar photos.
- [ ] **6. Full-Screen Photo Viewer**
    - Implement a lightbox view when clicking a photo in the grid.
    - Add keyboard navigation (Left/Right arrows) and basic zoom controls.
- [ ] **7. Video Support (.MP4, .MOV)**
    - Extend backend managers to recognize video formats.
    - Use `ffmpeg` (if available) or native Go libraries to extract frame thumbnails.

## Low Priority (Roadmap)

- [ ] **8. Advanced Filtering & Search**
    - Add search by camera model, date range, or location.
    - Implement "Smart Collections" (e.g., "Photos from last Summer").
- [ ] **9. Export Functionality**
    - Allow users to select multiple photos and export them to an external folder.
    - Add options for resizing or stripping metadata on export.
- [ ] **10. Selective Deletion & Pruning**
    - Implement a "Trash" workflow to mark photos for removal.
    - Permanently delete files from disk and records from SQLite upon confirmation.

## Technical Debt / Completed
- [x] **Fix Thumbnail Visibility (Asset Routing):** Resolved the `431` error by switching to Base64 serving and removing the Vite proxy.
- [x] **Autonomous Testing Mandate:** Updated `GEMINI.md` to ensure all future UI features include automation hooks.
- [x] **Large Library Optimization (Memory/Performance):**
    - Implemented backend worker pool (semaphore) for thumbnail generation to prevent memory spikes.
    - Switched from base64 to HTTP URLs for thumbnails to leverage browser caching.
    - Added frontend pagination (Load More) and backend paged queries to handle 7000+ photos.
    - Improved `check.sh` robustness for ARM64 environments.
