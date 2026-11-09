import React from 'react'
import {createRoot} from 'react-dom/client'
import './style.css'
import App from './App'
import { GlobalProvider } from './global'

const container = document.getElementById('root')

const root = createRoot(container!)

root.render(
    <React.StrictMode>
                <GlobalProvider>
            <App />
        </GlobalProvider>

    </React.StrictMode>
)
