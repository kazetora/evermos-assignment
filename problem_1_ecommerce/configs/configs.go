package configs

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// InitConfig initialize database connection and return the connection object
func InitConfig(filename string) {

	if err := setEnv(filename); err != nil {
		log.Printf("Error loading config: %s", err.Error())
		panic("Failed to load config")
	}

}

func setEnv(filename string) error {
	vi := viper.New()

	vi.SetConfigName(filename)
	vi.AddConfigPath("./configs")
	vi.AutomaticEnv()

	err := vi.ReadInConfig()
	if err != nil {
		return err
	}

	configs := []string{"app_config", "db_config", "redis_config", "test_config"}
	for _, conf := range configs {
		if err = setupConfig(vi, conf); err != nil {
			return err
		}
	}
	return nil
}

func setupConfig(vi *viper.Viper, name string) error {
	conf := vi.GetStringMapString(name)
	for k := range conf {
		key := strings.ToUpper(k)
		fmt.Printf("%s, %s\n", key, conf[k])
		err := os.Setenv(key, conf[k])
		if err != nil {
			return err
		}
	}

	return nil
}
