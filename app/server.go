package app

import (
	"flag"
	"log"
	"os"

	"github.com/VegaASantoso/goToko/app/controllers"
	"github.com/joho/godotenv"
)

//func pemberian default value pada .env
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok{
		return value
	}

	return fallback
}

// func pemanggilan pertama kali dieksekusi
func Run(){
	var server = controllers.Server{}
	var appConfig = controllers.AppConfig{}
	var dbConfig = controllers.DBConfig{}

	if err := godotenv.Load(); err != nil{
		log.Fatal("Error on loading .env file")
	}

	// Penggunaan .env
	appConfig.AppName = getEnv("APP_NAME", "GoTokoWeb")
	appConfig.AppEnv = getEnv("APP_ENV", "development")
	appConfig.AppPort = getEnv("APP_PORT", "9000")

	// .env pada database
	dbConfig.DBHost = getEnv("DB_HOST", "localhost")
	dbConfig.DBUser = getEnv("DB_USER", "vega")
	dbConfig.DBPassword = getEnv("DB_PASSWORD", "password")
	dbConfig.DBName = getEnv("DB_NAME", "latihan_gotokodb")
	dbConfig.DBPort = getEnv("DB_PORT", "5432")

	flag.Parse()
	arg := flag.Arg(0)
	if arg != ""{
		server.InitComands(appConfig, dbConfig)
	}else{
		server.Initialize(appConfig, dbConfig)
		server.Run(":" + appConfig.AppPort)
	}

}

