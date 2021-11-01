package pages

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yaron/wishlist-images/src/utils"
)

// Parse a page to search for images.
func Parse(c *gin.Context) {
	var json utils.Page
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("Warning: " + err.Error())
		return
	}
	images, err := utils.ParsePage(json.URI)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("Warning: " + err.Error())
		return
	}

	c.JSON(200, gin.H{"status": "Found image", "images": images})
}
