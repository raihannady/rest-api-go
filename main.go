package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Order struct {
    ID           uint      `gorm:"primary_key" json:"id"`
    CustomerName string    `json:"customerName"`
    OrderedAt    time.Time `json:"orderedAt"`
    Items        []Item    `json:"items" gorm:"foreignkey:OrderID"`
}

type Item struct {
    ID          uint   `gorm:"primary_key" json:"id"`
    ItemCode    string `json:"itemCode"`
    Description string `json:"description"`
    Quantity    int    `json:"quantity"`
    OrderID     uint   `json:"-"`
}

type CreateOrderRequest struct {
    OrderedAt    time.Time `json:"orderedAt"`
    CustomerName string    `json:"customerName"`
    Items        []Item    `json:"items"`
}

type CreateOrderResponse struct {
    ID           uint      `json:"id"`
    OrderedAt    time.Time `json:"orderedAt"`
    CustomerName string    `json:"customerName"`
    Items        []Item    `json:"items"`
}

var db *gorm.DB

func main() {
    var err error
    db, err = gorm.Open("postgres", "user=postgres port=5432 dbname=postgres password=admin sslmode=disable")
    if err != nil {
        log.Fatal("Failed to connect to database: ", err)
    }
    defer db.Close()

    db.AutoMigrate(&Order{}, &Item{})

    r := gin.Default()

    // Routes
    r.POST("/orders", createOrder)
    r.GET("/orders", getOrders)
	r.GET("/orders/:id", getOrderByID)
    r.PUT("/orders/:id", updateOrder)
    r.DELETE("/orders/:id", deleteOrder)


    if err := r.Run(":8080"); err != nil {
        log.Fatal("Failed to start server: ", err)
    }
}

// Create
func createOrder(c *gin.Context) {
    var request CreateOrderRequest
    if err := c.BindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    order := Order{
        CustomerName: request.CustomerName,
        OrderedAt:    request.OrderedAt,
        Items:        request.Items,
    }
    db.Create(&order)

    response := CreateOrderResponse{
        ID:           order.ID,
        CustomerName: order.CustomerName,
        OrderedAt:    order.OrderedAt,
        Items:        order.Items,
    }
    c.JSON(http.StatusCreated, response)
}

// Get
func getOrders(c *gin.Context) {
    var orders []Order
    db.Preload("Items").Find(&orders)

    c.JSON(http.StatusOK, orders)
}

// Get by ID
func getOrderByID(c *gin.Context) {
    id := c.Param("id")

    var order Order
    if db.Preload("Items").First(&order, id).RecordNotFound() {
        c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
        return
    }

    c.JSON(http.StatusOK, order)
}

// Update
func updateOrder(c *gin.Context) {
    id := c.Param("id")

    var order Order
    if db.First(&order, id).RecordNotFound() {
        c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
        return
    }

    var request CreateOrderRequest
    if err := c.BindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    order.CustomerName = request.CustomerName
    order.OrderedAt = request.OrderedAt
    order.Items = request.Items

    db.Save(&order)

    response := CreateOrderResponse{
        ID:           order.ID,
        CustomerName: order.CustomerName,
        OrderedAt:    order.OrderedAt,
        Items:        order.Items,
    }
    c.JSON(http.StatusOK, response)
}

// Delete
func deleteOrder(c *gin.Context) {
    id := c.Param("id")

    var order Order
    if db.First(&order, id).RecordNotFound() {
        c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
        return
    }

    db.Delete(&order)

    c.JSON(http.StatusOK, gin.H{"message": "Success delete"})
}
