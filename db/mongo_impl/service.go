package mongoimpl

import (
	"context"
	"time"

	"github.com/vinhphuctadang/go-interface/db"
	"github.com/vinhphuctadang/go-interface/db/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbServiceMongoBackend struct {
	db.DBService
	client *mongo.Client
}

func NewDbServiceMongoBackend(connectionUri string) (db.DBService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionUri))
	if err != nil {
		return nil, err
	}

	return &dbServiceMongoBackend{
		client: client,
	}, nil
}

func (d *dbServiceMongoBackend) CreateAccount(username, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	d.client.
		Database("go_interface_example").
		Collection("accounts").
		InsertOne(
			ctx,
			bson.D{
				{"username", username},
				{"password", password},
			},
		)
	return nil
}

func (d *dbServiceMongoBackend) ListAccount(pageIndex, pageSize int) (accounts []model.Account, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := d.client.
		Database("go_interface_example").
		Collection("accounts").
		Aggregate(
			ctx,
			mongo.Pipeline{
				bson.D{{"$skip", pageIndex}},
				bson.D{{"$limit", pageSize}},
			},
		)
	if err != nil {
		return nil, err
	}

	var result []bson.M
	if err := cursor.All(ctx, result); err != nil {
		return nil, err
	}

	for _, r := range result {
		accounts = append(accounts, model.Account{
			Username: r["username"].(string), // not good when: db upgrade/migration
			Password: r["password"].(string),
		})
	}

	return accounts, nil
}

func (d *dbServiceMongoBackend) Disconnect(ctx context.Context) error {
	return d.client.Disconnect(ctx)
}
