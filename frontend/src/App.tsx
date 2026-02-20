import {useState, useEffect} from 'react';
import './App.css';
import {GetPhotos, SelectFolder, ImportFromFolder} from "../wailsjs/go/main/App";
import {models} from "../wailsjs/go/models";

function App() {
    const [photos, setPhotos] = useState<models.Photo[]>([]);
    const [isImporting, setIsImporting] = useState(false);

    const loadPhotos = () => {
        GetPhotos().then(setPhotos).catch(console.error);
    };

    useEffect(() => {
        loadPhotos();
    }, []);

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
            <main className="main">
                <div className="grid">
                    {photos.map(photo => (
                        <div key={photo.id} className="photo-card">
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
        </div>
    )
}

export default App
