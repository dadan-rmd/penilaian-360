package server

import (
	"central-auth/internal/app/handler"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/cast"
)

func Router(opt handler.HandlerOption) *gin.Engine {
	healthyHandler := handler.HealthyCheckHandler{
		HandlerOption: opt,
	}

	authHandler := handler.AuthHandler{
		HandlerOption: opt,
	}

	userHandler := handler.UserHandler{
		HandlerOption: opt,
	}

	setMode := cast.ToBool(os.Getenv("DEBUG_MODE"))
	if setMode {
		gin.SetMode(gin.ReleaseMode)
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		gin.SetMode(gin.DebugMode)
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	corsOrigin := strings.Split(os.Getenv("CORS_HEADER"), ",")
	//routes
	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins: corsOrigin,
		// AllowAllOrigins:        true,
		AllowMethods:           []string{"POST", "DELETE", "GET", "OPTIONS", "PUT"},
		AllowHeaders:           []string{"Origin", "Content-Type", "Authorization", "userid", "REQUEST-ID", "X-SIGNATURE", "Referer", "User-Agent"},
		AllowCredentials:       true,
		ExposeHeaders:          []string{"Content-Length"},
		MaxAge:                 120 * time.Second,
		AllowWildcard:          true,
		AllowBrowserExtensions: true,
		AllowWebSockets:        true,
		AllowFiles:             true,
	}))

	r.Use(gin.Recovery())

	//Maximum memory limit for Multipart forms
	// r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.MaxMultipartMemory = 100 * 1024 * 1024 // 100MB

	apiGroup := r.Group("/api/v1")
	{
		apiGroup.GET("/healthy-check", healthyHandler.HealthyCheck)

		authGroup := apiGroup.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/forgot-password", authHandler.ForgotPass)
			authGroup.POST("/register", opt.AuthMiddleware.BasicAuthenticate(), userHandler.AddUser)
		}
	}

	return r
}
