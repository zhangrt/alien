package alien

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Application struct {
	port string
	router *Router
}

func BuildApplication() *Application  {
	application := &Application{port: ":80", router: NewRouter()}
	err := http.ListenAndServe(application.GetPort(), application.GetRouter())
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	return application
}

func GetApplication() *Application  {
	return application
}

var application *Application

func (app *Application) GetPort() (port string)  {
	port = app.port
	return
}

func (app *Application) GetRouter() (router *Router)  {
	router = app.router
	return
}

func sayHelloName(ctx *Context)  {
	ctx.r.ParseForm()
	fmt.Println(ctx.r.Form)
	fmt.Println("path", ctx.r.URL.Path)
	fmt.Println("scheme", ctx.r.URL.Scheme)
	fmt.Println(ctx.r.Form["url_long"])
	for k, v := range ctx.r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(ctx.w, "Hello astaxie")
}
