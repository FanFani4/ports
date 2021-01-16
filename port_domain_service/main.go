package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/FanFani4/ports/port_domain_service/server"
	"github.com/FanFani4/ports/ports"
	"google.golang.org/grpc"
)

func getMongoCollection(log *logrus.Logger) *mongo.Collection {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCtx()

	fmt.Println("mongodb://" + os.Getenv("MONGO"))
	mgo, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI("mongodb://"+os.Getenv("MONGO")),
	)

	if err != nil {
		log.Fatal("Mongo Client error: " + err.Error())
	}

	a, _ := mgo.ListDatabaseNames(ctx, nil)
	fmt.Println(a)

	err = mgo.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Ping Mongo Server error: "+err.Error(), http.StatusInternalServerError)
	}

	return mgo.Database(os.Getenv("DATABASE")).Collection(os.Getenv("COLLECTION"))
}

func main() {
	log := logrus.New()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", ":"+os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}

	db := getMongoCollection(log)

	s := grpc.NewServer()
	ports.RegisterPortDomainServiceServer(s, server.NewServer(db))

	if err := s.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
