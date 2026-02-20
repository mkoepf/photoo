# Photoo

Photoo is a high-performance, native desktop photo management application designed to help you organize and de-duplicate your digital photographs from various sources (Google Photos, iCloud, Digital Cameras).

## Key Features

- **Native Performance:** Built with Go and Wails for a fast, local-first experience.
- **Intelligent Import:** Automatically de-duplicates photos using SHA-256 hashing.
- **Standardized Library:** Copies photos into a managed "One Big Folder" with unique, date-based naming (`YYYY-MM-DD_HH-mm-ss.ext`).
- **Metadata Mastery:** 
    - Extracts EXIF data (Date Taken, Camera Model, GPS).
    - Supports Google Photos sidecar `.json` files to recover missing timestamps and location data.
    - Preserves metadata history in a local SQLite database.
- **HEIC Support:** Handles modern Apple photo formats and extracts their metadata.
- **Timeline View:** A clean, responsive grid view of your entire library sorted by date.

## Tech Stack

- **Backend:** [Go](https://go.dev/) (1.25+)
- **Desktop Framework:** [Wails v2](https://wails.io/)
- **Database:** [SQLite](https://sqlite.org/) (Pure Go via `modernc.org/sqlite`)
- **Frontend:** [React](https://reactjs.org/) + [TypeScript](https://www.typescriptlang.org/) + [Tailwind CSS](https://tailwindcss.com/)
- **Image Processing:** [Imaging](https://github.com/disintegration/imaging)

## Getting Started

### Prerequisites

- Go 1.25 or later
- Node.js & NPM
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)

### Development

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/mkoepf/photoo.git
    cd photoo
    ```

2.  **Run in development mode:**
    ```bash
    wails dev
    ```
    This will start the application with hot-reloading for both the Go backend and the React frontend.

3.  **Run Quality Checks:**
    ```bash
    ./scripts/check.sh
    ```

## Project Structure

- `app.go`: Main application logic and Wails bindings.
- `internal/library/`: Core logic for file importing and de-duplication.
- `internal/exif/`: Metadata extraction and sidecar parsing.
- `internal/db/`: SQLite schema and database initialization.
- `frontend/`: React-based user interface.
- `spec/`: Vision, architectural decisions, and implementation plans.

## Roadmap

- [ ] HEIC Thumbnail generation support.
- [ ] Metadata editing and write-back to EXIF.
- [ ] Video support (MP4/MOV).
- [ ] Persistent thumbnail caching.

## License

See the LICENSE file (if applicable) for details.
