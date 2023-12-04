package methods

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/goccy/go-json"
	"log"
	"net/http"
	"ohMyNATS/internal/database/models"
)

func Get() {
	// todo: read action
}

func GetCacheByID(client *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		data, err := client.Get(id).Result()
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get data from cache"})
			return
		}

		var order models.Order
		err = json.Unmarshal([]byte(data), &order)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal order data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": order})
	}
}
