class ShopApp {
    constructor() {
        this.products = [];
        this.cart = [];
        this.apiBase = window.appConfig ?
            `${window.appConfig.backendUrl}/api` :
            '/api';
        this.init();
    }

    async init() {
        await this.loadProducts();
        await this.loadCart();
        this.render();
    }

    async fetchAPI(endpoint, options = {}) {
        try {
            const response = await fetch(`${this.apiBase}${endpoint}`, {
                headers: {
                    'Content-Type': 'application/json',
                    ...options.headers,
                },
                ...options,
            });
            
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            
            return await response.json();
        } catch (error) {
            this.showError(`Ошибка: ${error.message}`);
            throw error;
        }
    }

    async loadProducts() {
        try {
            this.products = await this.fetchAPI('/products');
        } catch (error) {
            console.error('Failed to load products:', error);
        }
    }

    async loadCart() {
        try {
            this.cart = await this.fetchAPI('/cart');
        } catch (error) {
            console.error('Failed to load cart:', error);
        }
    }

    async addToCart(productId) {
        try {
            await this.fetchAPI('/cart', {
                method: 'POST',
                body: JSON.stringify({
                    product_id: productId,
                    quantity: 1
                })
            });
            await this.loadCart();
            this.render();
            this.showSuccess('Товар добавлен в корзину!');
        } catch (error) {
            console.error('Failed to add to cart:', error);
        }
    }

    async removeFromCart(productId) {
        try {
            await this.fetchAPI(`/cart/${productId}`, {
                method: 'DELETE'
            });
            await this.loadCart();
            this.render();
            this.showSuccess('Товар удален из корзины!');
        } catch (error) {
            console.error('Failed to remove from cart:', error);
        }
    }

    async checkout() {
        try {
            const result = await this.fetchAPI('/checkout', {
                method: 'POST'
            });
            await this.loadCart();
            this.render();
            this.showSuccess(`Заказ оформлен! Сумма: $${result.total}`);
        } catch (error) {
            console.error('Failed to checkout:', error);
        }
    }

    render() {
        this.renderProducts();
        this.renderCart();
    }

    renderProducts() {
        const container = document.getElementById('products-container');
        if (!container) return;

        container.innerHTML = this.products.map(product => `
            <div class="product-card">
                <div class="product-image">
                    ${product.image_url ? 
                        `<img src="${product.image_url}" alt="${product.name}" style="max-width: 100%; max-height: 100%;">` : 
                        '📱'
                    }
                </div>
                <div class="product-name">${product.name}</div>
                <div class="product-description">${product.description}</div>
                <div class="product-price">$${product.price}</div>
                <div class="product-stock">В наличии: ${product.stock} шт.</div>
                <button class="btn btn-primary" onclick="app.addToCart(${product.id})">
                    Добавить в корзину
                </button>
            </div>
        `).join('');
    }

    renderCart() {
        const container = document.getElementById('cart-container');
        const emptyCart = document.getElementById('empty-cart');
        const cartItems = document.getElementById('cart-items');
        const totalElement = document.getElementById('cart-total');
        const checkoutBtn = document.getElementById('checkout-btn');

        if (!container) return;

        console.log(this.cart)

        if (this.cart.length === 0) {
            emptyCart.classList.remove('hidden');
            cartItems.classList.add('hidden');
            totalElement.classList.add('hidden');
            checkoutBtn.classList.add('hidden');
            return;
        }

        emptyCart.classList.add('hidden');
        cartItems.classList.remove('hidden');
        totalElement.classList.remove('hidden');
        checkoutBtn.classList.remove('hidden');

        cartItems.innerHTML = this.cart.map(item => `
            <div class="cart-item">
                <div>
                    <strong>${item.product.name}</strong><br>
                    <small>$${item.product.price} x ${item.quantity}</small>
                </div>
                <div>
                    <strong>$${(item.product.price * item.quantity).toFixed(2)}</strong>
                    <button class="btn btn-danger" onclick="app.removeFromCart(${item.product.id})" style="margin-left: 10px;">
                        Удалить
                    </button>
                </div>
            </div>
        `).join('');

        const total = this.cart.reduce((sum, item) => sum + (item.product.price * item.quantity), 0);
        totalElement.textContent = `Общая сумма: $${total.toFixed(2)}`;
    }

    showError(message) {
        this.showMessage(message, 'error');
    }

    showSuccess(message) {
        this.showMessage(message, 'success');
    }

    showMessage(message, type) {
        const messageDiv = document.createElement('div');
        messageDiv.className = type;
        messageDiv.textContent = message;
        
        const container = document.querySelector('.container');
        container.insertBefore(messageDiv, container.firstChild);
        
        setTimeout(() => {
            messageDiv.remove();
        }, 5000);
    }
}

const app = new ShopApp();
