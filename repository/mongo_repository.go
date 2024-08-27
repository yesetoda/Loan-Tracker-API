package repository

import (
	"context"
	"example/b/Loan-Tracker-API/domain"
	"example/b/Loan-Tracker-API/infrastructures/password_service"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepo struct {
	UserCollection *mongo.Collection
}

func NewMongoRepo(uc *mongo.Collection) GeneralRepository {
	return &MongoRepo{
		UserCollection: uc,
	}
}

func (mr *MongoRepo) CreateUser(user domain.User) (domain.User, error) {
	cnt, err := mr.UserCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return domain.User{}, err
	}
	if cnt == 0 {
		user.IsAdmin = true
		user.Verified = true
	}

	_, err = mr.UserCollection.InsertOne(context.TODO(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.User{}, fmt.Errorf("username or email already exists")
		}
	}
	return user, nil
}

func (mr *MongoRepo) FindUserByEmail(email string) (domain.User, error) {
	var user domain.User
	err := mr.UserCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return domain.User{}, fmt.Errorf("no such user")
	}
	return user, nil
}
func (mr *MongoRepo) FinduserById(id string) (domain.User, error) {
	var user domain.User
	err := mr.UserCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return domain.User{}, fmt.Errorf("no such user")
	}
	return user, nil
}
func (mr *MongoRepo) UpdateUser(id string, user domain.User) error {
	result, err := mr.UserCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, user)
	if err != nil {
		return fmt.Errorf("update failed")
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("update failed")

	}
	return nil
}
func (mr *MongoRepo) DeleteUser(id string) error {
	result, err := mr.UserCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("delete failed")
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("delete failed")
	}
	return nil
}
func (mr *MongoRepo) VerifiyUser(id string) error {
	result, err := mr.UserCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"verified": true})
	if err != nil {
		return fmt.Errorf("verification failed")
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("verification failed")

	}
	return nil
}
func (mr *MongoRepo) AuthenticateUser(id string, password string) error {
	var user domain.User
	err := mr.UserCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return fmt.Errorf("invalid credential")
	}
	if err = password_service.CheckPasswordHash(password, user.Password); err != nil {
		return fmt.Errorf("invalid credential")
	}
	return nil
}
func (mr *MongoRepo) ListAllUsers() []domain.User {
	var users []domain.User
	cursor, err := mr.UserCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return users
	}
	for cursor.Next(context.TODO()) {
		var user domain.User
		cursor.Decode(&user)
		users = append(users, user)
	}
	return users
}
func (mr *MongoRepo) ResetPassword(id string, password string) error {
	result, err := mr.UserCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"password": password})
	if err != nil {
		return fmt.Errorf("password Reset failed")
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("password Reset failed")
	}
	return nil
}
