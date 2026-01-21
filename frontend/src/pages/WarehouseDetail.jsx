import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { api } from '../services/api';
import { randomData } from '../utils/random';
import Modal from '../components/ui/Modal';
import MagicButton from '../components/MagicButton';
import PurchaseModal from '../components/PurchaseModal';
import { motion } from 'framer-motion';
import { Plus, ShoppingCart, DollarSign, Edit, Check } from 'lucide-react';

const WarehouseDetail = () => {
    const { id } = useParams();
    const [inventory, setInventory] = useState([]);
    const [products, setProducts] = useState([]);
    const [analytics, setAnalytics] = useState(null);
    const [loading, setLoading] = useState(true);

    // Modals
    const [isInventoryModalOpen, setInventoryModalOpen] = useState(false);
    const [isPurchaseModalOpen, setPurchaseModalOpen] = useState(false);

    // Editing State
    const [editingItem, setEditingItem] = useState(null);
    const [editForm, setEditForm] = useState({ quantity: 0, discount: 0 });

    // Inventory Form
    const [invForm, setInvForm] = useState({
        product_id: '',
        quantity: '',
        price: '',
        discount: ''
    });

    const fetchData = async () => {
        try {
            const [invData, prodData, analyticsData] = await Promise.all([
                api.getWarehouseProducts(id),
                api.getProducts(),
                api.getWarehouseAnalytics(id)
            ]);
            setInventory(invData || []);
            setProducts(prodData || []);
            setAnalytics(analyticsData);
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchData();
    }, [id]);

    const handleMagicFillInventory = () => {
        const randomProduct = products[Math.floor(Math.random() * products.length)];
        const data = randomData.generateInventory();

        setInvForm({
            product_id: randomProduct ? randomProduct.id : '',
            quantity: data.quantity,
            price: data.price,
            discount: data.discount
        });
    };

    const handleAddInventory = async (e) => {
        e.preventDefault();
        try {
            await api.createInventory({
                warehouse_id: id,
                product_id: invForm.product_id,
                quantity: parseInt(invForm.quantity),
                price: parseFloat(invForm.price),
                discount: parseFloat(invForm.discount)
            });
            setInventoryModalOpen(false);
            setInvForm({ product_id: '', quantity: '', price: '', discount: '' });
            fetchData();
        } catch (error) {
            alert('Error adding inventory');
        }
    };

    const startEditing = (item) => {
        setEditingItem(item.id);
        setEditForm({ quantity: item.quantity, discount: item.discount });
    };

    const saveEdit = async (item) => {
        try {
            if (editForm.quantity !== item.quantity) {
                await api.updateQuantity(id, item.product_id, editForm.quantity);
            }
            if (editForm.discount !== item.discount) {
                await api.updateDiscount(id, item.product_id, editForm.discount);
            }
            setEditingItem(null);
            fetchData();
        } catch (e) {
            alert('Ошибка обновления: ' + e.message);
        }
    };

    return (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="page-transition">
            <div style={{ marginBottom: '2rem' }}>
                <h1 style={{ fontSize: '2rem', marginBottom: '0.5rem' }}>Детали склада</h1>
                <p style={{ color: 'var(--text-secondary)' }}>ID: {id}</p>
            </div>

            {/* Stats Cards */}
            {analytics && (
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))', gap: '1rem', marginBottom: '2rem' }}>
                    <div className="glass-panel" style={{ padding: '1.5rem', display: 'flex', alignItems: 'center', gap: '1rem' }}>
                        <div style={{ background: 'rgba(99, 102, 241, 0.2)', padding: '10px', borderRadius: '50%' }}>
                            <DollarSign color="#818cf8" />
                        </div>
                        <div>
                            <h3 style={{ fontSize: '0.9rem', color: 'var(--text-secondary)' }}>Общая выручка</h3>
                            <p style={{ fontSize: '1.5rem', fontWeight: 'bold' }}>{analytics.total_sum.toLocaleString()} ₽</p>
                        </div>
                    </div>
                </div>
            )}

            {/* Actions */}
            <div style={{ display: 'flex', gap: '1rem', marginBottom: '2rem' }}>
                <button className="btn btn-primary" onClick={() => setInventoryModalOpen(true)}>
                    <Plus size={18} /> Добавить позицию
                </button>
                <button className="btn" style={{ background: 'var(--accent-gradient)' }} onClick={() => setPurchaseModalOpen(true)}>
                    <ShoppingCart size={18} /> Симуляция покупки
                </button>
            </div>

            {/* Inventory List */}
            <h2 style={{ marginBottom: '1rem' }}>Текущий инвентарь</h2>
            <div className="glass-panel" style={{ overflow: 'hidden' }}>
                <table style={{ width: '100%', borderCollapse: 'collapse' }}>
                    <thead>
                        <tr style={{ borderBottom: '1px solid var(--glass-border)', textAlign: 'left' }}>
                            <th style={{ padding: '1rem' }}>Товар</th>
                            <th style={{ padding: '1rem' }}>Кол-во</th>
                            <th style={{ padding: '1rem' }}>Цена</th>
                            <th style={{ padding: '1rem' }}>Скидка %</th>
                            <th style={{ padding: '1rem' }}>Действия</th>
                        </tr>
                    </thead>
                    <tbody>
                        {inventory.map((item) => (
                            <tr key={item.id} style={{ borderBottom: '1px solid var(--glass-border)' }}>
                                <td style={{ padding: '1rem' }}>
                                    <div style={{ fontWeight: 'bold' }}>{item.product?.name || 'Неизвестный товар'}</div>
                                    <div style={{ fontSize: '0.8rem', color: 'var(--text-secondary)' }}>{item.product?.barcode}</div>
                                </td>
                                <td style={{ padding: '1rem' }}>
                                    {editingItem === item.id ? (
                                        <input
                                            type="number"
                                            value={editForm.quantity}
                                            onChange={(e) => setEditForm({ ...editForm, quantity: e.target.value })}
                                            style={{ width: '80px' }}
                                        />
                                    ) : item.quantity}
                                </td>
                                <td style={{ padding: '1rem' }}>{item.price} ₽</td>
                                <td style={{ padding: '1rem' }}>
                                    {editingItem === item.id ? (
                                        <input
                                            type="number"
                                            value={editForm.discount}
                                            onChange={(e) => setEditForm({ ...editForm, discount: e.target.value })}
                                            style={{ width: '80px' }}
                                        />
                                    ) : item.discount}
                                </td>
                                <td style={{ padding: '1rem' }}>
                                    {editingItem === item.id ? (
                                        <button onClick={() => saveEdit(item)} className="btn-icon" style={{ background: 'transparent', border: 'none', color: '#4ade80', cursor: 'pointer' }}>
                                            <Check size={18} />
                                        </button>
                                    ) : (
                                        <button onClick={() => startEditing(item)} className="btn-icon" style={{ background: 'transparent', border: 'none', color: '#60a5fa', cursor: 'pointer' }}>
                                            <Edit size={18} />
                                        </button>
                                    )}
                                </td>
                            </tr>
                        ))}
                        {inventory.length === 0 && (
                            <tr>
                                <td colSpan="5" style={{ padding: '2rem', textAlign: 'center', color: 'var(--text-secondary)' }}>
                                    Нет данных. Добавьте товары!
                                </td>
                            </tr>
                        )}
                    </tbody>
                </table>
            </div>

            {/* Modals */}
            <Modal isOpen={isInventoryModalOpen} onClose={() => setInventoryModalOpen(false)} title="Добавить позицию">
                <form onSubmit={handleAddInventory} style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                    <select
                        value={invForm.product_id}
                        onChange={e => setInvForm({ ...invForm, product_id: e.target.value })}
                        required
                        style={{ background: 'rgba(0,0,0,0.3)', color: 'white' }}
                    >
                        <option value="">Выберите товар...</option>
                        {products.map(p => (
                            <option key={p.id} value={p.id}>{p.name}</option>
                        ))}
                    </select>

                    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: '1rem' }}>
                        <input type="number" placeholder="Кол-во" value={invForm.quantity} onChange={e => setInvForm({ ...invForm, quantity: e.target.value })} required />
                        <input type="number" placeholder="Цена" step="0.01" value={invForm.price} onChange={e => setInvForm({ ...invForm, price: e.target.value })} required />
                        <input type="number" placeholder="Скидка %" step="0.1" value={invForm.discount} onChange={e => setInvForm({ ...invForm, discount: e.target.value })} required />
                    </div>

                    <div style={{ display: 'flex', gap: '1rem', marginTop: '0.5rem' }}>
                        <div style={{ flex: 1 }}>
                            <MagicButton onFill={handleMagicFillInventory} />
                        </div>
                        <button type="submit" className="btn btn-primary" style={{ flex: 1 }}>Добавить</button>
                    </div>
                </form>
            </Modal>

            <PurchaseModal
                isOpen={isPurchaseModalOpen}
                onClose={() => setPurchaseModalOpen(false)}
                warehouseId={id}
                inventory={inventory}
                onSuccess={fetchData}
            />

        </motion.div>
    );
};

export default WarehouseDetail;
