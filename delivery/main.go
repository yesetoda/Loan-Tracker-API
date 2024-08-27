package main

import (
	"context"
	"example/b/Loan-Tracker-API/controller"
	// "example/b/Loan-Tracker-API/infrastructures"
	"example/b/Loan-Tracker-API/repository"
	"example/b/Loan-Tracker-API/router"
	"example/b/Loan-Tracker-API/usecase"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	userCollection := client.Database("Loan").Collection("User")
	err = repository.CreateUniqueIndexes(userCollection)
	if err != nil {
		log.Fatal("Failed to create unique indexes:", err)
	}

	userRepo := repository.NewMongoRepo(userCollection)
	userUsecase := usecase.NewUsecase(userRepo)
	userController := controller.NewUserController(&userUsecase)
	// auth := infrastructures.NewAuthController(userRepo)
	router := router.NewUserRouter(userController)
	router.SetupRoutes()
}
