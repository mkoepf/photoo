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
    - Write modified metadata back to the file's EXIF data.
