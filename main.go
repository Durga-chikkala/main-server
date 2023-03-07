package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/main-server/handlers/query"
	"github.com/main-server/handlers/user"
	middlewares "github.com/main-server/middleware"
	serviceQuery "github.com/main-server/services/query"
	serviceUser "github.com/main-server/services/user"

	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	PORT := ":8000"

	client := resty.New()

	//query Injection
	queryService := serviceQuery.New(client)
	queryHandler := query.New(queryService)

	//User Injection
	userService := serviceUser.New(client)
	userHandler := user.New(userService)

	app.Use(middlewares.CORSMiddleware())
	app.POST("/user/signup", userHandler.Create)
	app.GET("/user/login", userHandler.Get)

	app.Use(middlewares.Auth())

	//Questions Endpoints
	app.POST("/chatbot", queryHandler.Create)
	app.GET("/chatbot", queryHandler.Get)
	app.GET("/chatbot/:question", queryHandler.GetByQuestion)
	app.GET("/chatbot/frequentQuestions", queryHandler.GetFrequentQuestions)
	app.PATCH("/chatbot/:question", queryHandler.PatchByQuestion)

	// User Details
	app.GET("user/logout", userHandler.Logout)
	app.GET("user/:id", userHandler.GetByID)
	app.PATCH("user/:id", userHandler.PatchByID)

	log.Printf("The server is running at port:%v", PORT)
	err := app.Run(PORT)
	if err != nil {
		log.Printf("Port is already in use")
		return
	}
}
