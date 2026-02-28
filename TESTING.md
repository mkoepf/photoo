# Testing Photoo

This document provides instructions on how to verify the current state of the project using automated and manual methods.

## 1. Automated Quality Suite
The project includes a comprehensive check script that validates the frontend build, Go formatting, backend tests, and security vulnerabilities.

**To run the full suite:**
```bash
./scripts/check.sh
```

**Success Indicators:**
- `Checking Frontend... DONE`
- `Checking Go formatting... PASSED`
- `Running Go tests... Go tests PASSED`
- `All Photoo Quality Checks PASSED`

---

## 2. Manual CLI Import Test
If you want to test the core library management logic (hashing, deduplication, and renaming) without launching the GUI, use the test-import utility.

**To run the import tool:**
```bash
go run cmd/test-import/main.go
```

**What it does:**
- Scans the `test_data/` directory.
- Imports images into the `library/` folder.
- Populates the `photoo.db` SQLite database.
- Outputs success/error messages for each file.

---

## 3. Desktop UI Testing
To test the full application experience, including the interactive grid and folder selection.

**To start the app in development mode:**
```bash
wails dev
```

**Testing Steps:**
1. **Launch:** Verify the window opens and the "Photoo" logo appears.
2. **Import:** Click **Import Folder** and select `test_data/source_digital_camera`.
3. **Verify Grid:** Confirm that thumbnails appear in the grid with dates underneath.
4. **Deduplication:** Import the same folder again; verify that the "Importing..." state finishes quickly and no duplicate entries appear in the grid.

---

## 4. Test Data Scenarios
We have provided three specific datasets in `test_data/` to verify different logic paths:

- **source_digital_camera:** Standard JPEGs with internal EXIF data.
- **source_google_photos:** Images paired with `.supplemental-metadata.json` files. Verify that the "Date Taken" matches the JSON timestamp rather than the file modification time.
- **source_icloud:** Contains `.HEIC` files. Verify they are imported and renamed correctly (note: thumbnails for these are not yet supported).
