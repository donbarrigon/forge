package db

import (
	"context"
	"fmt"
	"time"

	"github.com/donbarrigon/forge/env"
	"github.com/donbarrigon/forge/errs"
	"github.com/donbarrigon/forge/logs"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client
var Mongo *mongo.Database

func Col(col string) (*mongo.Collection, *errs.Error) {
	if Mongo == nil {
		logs.Error("db is not initialized. Call InitMongoDB() first.")
		return nil, errs.InternalMsg("db is not initialized", nil)
	}
	return Mongo.Collection(col), nil
}

func InitMongoDB() *errs.Error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(env.DB.ConnectionString).SetServerAPIOptions(serverAPI)
	clientOptions.SetMaxPoolSize(env.DB.ClientOptions.MaxPoolSize)
	clientOptions.SetMinPoolSize(env.DB.ClientOptions.MinPoolSize)
	clientOptions.SetRetryWrites(env.DB.ClientOptions.RetryWrites)
	clientOptions.SetTimeout(time.Duration(env.DB.ClientOptions.Timeout) * time.Second)

	var err error

	Client, err = mongo.Connect(clientOptions)
	if err != nil {
		msg := fmt.Sprintf("🔴💥 Fail to connect db %s: %s", env.DB.Name, env.DB.ConnectionString)
		logs.Error(msg)
		return errs.InternalMsg(msg, err)
	}
	Mongo = Client.Database(env.DB.Name)

	ctx := context.TODO()
	if err = Client.Ping(ctx, nil); err != nil {
		msg := fmt.Sprintf("🔴💥 Fail to ping db %s: %s", env.DB.Name, env.DB.ConnectionString)
		logs.Error(msg)
		CloseMongoDB()
		return errs.InternalMsg(msg, err)
	}

	logs.Info("🍃 Successful connection to %s: %s", env.DB.Name, env.DB.ConnectionString)
	return nil
}

func CloseMongoDB() *errs.Error {
	ctx := context.TODO()
	if err := Client.Disconnect(ctx); err != nil {
		logs.Error("Failed to disconnect from MongoDB: %v", err)
		return errs.InternalMsg("Failed to disconnect from MongoDB", err)
	}
	return nil
}
