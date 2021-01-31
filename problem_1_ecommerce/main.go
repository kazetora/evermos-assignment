package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/configs"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/database"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/routers"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/taskqueue"
)

func main() {

	// read config file
	configs.InitConfig("app_config")

	// setup databse connection
	database.Setup()

	// setup redis connection
	// storage.Setup()

	// setup task queue
	taskqueue.Setup()

	// start consumer queues
	go startConsumerQueues()

	// start API Service
	startAPIService()

}

func startAPIService() {
	r := routers.InitRouter()

	fmt.Println("Server listen at: " + os.Getenv("APP_PORT"))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("APP_PORT"), r))

}

func startConsumerQueues() {
	if err := taskqueue.StartConsumers(); err != nil {
		panic(err.Error())
	}
}
