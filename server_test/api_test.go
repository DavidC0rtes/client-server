package server_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DavidC0rtes/client-server/server"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
)

func setUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

// Tests thats the /info endpoint returns [] because the server is not running.
func TestInfoHandler(t *testing.T) {

	mockResponse := `[]`

	r := setUpRouter()
	r.GET("/info", server.InfoHandler)

	req, _ := http.NewRequest("GET", "/info", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resData, _ := ioutil.ReadAll(w.Body)
	assert.Equal(t, mockResponse, string(resData))
	assert.Equal(t, http.StatusOK, w.Code)
}
