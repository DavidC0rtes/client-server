package api

import (
	"github.com/DavidC0rtes/client-server/server-tcp"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Start() {
	router := gin.Default()
	router.GET("/info", getInfo)

	router.Run("localhost:8080")
}

func getInfo(c *gin.Context) {
	//info := copyMap(*server_tcp.GetData())
	c.IndentedJSON(http.StatusOK, *server_tcp.GetData())
	/*b, err := json.Marshal(info)
	checkError(err)
	fmt.Printf("%v\n", server_tcp.Data)
	c.Data(200, "application/json", b)*/
}

func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func copyMap(og map[int]server_tcp.Info) map[int]server_tcp.Info {
	targetMap := make(map[int]server_tcp.Info)
	for k, v := range og {
		targetMap[k] = v
	}
	return targetMap
}
