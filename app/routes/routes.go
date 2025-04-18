// app/routes/routes.go
package routes

import (
	"net/http"

	"app/internal/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Router sets up all the routes for the application
func Router(app *gin.Engine) {

	// Health check endpoint
	app.GET("/healthz", func(ctx *gin.Context) {
		logger.Infof("Health check passed")
		ctx.JSON(http.StatusOK, gin.H{"status": "Health check passed.", "message": "Welcome to Project-k API."})
	})

	// Middleware
	app.Use(otelgin.Middleware(viper.GetString("APP_NAME")))
	app.Use(cors.New(cors.Config{
		AllowAllOrigins:        true,
		AllowMethods:           []string{"*"},
		AllowHeaders:           []string{"*"},
		AllowCredentials:       true,
		AllowWildcard:          true,
		AllowBrowserExtensions: true,
		AllowWebSockets:        true,
		AllowFiles:             false,
	}))

	// Create a new group for /api/v1
	apiV1 := app.Group("/api/v1")

	// Define groups of routes under /api/v1
	Product(apiV1.Group("/products"))
	User(apiV1.Group("/users"))
	Room(apiV1.Group("/rooms"))
	Building((apiV1.Group("/buildings")))
	Building_Room((apiV1.Group("/buildingRooms")))
	Booking((apiV1.Group("/bookings")))
}
