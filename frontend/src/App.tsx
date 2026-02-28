import {useState, useEffect, Component, ErrorInfo, ReactNode} from 'react';
import './App.css';
import {GetPhotos, SelectFolder, ImportFromFolder, UpdatePhotoDate} from "../wailsjs/go/main/App";
import {models} from "../wailsjs/go/models";

interface Props {
    children?: ReactNode;
}

interface State {
    hasError: boolean;
    error?: Error;
}

class ErrorBoundary extends Component<Props, State> {
    public state: State = {
        hasError: false
    };

    public static getDerivedStateFromError(error: Error): State {
        return { hasError: true, error };
    }

    public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
        console.error("Uncaught error:", error, errorInfo);
    }

    public render() {
        if (this.state.hasError) {
            return (
                <div className="error-fallback">
                    <h1>Something went wrong.</h1>
                    <pre>{this.state.error?.message}</pre>
                    <button onClick={() => window.location.reload()}>Reload App</button>
                </div>
            );
        }

        return this.props.children;
    }
}

function App() {
    const [photos, setPhotos] = useState<models.Photo[]>([]);
    const [isImporting, setIsImporting] = useState(false);
    const [selectedPhoto, setSelectedPhoto] = useState<models.Photo | null>(null);
    const [isEditing, setIsEditing] = useState(false);
    const [editDate, setEditDate] = useState("");

    const safeFormatDate = (dateVal: any, local = true) => {
        try {
            const d = new Date(dateVal);
            if (isNaN(d.getTime())) return "Invalid Date";
            return local ? d.toLocaleString() : d.toLocaleDateString();
        } catch (e) {
            return "Invalid Date";
        }
    };

    const loadPhotos = () => {
        console.log("Fetching photos...");
        GetPhotos().then(data => {
            console.log("Photos received:", data);
            setPhotos(data);
        }).catch(err => {
            console.error("Failed to fetch photos:", err);
        });
    };

    useEffect(() => {
        loadPhotos();
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
                loadPhotos();
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
            loadPhotos();
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
        <ErrorBoundary>
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
            </div>
        </ErrorBoundary>
    )
}

export default App
