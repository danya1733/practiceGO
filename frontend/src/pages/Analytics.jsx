import React, { useState, useEffect } from 'react';
import { api } from '../services/api';
import { motion } from 'framer-motion';
import { TrendingUp, Award } from 'lucide-react';

const Analytics = () => {
    const [topWarehouses, setTopWarehouses] = useState([]);

    useEffect(() => {
        api.getTopWarehouses().then(data => setTopWarehouses(data || []));
    }, []);

    return (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="page-transition">
            <h1 style={{ fontSize: '2rem', marginBottom: '2rem', display: 'flex', alignItems: 'center', gap: '1rem' }}>
                <TrendingUp /> Аналитическая панель
            </h1>

            <div className="glass-panel" style={{ padding: '2rem' }}>
                <h2 style={{ marginBottom: '1.5rem', display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                    <Award color="#f59e0b" /> Топ Склады
                </h2>

                <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                    {topWarehouses.map((w, index) => (
                        <div key={w.warehouse_id} style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                            <div style={{
                                width: '30px', height: '30px',
                                borderRadius: '50%', background: index === 0 ? '#f59e0b' : index === 1 ? '#94a3b8' : '#b45309',
                                display: 'flex', alignItems: 'center', justifyContent: 'center', fontWeight: 'bold'
                            }}>
                                {index + 1}
                            </div>
                            <div style={{ flex: 1 }}>
                                <div style={{ marginBottom: '0.25rem' }}>{w.address}</div>
                                <div style={{ width: '100%', height: '8px', background: 'rgba(255,255,255,0.1)', borderRadius: '4px', overflow: 'hidden' }}>
                                    <motion.div
                                        initial={{ width: 0 }}
                                        animate={{ width: `${(w.total_sum / (topWarehouses[0]?.total_sum || 1)) * 100}%` }}
                                        transition={{ duration: 1, ease: 'easeOut' }}
                                        style={{ height: '100%', background: 'var(--primary-gradient)' }}
                                    />
                                </div>
                            </div>
                            <div style={{ fontWeight: 'bold', minWidth: '100px', textAlign: 'right' }}>
                                {w.total_sum.toLocaleString()} ₽
                            </div>
                        </div>
                    ))}
                    {topWarehouses.length === 0 && <p style={{ color: 'var(--text-secondary)' }}>Нет данных для аналитики.</p>}
                </div>
            </div>
        </motion.div>
    );
};

export default Analytics;
