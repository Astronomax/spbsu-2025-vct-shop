package database

import (
	"backend/models"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "shopdb"
)

func InitDB() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	if os.Getenv("DB_HOST") != "" {
		psqlInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"))
	}

	var err error
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	log.Println("Successfully connected to PostgreSQL!")

	if err := createTables(); err != nil {
		return err
	}

	return seedProducts()
}

func createTables() error {
	productsTable := `
        CREATE TABLE IF NOT EXISTS products (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            description TEXT,
            price DECIMAL(10,2) NOT NULL,
            image_url TEXT,
            stock INTEGER DEFAULT 0,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `

	cartTable := `
        CREATE TABLE IF NOT EXISTS cart (
            product_id INTEGER PRIMARY KEY,
            quantity INTEGER NOT NULL DEFAULT 0,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
        )
    `

	ordersTable := `
        CREATE TABLE IF NOT EXISTS orders (
            id SERIAL PRIMARY KEY,
            items TEXT NOT NULL,
            total_price DECIMAL(10,2) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `

	tables := []string{productsTable, cartTable, ordersTable}

	for _, table := range tables {
		_, err := DB.Exec(table)
		if err != nil {
			return fmt.Errorf("error creating table: %v", err)
		}
	}

	log.Println("Database tables created successfully!")
	return nil
}

func seedProducts() error {
	products := []models.Product{
		{Name: "iPhone 15", Description: "Новый смартфон от Apple", Price: 999.99, ImageURL: "/static/images/iphone.jpg", Stock: 10},
		{Name: "MacBook Pro", Description: "Мощный ноутбук для работы", Price: 1999.99, ImageURL: "/static/images/macbook.jpg", Stock: 5},
		{Name: "AirPods Pro", Description: "Беспроводные наушники", Price: 249.99, ImageURL: "/static/images/airpods.jpg", Stock: 20},
		{Name: "iPad Air", Description: "Универсальный планшет", Price: 599.99, ImageURL: "/static/images/ipad.jpg", Stock: 8},
	}

	for _, product := range products {
		var exists bool
		err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE name = $1)", product.Name).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			_, err = DB.Exec(
				"INSERT INTO products (name, description, price, image_url, stock) VALUES ($1, $2, $3, $4, $5)",
				product.Name, product.Description, product.Price, product.ImageURL, product.Stock,
			)
			if err != nil {
				return err
			}
		}
	}

	log.Println("Test products seeded successfully!")
	return nil
}
