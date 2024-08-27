package main

import (
	"context"
	"example/b/Loan-Tracker-API/config"
	"example/b/Loan-Tracker-API/controller"
	"fmt"

	// "example/b/Loan-Tracker-API/infrastructures"
	"example/b/Loan-Tracker-API/repository"
	"example/b/Loan-Tracker-API/router"
	"example/b/Loan-Tracker-API/usecase"
	"log"

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
	// auth := infrastructures.NewAuthController(userRepo)
	router := router.NewUserRouter(userController)
	router.SetupRoutes()
}
