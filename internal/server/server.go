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
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/multimoml/tracker/docs"

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
	router.GET("/tracker/openapi", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/tracker/openapi/index.html")
	})
	router.GET("/tracker/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/tracker/openapi/index.html")
	})
	router.GET("/tracker/openapi/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1
	v1 := router.Group("/tracker/v1")
	{
		v1.GET("/new", NewProducts)
	}

	// Start HTTP server
	log.Fatal(router.Run(fmt.Sprintf(":%s", cfg.Port)))

}

// Liveness is a simple endpoint to check if the server is alive
// @Summary Get liveness status of the microservice
// @Description Get liveness status of the microservice
// @Tags Kubernetes
// @Success 200 {string} string
// @Router /live [get]
func Liveness(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"status": "alive"})
}

// Readiness is a simple endpoint to check if the server is ready
// @Summary Get readiness status of the microservice
// @Description Get readiness status of the microservice
// @Tags Kubernetes
// @Success 200 {string} string
// @Failure 503 {string} string
// @Router /ready [get]
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

// NewProducts returns a list of all new products
// @Summary Get a list of all new products
// @Description Get a list of all new products
// @Tags NewProducts
// @Produce json
// @Success 200 {array} object
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /v1/new [get]
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
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "No new products found."})
		return
	}

	c.IndentedJSON(http.StatusOK, newProducts)

}
