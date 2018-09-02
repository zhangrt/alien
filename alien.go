package alien

import (
	"errors"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type Application struct {
	port string
	router *Router
	db *gorm.DB
}

var application *Application

var config *Config

func BuildApplication() *Application  {
	port := flag.String("port", "", "http listen port")

	flag.Parse()

	config = new(Config)
	config.init()
	if *port == ""{
		if config.Port == "" {
			*port = ":80"
		} else {
			*port = config.Port
		}
	}

	application = &Application{port: *port, router: NewRouter(), db: nil}
	return application
}

func CreateDbConnPool() (db *gorm.DB, err error) {
	if config != nil {
		db, err = gorm.Open("postgres", config.Conn)
		if err == nil {
			application.db = db
			fmt.Println("数据库连接成功")
		}
	} else {
		err = errors.New("请正确配置数据库连接串")
	}
	return
}

func CloseDbConn()  {
	if application != nil && application.db != nil {
		application.db.Close()
	}
}

func GetApplication() *Application  {
	return application
}

func (app *Application) GetPort() (port string)  {
	port = app.port
	return
}

func (app *Application) GetRouter() (router *Router)  {
	router = app.router
	return
}

func (app *Application) Listen()  {
	err := http.ListenAndServe(application.GetPort(), application.GetRouter())
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
