package server_tcp

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func startAPI() {
	router := gin.Default()
	router.GET("/info", getInfo)

	router.Run("localhost:8080")
}

func getInfo(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, Data)
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
