import React, { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import { Link } from 'react-router-dom';
import { Warehouse, Package, TrendingUp, ArrowRight } from 'lucide-react';
import { api } from '../services/api';

const DashboardCard = ({ title, value, icon: Icon, to, color }) => (
    <Link to={to} style={{ textDecoration: 'none', color: 'inherit' }}>
        <motion.div
            whileHover={{ y: -5, scale: 1.02 }}
            className="glass-panel"
            style={{ padding: '2rem', height: '100%', position: 'relative', overflow: 'hidden' }}
        >
            <div style={{ position: 'absolute', right: '-20px', top: '-20px', opacity: 0.1, transform: 'rotate(15deg)' }}>
                <Icon size={120} />
            </div>

            <div style={{ display: 'flex', alignItems: 'center', gap: '1rem', marginBottom: '1rem' }}>
                <div style={{ padding: '1rem', background: `rgba(${color}, 0.2)`, borderRadius: '12px', color: `rgb(${color})` }}>
                    <Icon size={24} />
                </div>
                <h3 style={{ margin: 0 }}>{title}</h3>
            </div>

            <div style={{ fontSize: '2.5rem', fontWeight: 'bold', marginBottom: '1rem' }}>
                {value}
            </div>

            <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', color: `rgb(${color})`, fontSize: '0.9rem', fontWeight: 500 }}>
                Подробнее <ArrowRight size={16} />
            </div>
        </motion.div>
    </Link>
);

const Dashboard = () => {
    const [stats, setStats] = useState({ warehouses: 0, products: 0 });

    useEffect(() => {
        // Quick summary fetch
        Promise.all([
            api.getWarehouses().then(d => d ? d.length : 0),
            api.getProducts().then(d => d ? d.length : 0)
        ]).then(([w, p]) => setStats({ warehouses: w, products: p }));
    }, []);

    return (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="page-transition">
            <div style={{ textAlign: 'center', marginBottom: '3rem', padding: '2rem 0' }}>
                <h1 style={{ fontSize: '3rem', marginBottom: '1rem', background: 'var(--primary-gradient)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
                    Добро пожаловать в WMS Pro
                </h1>
                <p style={{ color: 'var(--text-secondary)', fontSize: '1.2rem', maxWidth: '600px', margin: '0 auto' }}>
                    Премиальное решение для управления складом. Отслеживайте запасы, управляйте товарами и анализируйте показатели со стилем.
                </p>
            </div>

            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '2rem' }}>
                <DashboardCard
                    title="Склады"
                    value={stats.warehouses}
                    icon={Warehouse}
                    to="/warehouses"
                    color="6, 182, 212"
                />
                <DashboardCard
                    title="Товары"
                    value={stats.products}
                    icon={Package}
                    to="/products"
                    color="244, 63, 94"
                />
                <DashboardCard
                    title="Аналитика"
                    value="Открыть"
                    icon={TrendingUp}
                    to="/analytics"
                    color="139, 92, 246"
                />
            </div>
        </motion.div>
    );
};

export default Dashboard;
