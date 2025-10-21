package handlers

import (
	"backend/database"
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, description, price, image_url, stock FROM products ORDER BY id")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.Stock); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		products = append(products, p)
	}

	c.JSON(http.StatusOK, products)
}

func GetCart(c *gin.Context) {
	rows, err := database.DB.Query(`
        SELECT p.id, p.name, p.description, p.price, p.image_url, c.quantity 
        FROM cart c 
        JOIN products p ON c.product_id = p.id
        ORDER BY p.id
    `)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	cartItems := make([]map[string]interface{}, 0)
	for rows.Next() {
		var product models.Product
		var quantity int
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.ImageURL, &quantity); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cartItems = append(cartItems, map[string]interface{}{
			"product":  product,
			"quantity": quantity,
		})
	}

	c.JSON(http.StatusOK, cartItems)
}

func AddToCart(c *gin.Context) {
	var item models.CartItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1 AND stock > 0)", item.ProductID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Товар не существует или нет в наличии"})
		return
	}

	_, err = database.DB.Exec(`
        INSERT INTO cart (product_id, quantity) 
        VALUES ($1, 1)
        ON CONFLICT (product_id) 
        DO UPDATE SET quantity = cart.quantity + 1
    `, item.ProductID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Товар добавлен в корзину"})
}

func RemoveFromCart(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	_, err = database.DB.Exec("DELETE FROM cart WHERE product_id = $1", productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Товар удален из корзины"})
}

func Checkout(c *gin.Context) {
	rows, err := database.DB.Query(`
        SELECT c.product_id, c.quantity, p.price 
        FROM cart c 
        JOIN products p ON c.product_id = p.id
    `)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var items []models.CartItem
	var totalPrice float64

	for rows.Next() {
		var item models.CartItem
		var price float64
		if err := rows.Scan(&item.ProductID, &item.Quantity, &price); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		items = append(items, item)
		totalPrice += price * float64(item.Quantity)
	}

	if len(items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Корзина пуста"})
		return
	}

	itemsJSON, _ := json.Marshal(items)
	_, err = database.DB.Exec("INSERT INTO orders (items, total_price) VALUES ($1, $2)", string(itemsJSON), totalPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = database.DB.Exec("DELETE FROM cart")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Заказ оформлен", "total": totalPrice})
}
