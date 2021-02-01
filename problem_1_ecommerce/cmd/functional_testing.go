package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/configs"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/controllers"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/database"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/models"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/requests"
)

var (
	successCount   []int
	failCount      []int
	inventoryExist bool
	inventoryCount int
)

const (
	testProductID   = 1
	testBuyQuantity = 1
)

func main() {
	testCases := []int{10, 100, 200, 300, 500, 1000}
	successCount = []int{0, 0, 0, 0, 0, 0}
	failCount = []int{0, 0, 0, 0, 0, 0}
	inventoryExist = false

	// read config file
	configs.InitConfig("test_config")

	// setup databse connection
	database.Setup()

	initTest()

	db := database.GetDatabase()

	for idx, testcase := range testCases {
		var wg sync.WaitGroup
		// setup the quantity to be half of request number (test case)
		db.Model(models.Inventories{}).Where("product_id = ?", testProductID).Update("quantity", testcase/2)
		fmt.Printf("Test case %d: %d Requests, %d Inventory Quantity\n", idx+1, testcase, testcase/2)
		for i := 0; i < testcase; i++ {
			wg.Add(1)
			go task(&wg, uint(idx+1), testProductID, testBuyQuantity, idx)
		}
		wg.Wait()

		fmt.Printf("Result: %d success, %d fail\n", successCount[idx], failCount[idx])
		if successCount[idx] <= testcase/2 {
			fmt.Println("SUCCESS")
		} else {
			fmt.Println("FAIL")
		}
	}

	cleanTest()

}

func initTest() {
	db := database.GetDatabase()

	var inventory models.Inventories
	queryRes := db.Where("product_id = ?", testProductID).First(&inventory)
	if queryRes.RecordNotFound() {
		newInventory := models.Inventories{
			ProductID: testProductID,
			Quantity:  0,
		}
		if err := db.Create(&newInventory).Error; err != nil {
			panic(err.Error())
		}
	} else {
		inventoryExist = true
		inventoryCount = inventory.Quantity
	}
}

func cleanTest() {
	db := database.GetDatabase()
	if !inventoryExist {
		db.Where("product_id = ?", testProductID).Delete(models.Inventories{})
	} else {
		db.Model(models.Inventories{}).Where("product_id = ?", testProductID).Update("quantity", inventoryCount)
	}
}

func task(wg *sync.WaitGroup, userID, productID uint, quantity, idx int) error {
	defer wg.Done()

	reqData := requests.AddToCartRequest{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	}

	reqJSON, err := json.Marshal(reqData)
	if err != nil {
		log.Printf("Error in json marshal: %s", err.Error())
		return err
	}

	response, err := http.Post(fmt.Sprintf("%s/api/v1/order/addToCart", os.Getenv("TEST_URL")), "application/json", bytes.NewBuffer(reqJSON))
	if err != nil {
		log.Printf("Error in http post: %s", err.Error())
		return err
	}

	decoder := json.NewDecoder(response.Body)
	transaction := map[string]string{}

	if err := decoder.Decode(&transaction); err != nil {
		log.Printf("Error in decoding response body: %s", err.Error())
		return err
	}

	for {
		response, err := http.Get(fmt.Sprintf("%s/api/v1/transaction/%s", os.Getenv("TEST_URL"), transaction["transaction_key"]))
		if err != nil {
			log.Printf("Error in http get: %s", err.Error())
			return err
		}

		if response.StatusCode != http.StatusOK {
			time.Sleep(1000 * time.Millisecond)
		} else {
			decoder := json.NewDecoder(response.Body)
			res := controllers.TransactionResponse{}

			if err := decoder.Decode(&res); err != nil {
				log.Printf("Error in decoding response body: %s", err.Error())
				return err
			}

			switch res.Status {
			case models.TransactionStatusSuccess:
				successCount[idx] = successCount[idx] + 1
			case models.TransactionStatusError:
				failCount[idx] = failCount[idx] + 1
			}

			break
		}
	}

	return nil
}
