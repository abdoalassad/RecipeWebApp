package Config

import (
	"context"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
)

func InitializingMongoDB() (context.Context, *mongo.Database) {
	ctx := context.Background()
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI("mongodb://admin:password@localhost:27017/test?authSource=admin"),
	)
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection := client.Database("demo")
	return ctx, collection
}

func InitializingRedisCache() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping()
	log.Println(status)
	return redisClient
}

//func Initialize() (context.Context, *mongo.Database, *redis.Client) {
//	ctx := context.Background()
//	client, err := mongo.Connect(
//		ctx,
//		options.Client().ApplyURI("mongodb://admin:password@localhost:27017/test?authSource=admin"),
//	)
//	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
//		log.Fatal(err)
//	}
//	log.Println("Connected to MongoDB")
//	collection := client.Database("demo")
//	//.Collection("recipes")
//
//	redisClient := redis.NewClient(&redis.Options{
//		Addr:     "localhost:6379",
//		Password: "",
//		DB:       0,
//	})
//	status := redisClient.Ping()
//	log.Println(status)
//	return ctx,collection,redisClient
//}
