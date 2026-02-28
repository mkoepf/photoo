import React from 'react'
import {createRoot} from 'react-dom/client'
import './style.css'
import App from './App'

console.log("Photoo frontend initializing...");
const container = document.getElementById('root')

if (container) {
    container.style.display = "block";
    container.style.minHeight = "100vh";
} else {
    console.error("Critical: #root container not found in document!");
}

const root = createRoot(container!)

root.render(
    <React.StrictMode>
        <App/>
    </React.StrictMode>
)
