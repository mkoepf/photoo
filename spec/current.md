# Current Status & Handover Context

## Project State
- **Core Engine:** Native Go/Wails desktop application.
- **Database:** SQLite (via `modernc.org/sqlite`) with schema for `photos` and `metadata_history`.
- **Import Logic:** 
    - Deduplication using SHA-256 hashes.
    - Standardized naming: `YYYY-MM-DD_HH-mm-ss_[suffix].ext`.
    - Library isolation: Files are copied to a managed `library/` folder.
- **Metadata:** 
    - EXIF extraction using `github.com/rwcarlsen/goexif/exif`.
    - Google Photos sidecar support (extracts true `photoTakenTime` from `.json` files).
- **UI:** 
    - React/TypeScript frontend with a virtualized timeline grid.
    - On-the-fly thumbnail generation via custom Wails asset handler (`/thumbnail/[filename]`).
- **Test Data:** Current `test_data/` has been successfully imported into the `library/` folder.

## Completed Tasks
- [x] Wails project scaffolded and moved to root.
- [x] Database schema initialized and persistent.
- [x] Library manager with deduplication and unique naming.
- [x] Google Photos JSON sidecar metadata extraction.
- [x] On-the-fly JPG/PNG thumbnailing.
- [x] "Import Folder" dialog and progress events.
- [x] `AfterTool` hook to run `scripts/check.sh` after file modifications for immediate feedback (Requires session restart to activate).

## Next Steps (Backlog)
1. **HEIC Thumbnails:** Currently, thumbnails for HEIC files fail because `imaging` (pure Go) doesn't support them. Need to integrate `libheif` or a similar solution.
2. **Metadata Writing:** Implement `internal/library/Manager.UpdateMetadata` to write changes back to the file's EXIF.
3. **Thumbnail Caching:** On-the-fly generation is slow for large libraries; implement a persistent thumbnail cache (e.g., in `.photoo/thumbnails`).
4. **Metadata History UI:** Create a view to browse the `metadata_history` table and revert changes.
5. **Video Support:** Extend the scanner and library logic to handle MP4/MOV files.

## Technical Notes for AI Agent
- **Main Entry:** `main.go` and `app.go`.
- **DB Operations:** `internal/db/db.go` and `internal/models/photo.go`.
- **Core Logic:** `internal/library/library.go` (Imports) and `internal/exif/exif.go` (Metadata).
- **Frontend:** `frontend/src/App.tsx`.
- **Wails Bindings:** Always run `wails generate module` after changing `App` struct methods.
- **Run Command:** `wails dev` to start with hot-reload (requires GUI environment).
