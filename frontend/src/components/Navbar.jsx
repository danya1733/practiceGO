import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { LayoutDashboard, Package, Warehouse, TrendingUp } from 'lucide-react';

const Navbar = () => {
    const location = useLocation();

    const isActive = (path) => location.pathname === path;

    const links = [
        { path: '/', label: 'Дашборд', icon: LayoutDashboard },
        { path: '/warehouses', label: 'Склады', icon: Warehouse },
        { path: '/products', label: 'Товары', icon: Package },
        { path: '/analytics', label: 'Аналитика', icon: TrendingUp },
    ];

    const [isHealthy, setIsHealthy] = React.useState(true);

    React.useEffect(() => {
        const checkHealth = async () => {
            try {
                const res = await fetch('/api/health'); // Direct call or use api.healthCheck
                setIsHealthy(res.ok);
            } catch {
                setIsHealthy(false);
            }
        };
        checkHealth();
        const interval = setInterval(checkHealth, 30000); // Check every 30s
        return () => clearInterval(interval);
    }, []);

    return (
        <nav className="glass-panel" style={{
            position: 'sticky',
            top: '1rem',
            zIndex: 100,
            marginBottom: '2rem',
            padding: '1rem 2rem',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center'
        }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                <div style={{ fontWeight: 'bold', fontSize: '1.5rem', background: 'var(--primary-gradient)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
                    WMS Pro
                </div>
                <div title={isHealthy ? "Online" : "Offline"} style={{ width: '10px', height: '10px', borderRadius: '50%', background: isHealthy ? '#4ade80' : '#f87171', boxShadow: `0 0 10px ${isHealthy ? '#4ade80' : '#f87171'}` }} />
            </div>

            <div style={{ display: 'flex', gap: '1rem' }}>
                {links.map(({ path, label, icon: Icon }) => (
                    <Link
                        key={path}
                        to={path}
                        style={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: '0.5rem',
                            textDecoration: 'none',
                            color: isActive(path) ? '#fff' : 'var(--text-secondary)',
                            padding: '0.5rem 1rem',
                            borderRadius: '8px',
                            background: isActive(path) ? 'rgba(255,255,255,0.1)' : 'transparent',
                            transition: 'all 0.3s ease'
                        }}
                    >
                        <Icon size={18} />
                        {label}
                    </Link>
                ))}
            </div>
        </nav>
    );
};

export default Navbar;
