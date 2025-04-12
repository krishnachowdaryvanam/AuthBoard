package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/krishnachowdaryvanam/authboard/auth_service/handlers"
)

func SetupRouter() *gin.Engine {
	// Create a new Gin router
	r := gin.Default()

	// Define routes for authentication (SignUp, Login)
	auth := r.Group("/auth")
	{
		auth.POST("/signup", handlers.SignUp)
		auth.POST("/login", handlers.Login)
	}

	// Add other routes if needed

	return r
}
