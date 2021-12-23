package models

import (
	"context"

	"github.com/108356037/v1/strategy-manager/internal/database/mongo"
	log "github.com/sirupsen/logrus"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateStrategy(userid, strategy string) error {

	existStrategy := mongo.StrategyDB.FindOne(context.Background(), bson.M{"strategy_name": strategy, "user_id": userid})

	// strategy already exists, therefore update
	if existStrategy.Err() == nil {
		result := mongo.StrategyDB.FindOneAndUpdate(
			context.Background(),
			bson.M{"strategy_name": strategy, "user_id": userid},
			bson.M{"$set": bson.M{
				"updated_at":  time.Now().Format("2006-01-02 15:04:05"),
				"schedule":    "",
				"cpu_request": "100m",
				"mem_request": "128Mi",
				"cpu_limit":   "750m",
				"mem_limit":   "512Mi"}},
			&options.FindOneAndUpdateOptions{})
		if result.Err() != nil {
			return result.Err()
		}
		bsonRaw, _ := result.DecodeBytes()
		log.Infof("Updated strategy object %v", bsonRaw)

		// strategy not found, therefore create
	} else {
		data := StrategyDoc{
			UserID:       userid,
			StrategyName: strategy,
			CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
			CpuRequest:   "100m",
			MemRequest:   "128Mi",
			CpuLimit:     "750m",
			MemLimit:     "512Mi",
		}
		result, err := mongo.StrategyDB.InsertOne(context.Background(), data)
		if err != nil {
			return err
		}
		log.Infof("Inserted strategy object %v", result.InsertedID)
	}

	return nil
}

func ListUserStrategies(userid string) (*[]StrategyDoc, error) {
	result := &[]StrategyDoc{}

	cursor, err := mongo.StrategyDB.Find(context.Background(), bson.M{"user_id": userid})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetSingleStrategyById(userid, strategyId string) (*StrategyDoc, error) {
	var result StrategyDoc

	objID, _ := primitive.ObjectIDFromHex(strategyId)
	queryResult := mongo.StrategyDB.FindOne(context.Background(), bson.M{"user_id": userid, "_id": objID})
	if queryResult.Err() != nil {
		return nil, queryResult.Err()
	}

	if err := queryResult.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func GetSingleStrategy(userid, strategy string) (*StrategyDoc, error) {
	var result StrategyDoc

	queryResult := mongo.StrategyDB.FindOne(context.Background(), bson.M{"user_id": userid, "strategy_name": strategy})
	if queryResult.Err() != nil {
		return nil, queryResult.Err()
	}

	if err := queryResult.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func DeleteUserStrategy(userid, strategy string) (*StrategyDoc, error) {
	var result StrategyDoc

	queryResult := mongo.StrategyDB.FindOneAndDelete(context.Background(), bson.M{"strategy_name": strategy, "user_id": userid})
	if queryResult.Err() != nil {
		return nil, queryResult.Err()
	}

	if err := queryResult.Decode(&result); err != nil {
		return nil, err
	}
	log.Info("Deleted strategy %s", result.ID)

	return &result, nil
}

// delete all strategies within a user
func DeleteUserStrategies(userid string) error {
	deleteResult, err := mongo.StrategyDB.DeleteMany(context.Background(), bson.M{"user_id": userid})
	if err != nil {
		return err
	}
	log.Infof("Deleted %v documents within user %s", deleteResult.DeletedCount, userid)
	return nil
}

//TODO: Update Strategy
func UpdateUserStrategy(userid, strategy string, updateInfo map[string]interface{}) (*StrategyDoc, error) {
	var result StrategyDoc
	updateBson := bson.M{}
	for k, v := range updateInfo {
		updateBson[k] = v
	}

	queryResult := mongo.StrategyDB.FindOneAndUpdate(
		context.Background(),
		bson.M{"strategy_name": strategy, "user_id": userid},
		bson.M{"$set": updateBson},
		&options.FindOneAndUpdateOptions{},
	)
	if queryResult.Err() != nil {
		return nil, queryResult.Err()
	}

	if err := queryResult.Decode(&result); err != nil {
		return nil, err
	}
	log.Info("Updated strategy %s", result.ID)

	return &result, nil
}
