package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongoDaoImpl struct {
	db *mongo.Database
}

var _ UserDao = &UserMongoDaoImpl{}

func WithTimeout(action func(context.Context) any) (result any, failedResult error) {
	var context, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer handlePanic(&failedResult)
	defer cancel()
	return action(context), failedResult
}

func handlePanic(failedResult *error) {
    if r := recover(); r != nil {
        *failedResult = r.(error)
    }
}

// CreateUser implements UserDao.
func (u *UserMongoDaoImpl) CreateUser(user User) (int, error) {
	result, err := WithTimeout(func(context context.Context) any {
		collection := u.createUserCollection(context)
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		insertResult, err := collection.InsertOne(context, user)
		if err != nil {
			panic(err)
		}

		// MongoDB 的 InsertedID 通常是 primitive.ObjectID 型別，需轉換為 int64 或適當處理
		if oid, ok := insertResult.InsertedID.(primitive.ObjectID); ok {
			// 將 ObjectID 轉換為 int64（取前8位）
			return int64(oid.Timestamp().Unix())
		}
		return -1
	})

	if err == nil {
		if id, ok := result.(int64); ok {
			return int(id), nil
		}
	} 
	return -1, err
}

func (u *UserMongoDaoImpl) createUserCollection(context context.Context) *mongo.Collection {
	collection := u.db.Collection("users")
	if collection == nil {
		failed := u.db.CreateCollection(context, "users")
		log.Printf("Collection creation result: %v", failed)
	}
	return collection
}

// DeleteUser implements UserDao.
// TODO: the id is not the primary key, it is the object id
func (u *UserMongoDaoImpl) DeleteUser(id int) error {
	result, err := WithTimeout(func(context context.Context) any {
		collection := u.createUserCollection(context)
		filter := bson.M{"_id": id}
		_, err := collection.DeleteOne(context, filter)
		return err
	})
	if err == nil {
		return result.(error)
	}
	return err
}

// FindUserByName implements UserDao.
func (u *UserMongoDaoImpl) FindUserByName(name string) (*User, error) {
	result, err := WithTimeout(func(context context.Context) any {
		collection := u.createUserCollection(context)
		filter := bson.M{"name": name}
		user := collection.FindOne(context, filter)
		return user.Decode(&User{})
	})
	if err == nil {
		return result.(*User), err
	}
	return nil, err
}

// GetUsers implements UserDao.
func (u *UserMongoDaoImpl) GetUsers() ([]User, error) {
	result, err := WithTimeout(func(context context.Context) any {
		collection := u.createUserCollection(context)
		filter := bson.M{}
		cursor, err := collection.Find(context, filter)
		if err != nil {
			panic(err)
		}
		var users []User
		err = cursor.All(context, users)
		if err != nil {
			panic(err)
		}

		return users
	})
	if err == nil {
		return result.([]User), err
	}
	return nil, err
}
