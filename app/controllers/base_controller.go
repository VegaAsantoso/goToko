package controllers

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/VegaASantoso/goToko/app/models"
	"github.com/VegaASantoso/goToko/database/seeders"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/urfave/cli"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
	AppConfig *AppConfig
}

// struct pada .env
type AppConfig struct {
	AppName string
	AppEnv  string
	AppPort string
	AppURL string
}

// struct .env database
type DBConfig struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
}

// pagination
type PageLink struct{
	Page int32
	Url string
	IsCurrentPage bool
}

type PaginationLinks struct{
	CurrentPage string
	NextPage string
	PrevPage string
	TotalRows int32
	TotalPages int32
	Links []PageLink
}

type PaginationParams struct{
	Path string
	TotalRows int32
	PerPage int32
	CurrentPage int32
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var sessionShoppingCart = "shopping-cart-session"

// method Initialize
func (server *Server) Initialize(appConfig AppConfig, dbConfig DBConfig) {
	fmt.Println("Welcome to " + appConfig.AppName)

	// connect ke database postgresSQL
	server.initializeDB(dbConfig)

	server.initializeAppConfig(appConfig)
	server.initializeRoutes()

	// seeders
	// seeders.DBSeed(server.DB)

}

// method Run
func (server *Server) Run(addr string){
	
	fmt.Printf("Listening to port %s", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

// Method initialize database
func (server *Server) initializeDB(dbConfig DBConfig) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbConfig.DBHost, dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBName, dbConfig.DBPort)
	server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed connection in database server")
	}

}

func(server *Server) initializeAppConfig(appConfig AppConfig){
	server.AppConfig = &appConfig
}

func (server *Server) dbMigrate() {
	for _, model := range models.RegisterModels() {
		err := server.DB.Debug().AutoMigrate(model.Model)

		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Database migrated successfully")
}

// pemisahan menggunakan cli
func (server *Server) InitComands(config AppConfig, dbConfig DBConfig) {
	server.initializeDB(dbConfig)

	cmdApp := cli.NewApp()
	cmdApp.Commands = []cli.Command{
		{
			Name: "db:migrate",
			Action: func(c *cli.Context) error {
				server.dbMigrate()
				return nil
			},
		},
		{
			Name: "db:seed",
			Action: func(c *cli.Context) error {
				err := seeders.DBSeed(server.DB)
				if err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
	}

	if err := cmdApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func GetPaginationLinks(config *AppConfig, params PaginationParams) (PaginationLinks, error){
	var links []PageLink

	totalPages := int32(math.Ceil(float64(params.TotalRows) / float64(params.PerPage)))

	for i := 1; int32(i)<= totalPages; i++{
		links = append(links, PageLink{
			Page: int32(i),
			Url: fmt.Sprintf("%s/%s?page=%s", config.AppURL, params.Path, fmt.Sprint(i)),
			IsCurrentPage: int32(i) == params.CurrentPage,
		})
	}

	var nextPage int32
	var prevPage int32

	prevPage = 1
	nextPage = totalPages

	if params.CurrentPage > 2{
		prevPage = params.CurrentPage - 1
	}
	
	if params.CurrentPage < totalPages{
		nextPage = params.CurrentPage + 1
	}

	return PaginationLinks{
		CurrentPage: fmt.Sprintf("%s/%s?page=%s", config.AppURL, params.Path, fmt.Sprint(params.CurrentPage)),
		NextPage: fmt.Sprintf("%s/%s?page=%s", config.AppURL, params.Path, fmt.Sprint(nextPage)),
		PrevPage: fmt.Sprintf("%s/%s?page=%s", config.AppURL, params.Path, fmt.Sprint(prevPage)),
		TotalRows: params.TotalRows,
		TotalPages: totalPages,
		Links: links,
	}, nil
}