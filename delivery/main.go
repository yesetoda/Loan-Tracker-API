package main

import (
	"context"
	"fmt"

	"github.com/yesetoda/b/Loan-Tracker-API/config"
	"github.com/yesetoda/b/Loan-Tracker-API/controller"

	"log"

	"github.com/yesetoda/b/Loan-Tracker-API/infrastructures"
	"github.com/yesetoda/b/Loan-Tracker-API/repository"
	"github.com/yesetoda/b/Loan-Tracker-API/router"
	"github.com/yesetoda/b/Loan-Tracker-API/usecase"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(config.Database.Uri)
	clientOptions := options.Client().ApplyURI(config.Database.Uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	userCollection := client.Database(config.Database.Name).Collection("User")
	LoanCollection := client.Database(config.Database.Name).Collection("Loan")
	err = repository.CreateUniqueIndexes(userCollection)
	if err != nil {
		log.Fatal("Failed to create unique indexes:", err)
	}

	userRepo := repository.NewMongoRepo(userCollection, LoanCollection)
	userUsecase := usecase.NewUsecase(userRepo)
	userController := controller.NewUserController(&userUsecase)
	auth := infrastructures.NewAuthController(userRepo)
	fmt.Println(auth)
	router := router.NewUserRouter(userController,auth)
	router.SetupRoutes()
}
