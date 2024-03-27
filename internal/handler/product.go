package handler

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/internal/model"
	"main.go/internal/repository"
)

type ProductHandler struct {
	repo *repository.ProductRepository
}

func NewProductHandler(db *sql.DB) *ProductHandler {
	return &ProductHandler{
		repo: repository.NewProductRepository(db),
	}
}

func (h *ProductHandler) RegisterProductRoutes(router *gin.Engine) {
	router.GET("/products", h.GetProducts)
	router.POST("/products", h.CreateProduct)
	router.PATCH("/products/:id", h.UpdateProduct)
	router.DELETE("/products/:id", h.DeleteProduct)
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	var products []model.Products
	rows, err := h.repo.DB.Query(`SELECT * FROM products ORDER BY id`)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Select Query Errror"})
		return
	}
	defer rows.Close()

	// log.Print(rows)
	for rows.Next() {
		var p model.Products
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Quantity, &p.Description); err != nil {
			log.Print(`Failed to get all products`)
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Select Scan Errror " + err.Error()})
			return
		}

		products = append(products, p)
	}

	c.IndentedJSON(http.StatusOK, products)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req model.Products

	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad POST request"})
		return
	}

	_, err := h.repo.DB.Exec(`INSERT INTO products(ID, name, price, quantity, description) VALUES ($1, $2, $3, $4, $5)`,
		&req.ID, &req.Name, &req.Price, &req.Quantity, &req.Description)

	if err != nil {
		// It's useful to log the error for debugging purposes
		log.Printf("Failed to insert new product: %v", err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert new product"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"Product added successfully": &req})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var req model.Products
	id := c.Param("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Bad PATCH request"})
		return
	}

	currentProduct, err := h.repo.GetProductById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Product not found!"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving product"})
		}
		return
	}

	if req.Name != "" {
		currentProduct.Name = req.Name
	}
	if req.Price != 0 {
		currentProduct.Price = req.Price
	}
	if req.Quantity != 0 {
		currentProduct.Quantity = req.Quantity
	}
	if req.Description != "" {
		currentProduct.Description = req.Description
	}

	_, err = h.repo.DB.Exec(`UPDATE products SET name = $1, price = $2, quantity = $3, description = $4 WHERE id = $5`, &currentProduct.Name, &currentProduct.Price, &currentProduct.Quantity, &currentProduct.Description, &id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed to update product"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"Updated product": currentProduct})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	_, err := h.repo.GetProductById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Product not found!"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving product"})
		}
		return
	}

	_, err = h.repo.DB.Exec(`DELETE FROM products WHERE id = $1`, &id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed to delete product"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"Deleted product": id})
}
