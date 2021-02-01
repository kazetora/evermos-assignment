package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/models"
	"github.com/ventu-io/go-shortid"
)

const transactionKeyFmt = "evermos:ec:transaction:%d:%s"

// RegisterTransaction registre transaction data in redis
func RegisterTransaction(data models.TransactionCache) (string, error) {

	genKey, err := shortid.Generate()
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	storeKey := fmt.Sprintf(transactionKeyFmt, data.UserID, genKey)

	dataJSON, err := json.Marshal(data)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	if err := redisClient.Set(ctx, storeKey, dataJSON, 180*time.Second).Err(); err != nil {
		log.Println(err.Error())
		return "", err
	}

	return storeKey, nil

}

// GetTransaction retrive transaction from redis
func GetTransaction(retKey string) (models.TransactionCache, error) {

	data, err := redisClient.Get(ctx, retKey).Bytes()
	if err != nil {
		log.Println(err.Error())
		return models.TransactionCache{}, err
	}

	var transaction models.TransactionCache

	err = json.Unmarshal(data, &transaction)
	if err != nil {
		log.Println(err.Error())
		return models.TransactionCache{}, err
	}

	return transaction, nil
}

// UpdateTransaction update transaction from redis
func UpdateTransaction(key string, value models.TransactionCache) error {
	dataJSON, err := json.Marshal(value)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if err := redisClient.Set(ctx, key, dataJSON, 180*time.Second).Err(); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil

}
