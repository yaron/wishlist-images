package pages

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yaron/wishlist-images/src/utils"
)

// Add creates new images on the local filesystem
func Add(c *gin.Context) {
	var json utils.Image
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("Warning: " + err.Error())
		return
	}
	filename, err := utils.AddItem(json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println("Warning: " + err.Error())
		return
	}
	c.JSON(200, gin.H{"status": "Item added", "localUri": filename})
}
