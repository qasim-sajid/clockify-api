package conf

import (
	"fmt"
	"os"

	"github.com/subosito/gotenv"
)

const (
	TIME_LAYOUT = "2006-01-02T15:05:05"
)

func GetServerAddress() string {
	return fmt.Sprintf(":%v", Configs.UserPort)
}

// Configuration specifies env variables
type Configuration struct {
	UserPort          string
	SigningKey        string
	RefreshSigningKey string
	DBName            string
	DBUser            string
	DBPassword        string
	DBHost            string
	DBPort            string
}

var (
	// Configs can be used gloablly to get env variables
	Configs *Configuration
)

// InitConfigs loads enviornment variables
func InitConfigs() {
	if err := gotenv.Load(); err != nil {
		fmt.Printf("gotenv: could not find .env file - Error: %v\n", err)
	}

	Configs = &Configuration{
		UserPort:          os.Getenv("USER_PORT"),
		SigningKey:        os.Getenv("SIGNING_KEY"),
		RefreshSigningKey: os.Getenv("REFRESH_SIGNING_KEY"),
		DBName:            os.Getenv("DB_NAME"),
		DBUser:            os.Getenv("DB_USER"),
		DBPassword:        os.Getenv("DB_PASSWORD"),
		DBHost:            os.Getenv("DB_HOST"),
		DBPort:            os.Getenv("DB_PORT"),
	}

	validate()
}

func validate() {
	message := "Missing env variable:"
	if Configs.UserPort == "" {
		panic(fmt.Sprintf("%v %v", message, "USER_PORT"))
	} else if Configs.SigningKey == "" {
		panic(fmt.Sprintf("%v %v", message, "SIGNING_KEY"))
	} else if Configs.RefreshSigningKey == "" {
		panic(fmt.Sprintf("%v %v", message, "REFRESH_SIGNING_KEY"))
	} else if Configs.DBName == "" {
		panic(fmt.Sprintf("%v %v", message, "DB_NAME"))
	} else if Configs.DBUser == "" {
		panic(fmt.Sprintf("%v %v", message, "DB_USER"))
	} else if Configs.DBPassword == "" {
		panic(fmt.Sprintf("%v %v", message, "DB_PASSWORD"))
	} else if Configs.DBHost == "" {
		panic(fmt.Sprintf("%v %v", message, "DB_HOST"))
	} else if Configs.DBPort == "" {
		panic(fmt.Sprintf("%v %v", message, "DB_PORT"))
	}
}
