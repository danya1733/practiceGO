import React from 'react';
import { Routes, Route } from 'react-router-dom';
import Navbar from './components/Navbar';
import Dashboard from './pages/Dashboard';
import Warehouses from './pages/Warehouses';
import Products from './pages/Products';
import WarehouseDetail from './pages/WarehouseDetail';
import Analytics from './pages/Analytics';

function App() {
    return (
        <div className="container" style={{ paddingBottom: '2rem' }}>
            <Navbar />
            <Routes>
                <Route path="/" element={<Dashboard />} />
                <Route path="/warehouses" element={<Warehouses />} />
                <Route path="/warehouses/:id" element={<WarehouseDetail />} />
                <Route path="/products" element={<Products />} />
                <Route path="/analytics" element={<Analytics />} />
            </Routes>
        </div>
    );
}

export default App;
