package methods

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"log"
	"ohMyNATS/internal/database/models"
)

func Other() {
	// todo: other action
}

func RestoreCacheFromDB(client *redis.Client, db *gorm.DB) {
	var records []models.Order

	db.Find(&records)

	for _, record := range records {
		jsonData, err := json.Marshal(record)
		if err != nil {
			log.Fatal(err)
		}
		err = client.Set(record.OrderUID, jsonData, 0).Err()
		if err != nil {
			log.Fatal(err)
		}
	}
}
