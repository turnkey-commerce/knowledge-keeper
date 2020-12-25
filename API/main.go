package main

import (
	"fmt"
	"strconv"

	echo "github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"

	_ "github.com/turnkey-commerce/knowledge-keeper/docs"

	"github.com/turnkey-commerce/knowledge-keeper/config"
	"github.com/turnkey-commerce/knowledge-keeper/database"
	"github.com/turnkey-commerce/knowledge-keeper/handlers"
	"github.com/turnkey-commerce/knowledge-keeper/middlewares"
)

// @title Knowledge Keeper API
// @version 1.0
// @description API docs for the Knowledge Keeper application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:5000
// @BasePath /
func main() {
	conf, err := config.GetConfig(".")
	if err != nil {
		fmt.Printf("Unable to load the configuration, %v", err)
	}

	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	db := database.New(database.GetConnectionString(conf.Database))
	fmt.Println(db.Stats().InUse)

	err = database.Seed(db, conf.SeedUser.Email, conf.SeedUser.InitialPassword)
	if err != nil {
		fmt.Printf("Unable to seed the database, %v", err)
	}

	middlewares.UseMiddleware(e)

	h := handlers.NewHandler(db, conf.Server.Secret)

	h.GetRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(conf.Server.Port)))
}
