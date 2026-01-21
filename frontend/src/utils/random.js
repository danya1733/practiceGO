export const randomData = {
  warehouses: [
    { address: "ул. Пушкина, д. 10, Москва" },
    { address: "пр. Невский, д. 25, Санкт-Петербург" },
    { address: "ул. Ленина, д. 5, Новосибирск" },
    { address: "ул. Гагарина, д. 15, Екатеринбург" },
    { address: "ул. Советская, д. 42, Казань" }
  ],
  products: [
    { name: "Умная колонка Алиса", description: "Голосовой помощник с отличным звуком.", weight: 0.5, barcode: "8934758934", characteristics: { color: "Черный", generation: "5-е" } },
    { name: "Смартфон ПроПиксель X1", description: "Флагманская камера и AI функции.", weight: 0.2, barcode: "1230981230", characteristics: { storage: "256ГБ", color: "Оникс" } },
    { name: "Монитор УльтраВью 4K", description: "32-дюймовый безрамочный дисплей.", weight: 5.4, barcode: "5675675676", characteristics: { resolution: "3840x2160", refreshRate: "144Гц" } },
    { name: "Клавиатура КиберДек", description: "Механическая клавиатура с подсветкой.", weight: 1.1, barcode: "9988776655", characteristics: { switches: "Синие", layout: "ANSI" } },
    { name: "Мышь Квантум", description: "Беспроводная эргономичная мышь.", weight: 0.15, barcode: "3344556677", characteristics: { dpi: "25000", buttons: "7" } }
  ],

  generateWarehouse: () => {
    const item = randomData.warehouses[Math.floor(Math.random() * randomData.warehouses.length)];
    // Add random suffix to make it unique
    return {
      address: `${item.address} #${Math.floor(Math.random() * 1000)}`
    };
  },

  generateProduct: () => {
    const item = randomData.products[Math.floor(Math.random() * randomData.products.length)];
    return {
      name: `${item.name} ${Math.floor(Math.random() * 100) > 50 ? 'Plus' : ''}`,
      description: item.description,
      weight: item.weight,
      barcode: `${Math.floor(Math.random() * 10000000000)}`,
      characteristics: JSON.stringify(item.characteristics, null, 2)
    };
  },

  generateInventory: () => {
    return {
      quantity: Math.floor(Math.random() * 100) + 1,
      price: (Math.random() * 1000).toFixed(2),
      discount: Math.floor(Math.random() * 20)
    };
  }
};
