package main

import (
	"fmt"
	"strconv"

	echo "github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"github.com/swaggo/echo-swagger"

	_ "github.com/turnkey-commerce/knowledge-keeper/docs"

	"github.com/turnkey-commerce/knowledge-keeper/config"
	"github.com/turnkey-commerce/knowledge-keeper/database"
	"github.com/turnkey-commerce/knowledge-keeper/handlers"
	"github.com/turnkey-commerce/knowledge-keeper/middlewares"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:5000
// @BasePath /
func main() {
	// Set the file name of the configurations file
	viper.SetConfigName("config")
	// Set the path to look for the configurations file
	viper.AddConfigPath(".")
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
	viper.SetConfigType("toml")
	var conf config.Configurations

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	// Set undefined variables
	viper.SetDefault("database.dbname", "knowledge-keeper")

	err := viper.Unmarshal(&conf)
	if err != nil {
		fmt.Printf("Unable to decode configurations into struct, %v", err)
	}

	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	db := database.New(conf.Database.DBName, conf.Database.Server, conf.Database.DBUser, conf.Database.DBPassword)
	fmt.Println(db.Stats().InUse)

	middlewares.UseMiddleware(e)

	h := handlers.NewHandler(db, conf.Server.Secret)

	h.GetRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(conf.Server.Port)))
}
