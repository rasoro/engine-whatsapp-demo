package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/weni/whatsapp-router/config"
	"github.com/weni/whatsapp-router/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewDB() *mongo.Database {
	dbConf := config.GetConfig().DB
	options := options.Client().ApplyURI(dbConf.URI)
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()
	connection, err := mongo.Connect(ctx, options)
	if err != nil {
		logger.Error("mongodb FAIL")
		panic(err.Error())
	}

	if err := connection.Ping(context.TODO(), readpref.Primary()); err != nil {
		logger.Error("mongodb FAIL")
		panic(err.Error())
	} else {
		logger.Info("mongodb OK")
	}

	db := connection.Database(dbConf.Name)
	return db
}

func NewTestDB() *mongo.Database {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%v/?appName=whatsapp-router", "admin", "admin", "localhost", 27017)
	options := options.Client().ApplyURI(uri)
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()
	connection, err := mongo.Connect(ctx, options)
	if err != nil {
		logger.Error(fmt.Sprintf("mongodb FAIL: %s", err.Error()))
		panic(err.Error())
	}

	if err := connection.Ping(context.TODO(), readpref.Primary()); err != nil {
		logger.Error(fmt.Sprintf("mongodb FAIL: %s", err.Error()))
		panic(err.Error())
	} else {
		logger.Info("mongodb OK")
	}

	db := connection.Database("whatsapp-router-test")
	return db
}

func CleanupDB(db *mongo.Database) {
	err := db.Drop(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDB(db *mongo.Database) {
	if err := db.Client().Disconnect(context.TODO()); err != nil {
		logger.Error(fmt.Sprintf("Error on close MongoDB: %v", err))
		panic(err.Error())
	}
}
