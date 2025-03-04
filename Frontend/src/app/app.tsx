import { Route, Routes } from 'react-router-dom'
import './app.css'
import LoginPage from '../login/login.tsx'
import ItemPage from "../item/item.tsx";
import ItemDetail from "../item/item-detail.tsx";
import {EventDashboard} from "../track-point/event_dashboard.tsx";

export default function App() {
    return (
        <div className='app'>
            <div className='content'>
                <Routes>
                    <Route path='/' element={<LoginPage />}/>
                    <Route path='/login' element={<LoginPage />} />
                    <Route path='/item-page' element={<ItemPage /> }/>
                    <Route path='/item-detail' element={<ItemDetail /> }/>
                    <Route path='/event-dashboard' element={<EventDashboard /> }/>
                </Routes>
            </div>
        </div>
    )
}
