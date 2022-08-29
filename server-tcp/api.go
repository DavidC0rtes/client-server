package server_tcp

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func startAPI() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	router.Use(cors.New(config))
	router.GET("/info", getInfo)

	router.Run("localhost:8080")
}

func getInfo(c *gin.Context) {
	// Converting the Data var to slice makes it easier to handle in the frontend
	dataArr := make([]Info, len(Data))
	for i, value := range Data {
		fmt.Println(value)
		dataArr[i] = value
	}
	c.IndentedJSON(http.StatusOK, dataArr)
}

func checkError(err error, id, channel int) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
