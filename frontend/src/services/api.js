const API_URL = import.meta.env.VITE_API_URL || '/api';

export const api = {
    // Health
    healthCheck: async () => {
        try {
            const res = await fetch(`${API_URL}/health`);
            return res.ok;
        } catch (e) {
            return false;
        }
    },

    // Warehouses
    getWarehouses: async () => {
        const res = await fetch(`${API_URL}/warehouses`);
        if (!res.ok) throw new Error('Failed to fetch warehouses');
        return res.json();
    },
    createWarehouse: async (data) => {
        const res = await fetch(`${API_URL}/warehouses`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        });
        if (!res.ok) throw new Error('Failed to create warehouse');
        return res.json();
    },

    // Products
    getProducts: async () => {
        const res = await fetch(`${API_URL}/products`);
        if (!res.ok) throw new Error('Failed to fetch products');
        return res.json();
    },
    createProduct: async (data) => {
        let payload = { ...data };
        try {
            if (typeof payload.characteristics === 'string') {
                payload.characteristics = JSON.parse(payload.characteristics);
            }
        } catch (e) {
            payload.characteristics = {};
        }

        const res = await fetch(`${API_URL}/products`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
        });
        if (!res.ok) throw new Error('Failed to create product');
        return res.json();
    },
    updateProduct: async (id, data) => {
        let payload = { ...data };
        try {
            if (typeof payload.characteristics === 'string') {
                payload.characteristics = JSON.parse(payload.characteristics);
            }
        } catch (e) {
            // keep as is or empty
        }

        const res = await fetch(`${API_URL}/products/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
        });
        if (!res.ok) throw new Error('Failed to update product');
        return res.json();
    },

    // Inventory
    getWarehouseProducts: async (id) => {
        const res = await fetch(`${API_URL}/warehouses/${id}/products`);
        if (!res.ok) throw new Error('Failed to fetch warehouse products');
        return res.json();
    },
    getWarehouseProduct: async (warehouseId, productId) => {
        const res = await fetch(`${API_URL}/warehouses/${warehouseId}/products/${productId}`);
        if (!res.ok) throw new Error('Failed to fetch warehouse product');
        return res.json();
    },
    createInventory: async (data) => {
        const res = await fetch(`${API_URL}/inventory`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        });
        if (!res.ok) throw new Error('Failed to add inventory');
        return res.json();
    },
    updateQuantity: async (warehouseId, productId, quantity) => {
        const res = await fetch(`${API_URL}/inventory/quantity`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ warehouse_id: warehouseId, product_id: productId, quantity: parseInt(quantity) }),
        });
        if (!res.ok) throw new Error('Failed to update quantity');
        return res.json();
    },
    updateDiscount: async (warehouseId, productId, discount) => {
        const res = await fetch(`${API_URL}/inventory/discount`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ warehouse_id: warehouseId, product_id: productId, discount: parseFloat(discount) }),
        });
        if (!res.ok) throw new Error('Failed to update discount');
        return res.json();
    },

    // Purchase & Calculation
    calculate: async (warehouseId, products) => {
        const res = await fetch(`${API_URL}/warehouses/calculate`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ warehouse_id: warehouseId, products }),
        });
        if (!res.ok) throw new Error('Failed to calculate price');
        return res.json();
    },
    purchase: async (warehouseId, products) => {
        const res = await fetch(`${API_URL}/warehouses/purchase`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ warehouse_id: warehouseId, products }),
        });
        if (!res.ok) throw new Error('Failed to purchase');
        return res.json();
    },

    // Analytics
    getWarehouseAnalytics: async (id) => {
        const res = await fetch(`${API_URL}/analytics/warehouses/${id}`);
        // Return null if 204 or 404
        if (!res.ok) return null;
        return res.json();
    },
    getTopWarehouses: async () => {
        const res = await fetch(`${API_URL}/analytics/warehouses/top`);
        if (!res.ok) return [];
        return res.json();
    }
};
