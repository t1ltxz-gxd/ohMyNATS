package nats

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/nats-io/stan.go"
	"log"
	"ohMyNATS/internal/database/models"
	"ohMyNATS/pkg/client"
)

type Nats struct {
	client.Builder
	Client    *redis.Client
	DB        *gorm.DB
	ClusterID string
	ClientID  string
}

func (nats *Nats) New() (stan.Subscription, error) {
	sc, err := stan.Connect(nats.ClusterID, nats.ClientID)
	if err != nil {
		log.Fatal(err)
	}

	sub, err := sc.Subscribe("test-channel", func(msg *stan.Msg) {
		err := nats.Client.Set(msg.Subject, string(msg.Data), 0).Err()
		if err != nil {
			return
		}

		record := models.Cache{ID: msg.Subject, Data: string(msg.Data)}
		nats.DB.Create(&record)
	})

	return sub, err
}
