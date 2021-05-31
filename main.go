package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yaron/wishlist-images/src/pages"
	"github.com/yaron/wishlist-images/src/utils"
)

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())
	//r.GET("/list", pages.List)
	authorized := r.Group("/", jWTAuth)
	authorized.POST("/add", pages.Add)
	//authorized.POST("/delete/:id", pages.Delete)
	r.Run()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func jWTAuth(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")
	token := strings.Replace(header, "Bearer ", "", 1)
	userID, err := utils.TestToken(token)
	if err != nil {
		log.Println("Warning: " + err.Error())
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set(gin.AuthUserKey, userID)

}
