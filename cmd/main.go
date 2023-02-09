package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	// "sync"

	_ "io/ioutil"
	_ "net/http"

	"github.com/BkSearch/BkSearch-Backend/api"
	"github.com/BkSearch/BkSearch-Backend/db"
	"github.com/BkSearch/BkSearch-Backend/elasticsearch"
	_ "github.com/BkSearch/BkSearch-Backend/node"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

var (
	port                                                           int
	host, userRead, userWrite, passwordRead, passwordWrite, dbName string
	elastic_search_port                                            string
)

func loadConfig() {
	//load host
	host = os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	// load port
	port, _ = strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 5432
	}
	// read user
	userRead = os.Getenv("UserRead")
	if userRead == "" {
		log.Fatal("Invalid read user")
	}

	passwordRead = os.Getenv("PasswordRead")

	// write user
	userWrite = os.Getenv("UserWrite")
	if userWrite == "" {
		log.Fatal("Invalid write user")
	}
	passwordWrite = os.Getenv("PasswordWrite")
	// load db
	dbName = os.Getenv("DBName")
	elastic_search_port = os.Getenv("ELASTIC_SEARCH_PORT")
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	loadConfig()
	dbWrite, err := db.InitSQLDB(port, host, userWrite, passwordWrite, dbName)
	if err != nil {
		log.Fatal(err)
	}
	itemDB := db.NewItemDB(dbWrite)
	es := elasticsearch.NewStackOverflow([]string{elastic_search_port})

	// newNode := node.NewNode(es, itemDB, 16)
 //  newNode.SynchronizeAnswerData()
	if err != nil {
	  fmt.Println("Error")
	  fmt.Println(err)
	}

	// fmt.Println(data)
	router := gin.Default()
	setup := api.Config{
		Version: "Test",
		Server:  router,
		ES:      es,
		ItemDB:  itemDB,
	}
	_, err = api.NewAPI(setup)
	if err != nil {
		fmt.Println("error")
	}

	router.Run(":5000")
}
