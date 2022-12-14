package server

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// startAPI, starts an API ;-) leveraging on gin router capabilities
// the only route is /info which returns a JSON of the Data variable.
// That is the only purpose of the API.
func startAPI() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	router.Use(cors.New(config))
	router.GET("/info", InfoHandler)

	router.Run("localhost:8080")
}

func InfoHandler(c *gin.Context) {
	// Converting the Data var to slice makes it easier to handle in the frontend
	dataArr := make([]Info, len(Data))
	for i, value := range Data {
		fmt.Println(value)
		dataArr[i] = value
	}
	c.IndentedJSON(http.StatusOK, dataArr)
}
