# Photoo Features

This document describes the key features of Photoo, how to try them out manually, and how to verify they are working correctly.

## 1. Desktop User Interface
Photoo provides a cross-platform desktop application built with Wails and React.

### How to try it out
1. Ensure you have the Wails CLI installed.
2. Run the application in development mode:
   ```bash
   wails dev
   ```
3. The Photoo window should appear.

### How to confirm it works
- The application window opens without errors.
- You can see the "Photoo" logo and navigation elements.
- The interface is responsive and interactive.

## 2. Photo Importing
Imports photos from any source folder into the managed Photoo library.

### How to try it out
1. Launch the app using `wails dev`.
2. Click the "Import" button (or "Select Folder" if available in the UI).
3. Select a folder containing images (e.g., `test_data/source_digital_camera`).
4. Wait for the import to complete.

### How to confirm it works
- The photo grid updates with new images after the import finishes.
- Check the `library/` directory in the project root; it should now contain copies of the imported files.
- The files in `library/` are renamed to a standardized format (`YYYY-MM-DD_HH-mm-ss.ext`).

## 3. Duplicate Detection
Prevents the same photo from being imported multiple times by comparing file hashes (SHA-256).

### How to try it out
1. Import a folder containing images.
2. Attempt to import the exact same folder again.

### How to confirm it works
- The application should finish importing quickly without adding new items to the grid.
- No new files should appear in the `library/` directory for the duplicate images.

## 4. Metadata Extraction (EXIF & Sidecars)
Automatically extracts date taken, camera model, and GPS coordinates from image files or supplemental JSON files (Google Photos style).

### How to try it out
1. Import the `test_data/source_google_photos` directory. This folder contains images and `.supplemental-metadata.json` files.
2. Import an image with standard EXIF data (e.g., from `test_data/source_digital_camera`).

### How to confirm it works
- In the Photoo UI, verify that the date shown below the photo matches the capture time (or the time in the JSON sidecar).
- *Note: While camera model and GPS coordinates are extracted and stored in the database, they are not yet displayed in the current version of the UI.*

## 5. Automatic Organization & Renaming
Organizes the library by renaming files based on their capture time to ensure a consistent structure.

### How to try it out
1. Import photos with various original filenames (e.g., `IMG_1234.JPG`, `DSC001.jpg`).

### How to confirm it works
- Navigate to the `library/` folder.
- All files follow the pattern `YYYY-MM-DD_HH-mm-ss.ext`.
- If two photos have the exact same timestamp, they are handled by adding a counter (e.g., `..._1.JPG`).

## 6. On-the-fly Thumbnails
Generates small versions of photos dynamically for the user interface.

### How to try it out
1. Import several high-resolution JPG images.
2. Browse the library in the UI.

### How to confirm it works
- The photo grid displays clear, resized versions of your photos.
- Inspect the network traffic in the Wails developer tools (F12) and look for requests to `/thumbnail/...`.
- *Note: On-the-fly thumbnail generation for HEIC files is currently not supported; these will display a placeholder or error.*

## 7. Command-Line Import Tool
A dedicated CLI utility for testing the import logic without launching the full UI.

### How to try it out
1. Run the test-import tool:
   ```bash
   go run cmd/test-import/main.go
   ```

### How to confirm it works
- The console output shows each file being scanned and imported.
- Success messages display the new library filename and extracted date.
- The `photoo.db` and `library/` directory are updated accordingly.
