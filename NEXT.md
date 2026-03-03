# Photoo - Next Steps

This document tracks the immediate development goals for the Photoo project.

- [x] **1. Metadata Display & Sidebar (UI)**
    - Add a sidebar to the photo grid.
    - Show capture date, camera model, and GPS coordinates when a photo is selected.
- [x] **2. HEIC Thumbnail Support (Backend)**
    - Integrate HEIF/HEIC decoding in `thumbnail.go`.
    - Ensure iPhone photos appear correctly in the grid.
- [x] **3. Persistent Thumbnail Cache (Backend/Performance)**
    - Implement a hidden storage for pre-rendered thumbnails.
    - Improve UI responsiveness for large libraries.
- [x] **4. Metadata Editing (Backend/UI)**
    - Allow users to modify capture dates in the sidebar.
- [ ] **5. Metadata Write-back (EXIF)**
    - Implement logic to write date/location changes back to the actual file's EXIF header.
    - Ensure the library file remains the "source of truth".
- [ ] **6. Virtualized Timeline Grid**
    - Group photos by Month/Year in the UI.
    - Use virtualization (e.g., `react-window`) to handle thousands of photos smoothly.
- [ ] **7. Metadata History & Undo UI**
    - Create a view to browse the `metadata_history` table.
    - Add functionality to revert changes to a previous state.
- [ ] **8. Video Support**
    - Support importing `.mp4` and `.mov` files.
    - Implement basic thumbnail extraction for video files.

## Technical Debt / Known Issues
- [ ] **Routing in `wails dev` (Scenario A):** Thumbnail requests (`/thumbnail/...`) are being intercepted by the frontend dev server/proxy and do not reach the Go backend custom handler. Needs investigation of Wails v2 asset serving priorities in development mode.
