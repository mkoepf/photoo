# Architectural Decisions

## 1. Desktop Environment
**Decision:** Native Desktop Application.
**Rationale:** Photoo is designed to run locally on the user's machine. A native application provides a more integrated experience (e.g., dock/taskbar icon) and easier access to local file systems.
**Technology:** **Wails v2** (Go backend with a native webview frontend).

## 2. Programming Language & Tech Stack
**Decision:** Go (Backend), React + Tailwind CSS (Frontend).
**Rationale:** Go is chosen for its performance, excellent concurrency for file processing, and strong standard library. React/Tailwind allows for a modern, responsive, and "simple viewer" UI.
**Database:** **SQLite** (via `modernc.org/sqlite`). A single-file, zero-dependency database ideal for local metadata storage.

## 3. Metadata & Image Handling
**Decision:** On-the-fly thumbnail generation; direct EXIF modification with history tracking.
**Rationale:** 
- **On-the-fly thumbnails:** Saves disk space and simplifies the initial import, with the option to cache for performance later.
- **Direct EXIF:** Ensures the imported files are "self-contained" and usable by other applications.
- **History Tracking:** A dedicated SQLite table will store every change made to a photo's metadata, allowing for reverts or audit trails.

## 4. Library Structure
**Decision:** "One Big Folder" with standardized naming.
**Rationale:** Simplifies file management. Filenames follow the pattern: `YYYY-MM-DD_HH-mm-ss_[suffix].ext`.

## 5. Scope & Scale
**Decision:** 
- **Scale:** Optimized for up to 50,000 photos.
- **AI:** No AI (face/object detection) in the MVP.
- **Formats:** JPEG, PNG, and **HEIC** (Apple). RAW and Video support will be added as needed.
- **HEIC Handling:** Use a high-performance library (e.g., `libheif` via CGO or a dedicated Go wrapper) to ensure seamless import and thumbnail generation for Apple device photos.
