package db

import (
	"context"
	"fmt"
	"time"

	"github.com/rem-gestion/rem-common/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongo(cfg config.MongoConfig) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Si cfg.URI incluye credenciales y la DB, ok; si no, armas:
	uri := cfg.URI
	if uri == "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%d",
			cfg.User, cfg.Password, cfg.Host, cfg.Port,
		)
	}
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return client, nil
}
