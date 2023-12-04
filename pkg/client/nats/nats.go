package nats

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
	"ohMyNATS/pkg/client"
)

type Nats struct {
	client.Builder
	Cache     *redis.Client
	DB        *gorm.DB
	ClusterID string
	ClientID  string
	User      string
	Password  string
	Channel   string
}

func (nts *Nats) New() (stan.Subscription, error) {
	connect, err := stan.Connect(
		nts.ClusterID,
		nts.ClientID,
		stan.NatsOptions(nats.UserInfo(nts.User, nts.Password)),
	)
	if err != nil {
		log.Fatal(err)
	}

	subscribe, err := connect.Subscribe(nts.Channel, func(msg *stan.Msg) {
		err := nts.Cache.Set(msg.Subject, string(msg.Data), 0).Err()
		if err != nil {
			return
		}

		//record := models.Cache{ID: msg.Subject, Data: string(msg.Data)}
		//nts.DB.Create(&record)
	})

	return subscribe, err
}
