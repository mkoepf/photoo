import {useState, useEffect} from 'react';
import './App.css';
import {GetPhotos, SelectFolder, ImportFromFolder, UpdatePhotoDate} from "../wailsjs/go/main/App";
import {models} from "../wailsjs/go/models";

function App() {
    const [photos, setPhotos] = useState<models.Photo[]>([]);
    const [isImporting, setIsImporting] = useState(false);
    const [selectedPhoto, setSelectedPhoto] = useState<models.Photo | null>(null);
    const [isEditing, setIsEditing] = useState(false);
    const [editDate, setEditDate] = useState("");

    const loadPhotos = () => {
        GetPhotos().then(setPhotos).catch(console.error);
    };

    useEffect(() => {
        loadPhotos();
    }, []);

    useEffect(() => {
        if (selectedPhoto) {
            // Format date for datetime-local input: YYYY-MM-DDThh:mm
            const d = new Date(selectedPhoto.date_taken);
            const formatted = d.toISOString().slice(0, 16);
            setEditDate(formatted);
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
                                    <span className="date">{new Date(photo.date_taken).toLocaleDateString()}</span>
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
                                        <span>{new Date(selectedPhoto.date_taken).toLocaleString()}</span>
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
    )
}

export default App
