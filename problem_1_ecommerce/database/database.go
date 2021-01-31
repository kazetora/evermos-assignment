package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kazetora/evermos-assignment/problem_1_ecommerce/models"
)

var db *gorm.DB

// Setup setup db connection
func Setup() {

	var err error

	dbSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	db, err = gorm.Open(os.Getenv("DB_TYPE"), dbSource)
	if err != nil {
		log.Println(err.Error())
		panic("Failed to connect database")
	}

	db.LogMode(false)

	// migrate database
	db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&models.Products{},
		&models.Inventories{},
		&models.LineItems{},
		&models.Orders{},
		&models.Transactions{})

}

// GetDatabase get db instance
func GetDatabase() *gorm.DB {
	return db
}
