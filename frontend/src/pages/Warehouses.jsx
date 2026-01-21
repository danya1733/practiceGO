import React, { useState, useEffect } from 'react';
import { Plus, MapPin } from 'lucide-react';
import { api } from '../services/api';
import { randomData } from '../utils/random';
import Modal from '../components/ui/Modal';
import MagicButton from '../components/MagicButton';
import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';

const Warehouses = () => {
    const [warehouses, setWarehouses] = useState([]);
    const [loading, setLoading] = useState(true);
    const [isModalOpen, setIsModalOpen] = useState(false);

    // Form State
    const [address, setAddress] = useState('');

    const fetchWarehouses = async () => {
        try {
            const data = await api.getWarehouses();
            setWarehouses(data || []);
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchWarehouses();
    }, []);

    const handleMagicFill = () => {
        const data = randomData.generateWarehouse();
        setAddress(data.address);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (!address) return;

        try {
            await api.createWarehouse({ address });
            setIsModalOpen(false);
            setAddress('');
            fetchWarehouses();
        } catch (error) {
            alert('Error creating warehouse');
        }
    };

    return (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="page-transition">
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem' }}>
                <h1 style={{ fontSize: '2rem', background: 'var(--secondary-gradient)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
                    Склады
                </h1>
                <button className="btn btn-primary" onClick={() => setIsModalOpen(true)}>
                    <Plus size={20} /> Добавить склад
                </button>
            </div>

            {loading ? (
                <div>Загрузка...</div>
            ) : (
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))', gap: '1.5rem' }}>
                    {warehouses.map((w) => (
                        <Link to={`/warehouses/${w.id}`} key={w.id} style={{ textDecoration: 'none', color: 'inherit' }}>
                            <motion.div
                                whileHover={{ y: -5 }}
                                className="glass-panel"
                                style={{ padding: '1.5rem', height: '100%', cursor: 'pointer' }}
                            >
                                <div style={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', marginBottom: '1rem' }}>
                                    <div style={{ padding: '0.75rem', background: 'rgba(6, 182, 212, 0.1)', borderRadius: '12px', color: '#06b6d4' }}>
                                        <MapPin size={24} />
                                    </div>
                                </div>
                                <h3 style={{ fontSize: '1.2rem', marginBottom: '0.5rem' }}>Склад</h3>
                                <p style={{ color: 'var(--text-secondary)' }}>{w.address}</p>
                                <div style={{ marginTop: '1rem', fontSize: '0.8rem', color: '#666' }}>
                                    ID: {w.id.slice(0, 8)}...
                                </div>
                            </motion.div>
                        </Link>
                    ))}
                </div>
            )}

            {/* Create Modal */}
            <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title="Новый склад">
                <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
                    <div>
                        <label style={{ display: 'block', marginBottom: '0.5rem', color: 'var(--text-secondary)' }}>Адрес</label>
                        <input
                            value={address}
                            onChange={(e) => setAddress(e.target.value)}
                            placeholder="Введите адрес склада"
                            required
                        />
                    </div>

                    <div style={{ display: 'flex', gap: '1rem', marginTop: '1rem' }}>
                        <div style={{ flex: 1 }}>
                            <MagicButton onFill={handleMagicFill} />
                        </div>
                        <button type="submit" className="btn btn-primary" style={{ flex: 1 }}>Создать</button>
                    </div>
                </form>
            </Modal>
        </motion.div>
    );
};

export default Warehouses;
