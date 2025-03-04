import './index.css'
import App from './app/app.tsx'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter, HashRouter } from 'react-router-dom'
import {readEventCommonParams} from "./track-point/event_params.tsx";


export const serverUrl = "http://localhost:5173"
const root = createRoot(document.getElementById('root')!)

root.render(
  <StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </StrictMode>,
)

readEventCommonParams()
