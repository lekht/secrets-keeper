package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func writeInternalError(c *gin.Context) {
	c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
}

func indexView(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func saveMessageView(c *gin.Context, keybuilder KeyBuilder, keeper Keeper) {
	message := c.PostForm("message")
	key, err := keybuilder.Get()
	if err != nil {
		writeInternalError(c)
		return
	}
	err = keeper.Set(key, message)
	if err != nil {
		writeInternalError(c)
		return
	}
	c.HTML(http.StatusOK, "key.html", gin.H{"key": fmt.Sprintf("http://%s/%s", c.Request.Host, key)})
}

func readMessageHandler(c *gin.Context, keybuilder KeyBuilder, keeper Keeper) {
	key := c.Param("key")
	msg, err := keeper.Get(key)
	if err != nil {
		if err.Error() == NotFoundError {
			c.HTML(http.StatusNotFound, "404.html", gin.H{})
			return
		}
		writeInternalError(c)
		return
	}
	err = keeper.Clean(key)
	if err != nil {
		writeInternalError(c)
		return
	}
	c.HTML(http.StatusOK, "message.html", gin.H{"message": msg})
}

func buildHandler(fn func(c *gin.Context, keybuilder KeyBuilder, keeper Keeper), keybuilder KeyBuilder, keeper Keeper) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(c, keybuilder, keeper)
	}
}

func getRouter(keybuilder KeyBuilder, keeper Keeper) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLFiles(
		"templates/index.html",
		"templates/key.html",
		"templates/message.html",
		"templates/404.html",
		"templates/500.html",
	)
	router.GET("/", indexView)
	router.POST("/", buildHandler(saveMessageView, keybuilder, keeper))
	router.GET("/:key", buildHandler(readMessageHandler, keybuilder, keeper))
	return router

}

func main() {
	keyBuilder := UUIDKeyBuilder{}
	keeper := getKeeper()
	router := getRouter(keyBuilder, keeper)
	router.Run("localhost:8080")
}
