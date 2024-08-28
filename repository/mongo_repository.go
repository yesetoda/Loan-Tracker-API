package repository

import (
	"context"
	"errors"
	"example/b/Loan-Tracker-API/domain"
	"example/b/Loan-Tracker-API/infrastructures/password_service"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepo struct {
	UserCollection *mongo.Collection
	LoanCollection *mongo.Collection
}

func NewMongoRepo(uc *mongo.Collection, lc *mongo.Collection) GeneralRepository {
	return &MongoRepo{
		UserCollection: uc,
		LoanCollection: lc,
	}
}

func IsValidObjectID(id string) (primitive.ObjectID, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return oid, nil
}

func (mr *MongoRepo) CreateUser(user domain.User) (domain.User, error) {
	// Set admin status if this is the first user
	count, err := mr.UserCollection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		return domain.User{}, err
	}
	if count == 0 {
		user.IsAdmin = true
	}

	_, err = mr.UserCollection.InsertOne(context.TODO(), user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.User{}, fmt.Errorf("username or email already exists")
		}
		return domain.User{}, err
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

func (mr *MongoRepo) FindUserById(id string) (domain.User, error) {
	objectID, err := IsValidObjectID(id)
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	err = mr.UserCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return domain.User{}, fmt.Errorf("no such user")
	}
	return user, nil
}

func (mr *MongoRepo) UpdateUser(id string, user domain.User) error {
	objectID, err := IsValidObjectID(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"email":            user.Email,
			"first_name":       user.FirstName,
			"last_name":        user.LastName,
			"verified":         user.Verified,
			"is_admin":         user.IsAdmin,
			"updated_at":       time.Now(),
			"verify_token":     user.VerifyToken,
			"verify_token_exp": user.VerfyTokenExp,
		},
	}
	if user.Password != "" {
		update = bson.M{
			"$set": bson.M{
				"email":            user.Email,
				"password":         user.Password,
				"first_name":       user.FirstName,
				"last_name":        user.LastName,
				"verified":         user.Verified,
				"is_admin":         user.IsAdmin,
				"updated_at":       time.Now(),
				"verify_token":     user.VerifyToken,
				"verify_token_exp": user.VerfyTokenExp,
			},
		}
	}

	result, err := mr.UserCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("update failed: no documents modified")
	}
	return nil
}

func (mr *MongoRepo) DeleteUser(id string) error {
	objectID, err := IsValidObjectID(id)
	if err != nil {
		return err
	}

	result, err := mr.UserCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("delete failed: no documents deleted")
	}
	return nil
}

func (mr *MongoRepo) VerifyUser(id string) error {
	objectID, err := IsValidObjectID(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"verified":   true,
			"updated_at": time.Now(),
		},
	}

	result, err := mr.UserCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("verification failed: no documents modified")
	}
	return nil
}

func (mr *MongoRepo) AuthenticateUser(id string, password string) error {
	objectID, err := IsValidObjectID(id)
	if err != nil {
		return err
	}

	var user domain.User
	err = mr.UserCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return fmt.Errorf("invalid credential: %w", err)
	}

	if err = password_service.CheckPasswordHash(password, user.Password); err != nil {
		return fmt.Errorf("invalid credential: %w", err)
	}

	return nil
}

func (mr *MongoRepo) ListAllUsers() ([]domain.User, error) {
	var users []domain.User
	cursor, err := mr.UserCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var user domain.User
		if err = cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (mr *MongoRepo) ResetPassword(id string, password string) error {
	objectID, err := IsValidObjectID(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"password":   password,
			"updated_at": time.Now(),
		},
	}

	result, err := mr.UserCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, update)
	if err != nil {
		return fmt.Errorf("password reset failed: %w", err)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("password reset failed: no documents modified")
	}
	return nil
}

func (mr *MongoRepo) CreateLoan(loan domain.Loan) (domain.Loan, error) {
	loan.ID = primitive.NewObjectID().Hex()
	loan.CreatedAt = time.Now()
	loan.UpdatedAt = time.Now()

	_, err := mr.LoanCollection.InsertOne(context.TODO(), loan)
	if err != nil {
		return domain.Loan{}, err
	}
	return loan, nil
}

func (mr *MongoRepo) FindLoanByID(id string) (domain.Loan, error) {
	objectID, err := IsValidObjectID(id)
	if err != nil {
		return domain.Loan{}, err
	}
	fmt.Println("this is the loanid", objectID)
	var loan domain.Loan
	err = mr.LoanCollection.FindOne(context.TODO(), bson.M{"_id": objectID.Hex()}).Decode(&loan)
	if err != nil {
		return domain.Loan{}, err
	}
	return loan, nil
}

func (mr *MongoRepo) FindLoansByUserID(userID string) ([]domain.Loan, error) {
	var loans []domain.Loan
	cursor, err := mr.LoanCollection.Find(context.TODO(), bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var loan domain.Loan
		if err = cursor.Decode(&loan); err != nil {
			return nil, err
		}
		loans = append(loans, loan)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return loans, nil
}

func (mr *MongoRepo) FindAllLoans(status, order string) ([]domain.Loan, error) {
	var loans []domain.Loan
	fmt.Println("FindAllLoans", status, order)
	findOptions := options.Find()
	if order == "desc" {
		findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	} else {
		findOptions.SetSort(bson.D{{Key: "created_at", Value: 1}})
	}

	filter := bson.M{}
	if status != "all" {
		filter["status"] = status
	}

	cursor, err := mr.LoanCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var loan domain.Loan
		if err = cursor.Decode(&loan); err != nil {
			return nil, err
		}
		loans = append(loans, loan)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return loans, nil
}

func (mr *MongoRepo) UpdateLoanStatus(id string, status string) (string, error) {
	objectID, err := IsValidObjectID(id)
	if err != nil {
		return "Update failed", err
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	result, err := mr.LoanCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID.Hex()}, update)
	if err != nil {
		return "Update failed", err
	}
	if result.ModifiedCount == 0 {
		return "Update failed: no documents modified", errors.New("update failed: no documents updated")
	}
	return "Loan status updated successfully", nil
}

func (mr *MongoRepo) DeleteLoan(id string) (string, error) {
	objectID, err := IsValidObjectID(id)
	if err != nil {
		return "Delete failed", err
	}

	result, err := mr.LoanCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID.Hex()})
	if err != nil {
		return "Delete failed", err
	}
	fmt.Println(result.DeletedCount, err)
	if result.DeletedCount == 0 {
		return " no documents deleted", errors.New("delete failed: no documents deleted")
	}
	return "Deleted loan with ID " + id, nil
}
