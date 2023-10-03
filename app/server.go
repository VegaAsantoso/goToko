package app

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/VegaASantoso/goToko/database/seeders"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
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
	// server.initializeDB(dbConfig)
	
	server.initializeRoutes()

	// seeders
	// seeders.DBSeed(server.DB)

}


// Method initialize database
func (server *Server) initializeDB(dbConfig DBConfig){
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)
	server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic ("Failed connection in database server")
	}

}

func (server *Server) dbMigrate(){
	for _, model := range RegisterModels(){
		err := server.DB.Debug().AutoMigrate(model.Model)

		if err != nil{
			log.Fatal(err)
		}
	}

	fmt.Println("Database migrated successfully")
}

// pemisahan menggunakan cli 
func (server *Server) initComands(config AppConfig, dbConfig DBConfig){
	server.initializeDB(dbConfig)

	cmdApp := cli.NewApp()
	cmdApp.Commands = []cli.Command{
		{
			Name : "db:migrate",
			Action : func(c *cli.Context) error {
				server.dbMigrate()
				return nil
			},
		},
		{
			Name : "db:seed",
			Action: func(c *cli.Context) error {
				err := seeders.DBSeed(server.DB)
				if err != nil{
					log.Fatal(err)
				}
				return nil
			},
			
		},
	}

	if err := cmdApp.Run(os.Args); err != nil{
		log.Fatal(err)
	}
}

//func pemberian default value pada .env
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok{
		return value
	}

	return fallback
}


// method Run
func (server *Server) Run(addr string){
	
	fmt.Printf("Listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}


// func pemanggilan pertama kali dieksekusi
func Run(){
	var server = Server{}
	var appConfig = AppConfig{}
	var dbConfig = DBConfig{}

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
		server.initComands(appConfig, dbConfig)
	}else{
		server.Initialize(appConfig, dbConfig)
		server.Run(":" + appConfig.AppPort)
	}

}

