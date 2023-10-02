package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct{
	DB *gorm.DB
	Router *mux.Router
}

// struct pada .env
type AppConfig struct{
	AppName string
	AppEnv string
	AppPort string
}

// struct .env database
type DBConfig struct{
	DBHost string
	DBUser string
	DBPassword string
	DBName string
	DBPort string
}

// method Initialize
func (server *Server) Initialize(appConfig AppConfig, dbConfig DBConfig) {
	fmt.Println("Welcome to " + appConfig.AppName)

	// connect ke database postgresSQL
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)
	server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic ("Failed connection in database server")
	}

	
	server.Router = mux.NewRouter()
	server.initializeRoutes()
}

// method Run
func (server *Server) Run(addr string){
	
	fmt.Printf("Listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

//func pemberian default value pada .env
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok{
		return value
	}

	return fallback
}

// func pemanggilan pertama kali dieksekusi
func Run(){
	var server = Server{}
	var appConfig = AppConfig{}
	var dbconfig = DBConfig{}

	if err := godotenv.Load(); err != nil{
		log.Fatal("Error on loading .env file")
	}

	// Penggunaan .env
	appConfig.AppName = getEnv("APP_NAME", "GoTokoWeb")
	appConfig.AppEnv = getEnv("APP_ENV", "development")
	appConfig.AppPort = getEnv("APP_PORT", "9000")

	// .env pada database
	dbconfig.DBHost = getEnv("DB_HOST", "localhost")
	dbconfig.DBUser = getEnv("DB_USER", "vega")
	dbconfig.DBPassword = getEnv("DB_PASSWORD", "password")
	dbconfig.DBName = getEnv("DB_NAME", "latihan_gotokodb")
	dbconfig.DBPort = getEnv("DB_PORT", "5432")

	server.Initialize(appConfig, dbconfig)
	server.Run(":" + appConfig.AppPort)
}

