import React, { useState, useEffect } from 'react';
import { Plus, Package, Tag } from 'lucide-react';
import { api } from '../services/api';
import { randomData } from '../utils/random';
import Modal from '../components/ui/Modal';
import MagicButton from '../components/MagicButton';
import { motion } from 'framer-motion';

const Products = () => {
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [isModalOpen, setIsModalOpen] = useState(false);

    // Form State
    const [form, setForm] = useState({
        name: '',
        description: '',
        weight: '',
        barcode: '',
        characteristics: ''
    });

    const fetchProducts = async () => {
        try {
            const data = await api.getProducts();
            setProducts(data || []);
        } catch (error) {
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchProducts();
    }, []);

    const handleMagicFill = () => {
        const data = randomData.generateProduct();
        setForm(data);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            await api.createProduct({
                ...form,
                weight: parseFloat(form.weight)
            });
            setIsModalOpen(false);
            setForm({ name: '', description: '', weight: '', barcode: '', characteristics: '' });
            fetchProducts();
        } catch (error) {
            alert('Error creating product');
        }
    };

    // Edit State
    const [editingProduct, setEditingProduct] = useState(null);
    const [isEditModalOpen, setEditModalOpen] = useState(false);

    const openEditModal = (product) => {
        setEditingProduct(product);
        setForm({
            name: product.name,
            description: product.description,
            weight: product.weight,
            barcode: product.barcode,
            characteristics: JSON.stringify(product.characteristics || {}, null, 2)
        });
        setEditModalOpen(true);
    };

    const handleUpdate = async (e) => {
        e.preventDefault();
        try {
            await api.updateProduct(editingProduct.id, {
                ...form,
                weight: parseFloat(form.weight)
            });
            setEditModalOpen(false);
            setEditingProduct(null);
            setForm({ name: '', description: '', weight: '', barcode: '', characteristics: '' });
            fetchProducts();
        } catch (error) {
            alert('Error updating product');
        }
    };

    return (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="page-transition">
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem' }}>
                <h1 style={{ fontSize: '2rem', background: 'var(--accent-gradient)', WebkitBackgroundClip: 'text', WebkitTextFillColor: 'transparent' }}>
                    Товары
                </h1>
                <button className="btn btn-primary" onClick={() => { setEditingProduct(null); setForm({ name: '', description: '', weight: '', barcode: '', characteristics: '' }); setIsModalOpen(true); }}>
                    <Plus size={20} /> Добавить товар
                </button>
            </div>

            {loading ? (
                <div>Загрузка...</div>
            ) : (
                <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(300px, 1fr))', gap: '1.5rem' }}>
                    {products.map((p) => (
                        <motion.div
                            key={p.id}
                            whileHover={{ y: -5 }}
                            className="glass-panel"
                            style={{ padding: '1.5rem', height: '100%', cursor: 'pointer', position: 'relative' }}
                            onClick={() => openEditModal(p)}
                        >
                            <div style={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', marginBottom: '1rem' }}>
                                <div style={{ padding: '0.75rem', background: 'rgba(244, 63, 94, 0.1)', borderRadius: '12px', color: '#f43f5e' }}>
                                    <Package size={24} />
                                </div>
                            </div>
                            <h3 style={{ fontSize: '1.2rem', marginBottom: '0.5rem' }}>{p.name}</h3>
                            <p style={{ color: 'var(--text-secondary)', fontSize: '0.9rem', marginBottom: '1rem' }}>{p.description}</p>

                            <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap' }}>
                                <span style={{ fontSize: '0.8rem', padding: '0.25rem 0.5rem', background: 'rgba(255,255,255,0.1)', borderRadius: '4px' }}>
                                    {p.weight} kg
                                </span>
                                <span style={{ fontSize: '0.8rem', padding: '0.25rem 0.5rem', background: 'rgba(255,255,255,0.1)', borderRadius: '4px', display: 'flex', alignItems: 'center', gap: '4px' }}>
                                    <Tag size={12} /> {p.barcode}
                                </span>
                            </div>
                        </motion.div>
                    ))}
                </div>
            )}

            {/* Create Modal */}
            <Modal isOpen={isModalOpen} onClose={() => setIsModalOpen(false)} title="Новый товар">
                <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                    <input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} placeholder="Название" required />
                    <textarea value={form.description} onChange={e => setForm({ ...form, description: e.target.value })} placeholder="Описание" rows={3} />
                    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                        <input type="number" step="0.1" value={form.weight} onChange={e => setForm({ ...form, weight: e.target.value })} placeholder="Вес (кг)" required />
                        <input value={form.barcode} onChange={e => setForm({ ...form, barcode: e.target.value })} placeholder="Штрихкод" required />
                    </div>
                    <textarea
                        value={form.characteristics}
                        onChange={e => setForm({ ...form, characteristics: e.target.value })}
                        placeholder='Характеристики (JSON) например {"цвет": "красный"}'
                        rows={3}
                        style={{ fontFamily: 'monospace', fontSize: '0.9rem' }}
                    />

                    <div style={{ display: 'flex', gap: '1rem', marginTop: '0.5rem' }}>
                        <div style={{ flex: 1 }}>
                            <MagicButton onFill={handleMagicFill} />
                        </div>
                        <button type="submit" className="btn btn-primary" style={{ flex: 1 }}>Создать</button>
                    </div>
                </form>
            </Modal>

            {/* Edit Modal */}
            <Modal isOpen={isEditModalOpen} onClose={() => setEditModalOpen(false)} title="Редактировать товар">
                <form onSubmit={handleUpdate} style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                    <input value={form.name} onChange={e => setForm({ ...form, name: e.target.value })} placeholder="Название" required />
                    <textarea value={form.description} onChange={e => setForm({ ...form, description: e.target.value })} placeholder="Описание" rows={3} />
                    <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
                        <input type="number" step="0.1" value={form.weight} onChange={e => setForm({ ...form, weight: e.target.value })} placeholder="Вес (кг)" required />
                        <input value={form.barcode} onChange={e => setForm({ ...form, barcode: e.target.value })} placeholder="Штрихкод" required />
                    </div>
                    <textarea
                        value={form.characteristics}
                        onChange={e => setForm({ ...form, characteristics: e.target.value })}
                        placeholder='Характеристики (JSON)'
                        rows={3}
                        style={{ fontFamily: 'monospace', fontSize: '0.9rem' }}
                    />

                    <div style={{ display: 'flex', gap: '1rem', marginTop: '0.5rem' }}>
                        <button type="submit" className="btn btn-primary" style={{ flex: 1 }}>Сохранить</button>
                    </div>
                </form>
            </Modal>
        </motion.div>
    );
};

export default Products;
