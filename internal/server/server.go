package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/multimoml/tracker/internal/config"
	"github.com/multimoml/tracker/internal/model"
)

func Run(ctx context.Context) {
	// Load environment variables
	cfg := config.LoadConfig()

	// Set up router
	router := gin.Default()

	// Endpoints
	router.GET("/tracker/live", Liveness)
	router.GET("/tracker/ready", Readiness)

	// API v1
	v1 := router.Group("/tracker/v1")
	{
		v1.GET("/all", NewProducts)
	}

	// Start HTTP server
	log.Fatal(router.Run(fmt.Sprintf(":%s", cfg.Port)))
}

func Liveness(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"status": "alive"})
}

func Readiness(c *gin.Context) {
	dispatcher := "http://dispatcher:6001"

	// if using dev environment access local tracker
	if os.Getenv("ACTIVE_ENV") == "dev" {
		dispatcher = "http://localhost:6001"
	}

	_, err := http.Get(dispatcher + "/products/ready")
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"status": "ready"})

}

func NewProducts(c *gin.Context) {
	dispatcher := "http://dispatcher:6001"

	// if using dev environment access local tracker
	if os.Getenv("ACTIVE_ENV") == "dev" {
		dispatcher = "http://localhost:6001"
	}

	res, err := http.Get(dispatcher + "/products/v1/all")
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// decode JSON response into products
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// decode body into products
	var products []model.Product
	err = json.Unmarshal(body, &products)
	if err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//
	var newProducts []model.Product
	for _, product := range products {
		if product.PriceInTime[0].IsNew {
			newProducts = append(newProducts, product)
		}
	}
	if len(newProducts) == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "no new products found"})
		return
	}

	c.IndentedJSON(http.StatusOK, newProducts)

}
