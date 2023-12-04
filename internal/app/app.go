package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/stan.go"
	"log"
	"ohMyNATS/internal/database/methods"
	"ohMyNATS/internal/database/models"
	"ohMyNATS/internal/misc"
	"ohMyNATS/pkg/client/database/postgres"
	"ohMyNATS/pkg/client/database/redis"
	"ohMyNATS/pkg/client/nats"
	"ohMyNATS/pkg/logger/zerolog"
	"os"
	"strconv"
)

func Serve(Env string, Port int) {
	logger := zerolog.LoggerBuilder(Env)
	logger.Debug().Msg("Debug mode is enabled!")
	logger.Info().Str("environment", Env).Msg("Starting app...")

	// Loading .env file.
	logger.Info().Str("env-file", ".env").Msg("Loading environment variables...")
	misc.LoadDotEnv(".env")

	// Build cache storage
	logger.Info().Msg("Initializing Redis client...")
	redisHost := os.Getenv("REDIS_HOST")
	redisPort, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Fatalf("Convertion error: %s", err)
	}
	redisClient := redis.Redis{
		Addr:     fmt.Sprintf("%s:%d", redisHost, redisPort),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	}
	cache, err := redisClient.New()
	if err != nil {
		log.Fatalf("Cache creation error: %s", err)
	}

	// Build database storage.
	logger.Info().Msg("Initializing Postgres client...")
	pgPort, err := strconv.Atoi(os.Getenv("PG_PORT"))
	if err != nil {
		log.Fatalf("Formating error: %s", err)
	}
	pgClient := postgres.Postgres{
		Host:     "localhost",
		Port:     pgPort,
		User:     os.Getenv("PG_USER"),
		Dbname:   os.Getenv("PG_DB"),
		Password: os.Getenv("PG_PASS"),
	}
	db, err := pgClient.New()
	if err != nil {
		log.Fatalf("Database creation error: %s", err)
	}
	logger.Info().Msg("Migrating models...")
	go db.AutoMigrate(&models.Order{})

	//orderJSON := `{
	//    "order_uid": "b563feb7b2b84b6test",
	//    "track_number": "WBILMTESTTRACK",
	//    "entry": "WBIL",
	//    "delivery": {
	//        "name": "Test Testov",
	//        "phone": "+9720000000",
	//        "zip": "2639809",
	//        "city": "Kiryat Mozkin",
	//        "address": "Ploshad Mira 15",
	//        "region": "Kraiot",
	//        "email": "test@gmail.com"
	//    },
	//    "payment": {
	//        "transaction": "b563feb7b2b84b6test",
	//        "request_id": "",
	//        "currency": "USD",
	//        "provider": "wbpay",
	//        "amount": 1817,
	//        "payment_dt": 1637907727,
	//        "bank": "alpha",
	//        "delivery_cost": 1500,
	//        "goods_total": 317,
	//        "custom_fee": 0
	//    },
	//    "items": [
	//        {
	//            "chrt_id": 9934930,
	//            "track_number": "WBILMTESTTRACK",
	//            "price": 453,
	//            "rid": "ab4219087a764ae0btest",
	//            "name": "Mascaras",
	//            "sale": 30,
	//            "size": "0",
	//            "total_price": 317,
	//            "nm_id": 2389212,
	//            "brand": "Vivienne Sabo",
	//            "status": 202
	//        }
	//    ],
	//    "locale": "en",
	//    "internal_signature": "",
	//    "customer_id": "test",
	//    "delivery_service": "meest",
	//    "shardkey": "9",
	//    "sm_id": 99,
	//    "date_created": "2021-11-26T06:22:19Z",
	//    "oof_shard": "1"
	//}`
	//
	//var order models.Order
	//err = json.Unmarshal([]byte(orderJSON), &order)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = db.Create(&order).Error
	//if err != nil {
	//	log.Fatal(err)
	//}

	logger.Info().Msg("Restoring cache from database...")
	go methods.RestoreCacheFromDB(cache, db)

	// Build Nats subscription.
	logger.Info().Msg("Initializing NATS client...")
	natsClient := nats.Nats{
		Cache:     cache,
		DB:        db,
		ClusterID: os.Getenv("CLUSTER_ID"),
		ClientID:  os.Getenv("CLIENT_ID"),
		User:      os.Getenv("NATS_USER"),
		Password:  os.Getenv("NATS_PASS"),
		Channel:   os.Getenv("NATS_CHANNEL"),
	}
	sub, err := natsClient.New()
	if err != nil {
		log.Fatalf("Subsribtion error: %s", err)
	}
	defer func(sub stan.Subscription) {
		err := sub.Unsubscribe()
		if err != nil {
			log.Fatalf("Unsubscribtion error: %s", err)
		}
	}(sub)

	router := gin.Default()

	// Implement handlers
	logger.Info().Msg("Initializing router paths...")
	router.GET("/orders/:id", methods.GetCacheByID(cache))

	// Run HTTP server
	logger.Info().Msg("Starting HTTP client...")
	if err := router.Run(fmt.Sprintf(":%d", Port)); err != nil {
		log.Fatal(err)
	}
}
