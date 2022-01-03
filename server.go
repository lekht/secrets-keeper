package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func indexView(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func saveMessageView(c *gin.Context) {
	message := c.PostForm("message")
	key := keyBuilder.Get()
	// TODO: handle error
	keeper.Set(key, message)
	c.HTML(http.StatusOK, "key.html", gin.H{"key": fmt.Sprintf("http://%s/%s", c.Request.Host, key)})
}

func readMessageHandler(c *gin.Context) {
	key := c.Param("key")
	msg, err := keeper.Get(key)
	if err != nil {
		if err.Error() == NotFoundError {
			c.HTML(http.StatusNotFound, "404.html", gin.H{})
			return
		}
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		return
	}
	c.HTML(http.StatusOK, "message.html", gin.H{"message": msg})
	err = keeper.Clean(key)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		return
	}
}

func getRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLFiles(
		"templates/index.html",
		"templates/key.html",
		"templates/message.html",
		"templates/404.html",
		"templates/500.html",
	)
	router.GET("/", indexView)
	router.POST("/", saveMessageView)
	router.GET("/:key", readMessageHandler)
	return router

}

func main() {
	router := getRouter()
	router.Run("localhost:8080")
}
