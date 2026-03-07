import {useState, useEffect} from 'react';
import './App.css';
import {GetPhotosPaged, SelectFolder, ImportFromFolder, UpdatePhotoDate, LogFrontendError, LogUIState} from "../wailsjs/go/main/App";
import {models} from "../wailsjs/go/models";

// Declare global Events interface for Wails runtime
declare global {
    interface Window {
        runtime: {
            EventsOn: (eventName: string, callback: (...data: any) => void) => () => void;
        }
    }
}

const PAGE_SIZE = 100;

function App() {
    const [photos, setPhotos] = useState<models.Photo[]>([]);
    const [page, setPage] = useState(0);
    const [hasMore, setHasMore] = useState(true);
    const [isImporting, setIsImporting] = useState(false);
    const [importStatus, setImportStatus] = useState({
        total: 0,
        current: 0,
        imported: 0,
        duplicates: 0,
        errors: 0,
        lastPath: "",
        isVisible: false
    });
    const [selectedPhoto, setSelectedPhoto] = useState<models.Photo | null>(null);
    const [isEditing, setIsEditing] = useState(false);
    const [editDate, setEditDate] = useState("");

    const safeFormatDate = (dateVal: any, local = true) => {
        try {
            if (!dateVal) return "Unknown Date";
            const d = new Date(dateVal);
            if (isNaN(d.getTime())) return "Invalid Date";
            return local ? d.toLocaleString() : d.toLocaleDateString();
        } catch (e) {
            return "Invalid Date";
        }
    };

    const loadPhotos = (reset = false) => {
        const newPage = reset ? 0 : page;
        const offset = reset ? 0 : photos.length;
        
        console.log(`Fetching photos offset=${offset}, limit=${PAGE_SIZE}...`);
        GetPhotosPaged(offset, PAGE_SIZE).then(data => {
            console.log("Photos received:", data);
            if (Array.isArray(data)) {
                if (reset) {
                    setPhotos(data);
                } else {
                    setPhotos(prev => [...prev, ...data]);
                }
                setHasMore(data.length === PAGE_SIZE);
                setPage(newPage + 1);
                LogUIState(`[AUTO] Loaded ${data.length} photos. Total: ${reset ? data.length : photos.length + data.length}`);
            } else {
                console.error("Received non-array data for photos:", data);
                if (reset) setPhotos([]);
            }
        }).catch(err => {
            console.error("Failed to fetch photos:", err);
        });
    };

    useEffect(() => {
        loadPhotos(true);

        // --- Automation: Command Listener (The "Hands") ---
        if (window.runtime) {
            window.runtime.EventsOn("automation:command", (payload: any) => {
                console.log("[AUTO] Received command:", payload);
                
                if (payload.action === "get_snapshot") {
                    const cards = document.querySelectorAll('.photo-card');
                    const images = document.querySelectorAll('img.thumbnail');
                    const total = images.length;
                    const broken = Array.from(images).filter(img => (img as HTMLImageElement).naturalWidth === 0).length;
                    const htmlSnapshot = document.querySelector('.grid')?.innerHTML || "No grid found";

                    LogUIState(JSON.stringify({
                        type: "snapshot",
                        totalCards: cards.length,
                        totalImages: total,
                        brokenImages: broken,
                        gridHtmlSample: htmlSnapshot.substring(0, 500), // First 500 chars to check structure
                        timestamp: new Date().toISOString()
                    }));
                }

                if (payload.action === "inspect_thumbnails") {
                    const images = Array.from(document.querySelectorAll('img.thumbnail'));
                    const details = images.map(img => ({
                        src: (img as HTMLImageElement).src.substring(0, 50), 
                        complete: (img as HTMLImageElement).complete,
                        naturalWidth: (img as HTMLImageElement).naturalWidth,
                        naturalHeight: (img as HTMLImageElement).naturalHeight,
                        alt: (img as HTMLImageElement).alt,
                        rect: img.getBoundingClientRect()
                    }));
                    LogUIState(JSON.stringify({
                        type: "thumbnail_inspection",
                        count: images.length,
                        details: details,
                        timestamp: new Date().toISOString()
                    }));
                }

                if (payload.action === "trigger_import") {
                    const path = payload.payload;
                    LogUIState(`[AUTO] Triggering import for path: ${path}`);
                    ImportFromFolder(path).then(count => {
                        LogUIState(`[AUTO] Import complete. Count: ${count}`);
                        loadPhotos(true);
                    }).catch(err => {
                        LogFrontendError(`[AUTO] Import failed: ${err}`);
                    });
                }
            });

            // --- Import Progress Listeners ---
            window.runtime.EventsOn("import:start", (data: any) => {
                setImportStatus({
                    total: data.total,
                    current: 0,
                    imported: 0,
                    duplicates: 0,
                    errors: 0,
                    lastPath: "Initializing...",
                    isVisible: true
                });
                setIsImporting(true);
            });

            window.runtime.EventsOn("import:progress", (data: any) => {
                setImportStatus(prev => ({
                    ...prev,
                    current: data.current,
                    imported: data.imported,
                    duplicates: data.duplicates,
                    errors: data.errors,
                    lastPath: data.lastPath
                }));
            });

            window.runtime.EventsOn("import:end", (data: any) => {
                setTimeout(() => {
                    setImportStatus(prev => ({ ...prev, isVisible: false }));
                    setIsImporting(false);
                    loadPhotos(true);
                }, 2000); // Keep visible for 2 seconds to show final stats
            });
        }

        // Automated Sanity Monitor (The "Eyes")
        const monitorInterval = setInterval(() => {
            const images = document.querySelectorAll('img.thumbnail');
            const total = images.length;
            const broken = Array.from(images).filter(img => (img as HTMLImageElement).naturalWidth === 0).length;
            
            if (total > 0) {
                LogUIState(JSON.stringify({
                    type: "sanity_check",
                    totalImages: total,
                    brokenImages: broken,
                    timestamp: new Date().toISOString()
                }));
            }
        }, 5000);

        return () => clearInterval(monitorInterval);
    }, []);

    useEffect(() => {
        if (selectedPhoto) {
            try {
                // Format date for datetime-local input: YYYY-MM-DDThh:mm
                const d = new Date(selectedPhoto.date_taken);
                if (!isNaN(d.getTime())) {
                    const formatted = d.toISOString().slice(0, 16);
                    setEditDate(formatted);
                } else {
                    setEditDate("");
                }
            } catch (e) {
                setEditDate("");
            }
            setIsEditing(false);
        }
    }, [selectedPhoto]);

    const handleImport = async () => {
        try {
            const folder = await SelectFolder();
            if (folder) {
                setIsImporting(true);
                await ImportFromFolder(folder);
                setIsImporting(false);
                loadPhotos(true);
            }
        } catch (error) {
            console.error(error);
            setIsImporting(false);
        }
    };

    const handleSaveDate = async () => {
        if (!selectedPhoto) return;
        try {
            await UpdatePhotoDate(selectedPhoto.id, editDate);
            setIsEditing(false);
            loadPhotos(true);
            // Update local selection to reflect change
            setSelectedPhoto(models.Photo.createFrom({
                ...selectedPhoto,
                date_taken: new Date(editDate).toISOString()
            }));
        } catch (error) {
            console.error("Failed to update date:", error);
            alert("Failed to update date");
        }
    };

    return (
        <div id="App">
            <header className="header">
                <div className="header-content">
                    <h1>Photoo</h1>
                    <button 
                        className="btn-import" 
                        onClick={handleImport}
                        disabled={isImporting}
                    >
                        {isImporting ? 'Importing...' : 'Import Folder'}
                    </button>
                </div>
            </header>
            <div className="content-wrapper">
                <main className="main">
                    <div className="grid">
                        {photos.map(photo => (
                            <div 
                                key={photo.id} 
                                className={`photo-card ${selectedPhoto?.id === photo.id ? 'selected' : ''}`}
                                onClick={() => setSelectedPhoto(photo)}
                            >
                                <img 
                                    src={`/thumbnail/${photo.filename}`} 
                                    alt={photo.filename} 
                                    className="thumbnail"
                                    loading="lazy"
                                    onError={(e) => {
                                        const target = e.target as HTMLImageElement;
                                        if (target.src) {
                                            LogFrontendError(`Failed to load thumbnail: ${photo.filename}`);
                                        }
                                    }}
                                />
                                <div className="photo-info">
                                    <span className="date">{safeFormatDate(photo.date_taken, false)}</span>
                                </div>
                            </div>
                        ))}
                        {photos.length === 0 && !isImporting && (
                            <div className="empty-state">No photos imported yet.</div>
                        )}
                    </div>
                    {hasMore && photos.length > 0 && (
                        <div className="load-more-container">
                            <button className="btn-load-more" onClick={() => loadPhotos()}>Load More</button>
                        </div>
                    )}
                </main>
                {selectedPhoto && (
                    <aside className="sidebar">
                        <div className="sidebar-header">
                            <h2>Metadata</h2>
                            <button className="btn-close" onClick={() => setSelectedPhoto(null)}>×</button>
                        </div>
                        <div className="sidebar-content">
                            <div className="meta-item">
                                <label>Filename</label>
                                <span>{selectedPhoto.filename}</span>
                            </div>
                            <div className="meta-item">
                                <label>Date Taken</label>
                                {isEditing ? (
                                    <div className="edit-group">
                                        <input 
                                            type="datetime-local" 
                                            value={editDate}
                                            onChange={(e) => setEditDate(e.target.value)}
                                        />
                                        <div className="edit-actions">
                                            <button className="btn-save" onClick={handleSaveDate}>Save</button>
                                            <button className="btn-cancel" onClick={() => setIsEditing(false)}>Cancel</button>
                                        </div>
                                    </div>
                                ) : (
                                    <div className="display-group">
                                        <span>{safeFormatDate(selectedPhoto.date_taken)}</span>
                                        <button className="btn-edit" onClick={() => setIsEditing(true)}>Edit</button>
                                    </div>
                                )}
                            </div>
                            <div className="meta-item">
                                <label>Camera</label>
                                <span>{selectedPhoto.camera_model || 'Unknown'}</span>
                            </div>
                            {selectedPhoto.latitude !== undefined && selectedPhoto.longitude !== undefined && selectedPhoto.latitude !== 0 && selectedPhoto.longitude !== 0 && (
                                <div className="meta-item">
                                    <label>Location</label>
                                    <span>{selectedPhoto.latitude.toFixed(4)}, {selectedPhoto.longitude.toFixed(4)}</span>
                                </div>
                            )}
                            <div className="meta-item">
                                <label>Original Path</label>
                                <span className="path">{selectedPhoto.original_path}</span>
                            </div>
                        </div>
                    </aside>
                )}
            </div>

            {importStatus.isVisible && (
                <div className="modal-overlay">
                    <div className="progress-modal">
                        <h2>Importing Photos...</h2>
                        <div className="progress-stats">
                            <span>Total: {importStatus.total}</span>
                            <span>Imported: {importStatus.imported}</span>
                            <span>Duplicates: {importStatus.duplicates}</span>
                            <span>Errors: {importStatus.errors}</span>
                        </div>
                        <div className="progress-bar-container">
                            <div 
                                className="progress-bar-fill" 
                                style={{ width: `${(importStatus.current / importStatus.total) * 100}%` }}
                            ></div>
                        </div>
                        <div className="progress-last-path">
                            {importStatus.lastPath}
                        </div>
                    </div>
                </div>
            )}
        </div>
    )
}

export default App
