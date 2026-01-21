import React, { useState, useEffect } from 'react';
import Modal from '../components/ui/Modal';
import { api } from '../services/api';
import { ShoppingCart } from 'lucide-react';

const PurchaseModal = ({ isOpen, onClose, warehouseId, inventory, onSuccess }) => {
    const [cart, setCart] = useState({}); // { productId: quantity }
    const [calculation, setCalculation] = useState(null);

    useEffect(() => {
        if (isOpen) {
            setCart({});
            setCalculation(null);
        }
    }, [isOpen]);

    const handleQuantityChange = (productId, qty) => {
        setCart(prev => ({
            ...prev,
            [productId]: qty > 0 ? qty : 0
        }));
    };

    const handleCalculate = async () => {
        const products = Object.entries(cart)
            .filter(([_, qty]) => qty > 0)
            .map(([productId, qty]) => ({ product_id: productId, quantity: parseInt(qty) }));

        if (products.length === 0) return;

        try {
            const data = await api.calculate(warehouseId, products);
            setCalculation(data);
        } catch (e) {
            alert('Ошибка расчета: ' + e.message);
        }
    };

    const handlePurchase = async () => {
        const products = Object.entries(cart)
            .filter(([_, qty]) => qty > 0)
            .map(([productId, qty]) => ({ product_id: productId, quantity: parseInt(qty) }));

        try {
            await api.purchase(warehouseId, products);
            alert('Покупка успешна!');
            onSuccess();
            onClose();
        } catch (e) {
            alert('Ошибка покупки: ' + e.message);
        }
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title="Оформление покупки">
            <div style={{ maxHeight: '60vh', overflowY: 'auto' }}>
                <div style={{ marginBottom: '1rem' }}>
                    {inventory.map(item => (
                        <div key={item.id} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '0.5rem', padding: '0.5rem', background: 'rgba(255,255,255,0.05)', borderRadius: '8px' }}>
                            <div>
                                <div style={{ fontWeight: 'bold' }}>{item.product?.name}</div>
                                <div style={{ fontSize: '0.8rem', color: 'var(--text-secondary)' }}>Доступно: {item.quantity} | Цена: {item.price} ₽ | Скидка: {item.discount}%</div>
                            </div>
                            <input
                                type="number"
                                min="0"
                                max={item.quantity}
                                value={cart[item.product_id] || ''}
                                onChange={(e) => handleQuantityChange(item.product_id, parseInt(e.target.value) || 0)}
                                placeholder="0"
                                style={{ width: '60px', textAlign: 'center' }}
                            />
                        </div>
                    ))}
                </div>

                <button onClick={handleCalculate} className="btn" style={{ width: '100%', marginBottom: '1rem', background: '#3b82f6' }}>
                    Рассчитать стоимость
                </button>

                {calculation && (
                    <div style={{ padding: '1rem', background: 'rgba(50, 205, 50, 0.1)', borderRadius: '8px', border: '1px solid rgba(50, 205, 50, 0.3)' }}>
                        <h3 style={{ marginBottom: '0.5rem' }}>Итого: {calculation.total_sum.toLocaleString()} ₽</h3>
                        <ul style={{ fontSize: '0.9rem', paddingLeft: '1.2rem' }}>
                            {calculation.items.map((item, idx) => (
                                <li key={idx}>
                                    {item.name}: {item.quantity} шт. x {item.price_with_discount} ₽ = {item.total_price} ₽
                                </li>
                            ))}
                        </ul>
                        <button onClick={handlePurchase} className="btn btn-primary" style={{ width: '100%', marginTop: '1rem' }}>
                            <ShoppingCart size={18} /> Купить
                        </button>
                    </div>
                )}
            </div>
        </Modal>
    );
};

export default PurchaseModal;
