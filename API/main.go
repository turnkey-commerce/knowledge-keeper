package main

import (
	"fmt"

	echo "github.com/labstack/echo/v4"

	"github.com/turnkey-commerce/knowledge-keeper/database"
	"github.com/turnkey-commerce/knowledge-keeper/handlers"
	"github.com/turnkey-commerce/knowledge-keeper/middlewares"
)

func main() {
	e := echo.New()

	db := database.New()
	fmt.Println(db.Stats().InUse)

	middlewares.UseMiddleware(e)

	h := handlers.NewHandler(db)

	h.GetRoutes(e)

	// Start server
	e.Logger.Fatal(e.Start(":5000"))
}
