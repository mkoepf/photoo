# Photoo - Implementation Plan

## Phase 1: Project & Database Infrastructure
1.  **Initialize Project:** Scaffold a Wails project with Go and React/Vite.
2.  **Schema Design:**
    *   `photos` Table: Current state (path, hash, date\_taken, location, etc.).
    *   `metadata_history` Table: Previous values for all fields to enable undo.
3.  **Library Service:** Implement the core Go service that manages the `library/` folder and ensures unique, date-based filenames.

## Phase 2: Import & De-duplication Engine
1.  **Scanner Service:** A background worker to walk source folders (e.g., `test_data/source_google_photos`).
2.  **HEIC Handler:** Implement support for HEIC/HEIF files (often from iPhones) to ensure they are handled alongside JPEGs.
3.  **Duplicate Detector:** SHA-256 hashing for bit-perfect file identification.
3.  **Import Logic:**
    *   Read original EXIF data.
    *   Copy unique files to the library.
    *   Rename files to `YYYY-MM-DD_HH-mm-ss_[suffix].ext`.
    *   Write normalized EXIF data back to the file.

## Phase 3: Metadata Normalization & Viewer UI
1.  **Metadata Service:** A Go worker that standardizes dates and locations across various EXIF formats.
2.  **Thumbnail API:** A custom Wails asset handler for on-the-fly, high-quality image resizing.
3.  **Timeline Grid:** A virtualized React component to display the entire library, grouped by Month/Year.
4.  **Metadata Sidebar:** A simple panel to view and edit metadata for individual or multiple photos.

## Phase 4: Validation & Refining
1.  **Consistency Tool:** Ensure the file system and database are always in sync.
2.  **Search Index:** Basic search/filter functionality for dates and camera models.
3.  **UI Polish:** Ensure the "Simple Viewer" aesthetic is maintained.

---

## Test Scenarios (Using `test_data/`)
- **Scenario 1: Cloud Sync Mess**
    - *Goal:* Verify duplicate detection across Google Photos and iCloud.
- **Scenario 2: Time Traveler**
    - *Goal:* Manually correct dates and verify both the filename and EXIF update.
- **Scenario 3: Naming & One-Big-Folder**
    - *Goal:* Ensure 100+ files are correctly imported and renamed without collisions.
