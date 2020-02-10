package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql database driver
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver

	"github.com/jongschneider/fullstack/api/models"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var DBURL string
	switch Dbdriver {
	case "mysql":
		DBURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
	case "postgres":
		DBURL = fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", Dbdriver, DbUser, DbPassword, DbHost, DbPort, DbName)
		// DBURL = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", DbHost, DbPort, DbUser, DbName, DbPassword)
	default:
		log.Fatal("Invalid DB driver type")
	}

	var err error
	server.DB, err = gorm.Open(Dbdriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", Dbdriver)
		spew.Dump(DBURL)
		log.Fatal("This is the error:", err)
	}

	fmt.Printf("We are connected to the %s database", Dbdriver)

	server.DB.Debug().AutoMigrate(&models.User{}, &models.Post{}) //database migration

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
