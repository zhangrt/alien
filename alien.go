package alien

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Application struct {
	port string
	router *Router
}

var application *Application

func BuildApplication() *Application  {
	port := flag.String("port", ":80", "http listen port")

	flag.Parse()
	application = &Application{port: *port, router: NewRouter()}
	return application
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

func sayHelloName(ctx *Context)  {
	ctx.Request.ParseForm()
	fmt.Println(ctx.Request.Form)
	fmt.Println("path", ctx.Request.URL.Path)
	fmt.Println("scheme", ctx.Request.URL.Scheme)
	fmt.Println(ctx.Request.Form["url_long"])
	for k, v := range ctx.Request.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(ctx.ResponseWriter, "Hello astaxie")
}
