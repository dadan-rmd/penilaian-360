package server

import (
	"os"
	"penilaian-360/internal/app/handler"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/cast"
)

func Router(opt handler.HandlerOption) *gin.Engine {

	departmentHandler := handler.DepartmentHandler{
		HandlerOption: opt,
	}
	employeeHandler := handler.EmployeeHandler{
		HandlerOption: opt,
	}
	formHistoryHandler := handler.FormHistoryHandler{
		HandlerOption: opt,
	}
	evaluationHandler := handler.EvaluationHandler{
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

	apiGroup := r.Group("/api/v1", opt.AuthMiddleware.AuthorizeEmployee())
	// apiGroup := r.Group("/api/v1")
	{
		apiGroup.GET("/departement", departmentHandler.GetDepartment)
		apiGroup.GET("/user", employeeHandler.GetUser)

		employeeGroup := apiGroup.Group("/employee")
		{
			employeeGroup.GET("", employeeHandler.GetEmployee)
			employeeGroup.GET("/emails", employeeHandler.GetEmployeeEmails)
		}

		formHistoryGroup := apiGroup.Group("/form-history")
		{
			formHistoryGroup.GET("", formHistoryHandler.GetFormHistory)
			formHistoryGroup.POST("", formHistoryHandler.FormHistory)
			formHistoryGroup.GET("/:id", formHistoryHandler.FormHistoryView)
			formHistoryGroup.DELETE("/:id", formHistoryHandler.FormHistoryDelete)
			formHistoryGroup.POST("/assignment", formHistoryHandler.FormHistoryAssignment)
			formHistoryGroup.GET("/detail", formHistoryHandler.FormHistoryDetail)
			formHistoryGroup.GET("/copy/:id", formHistoryHandler.FormHistoryCopy)
		}
		evaluationGroup := apiGroup.Group("/evaluation")
		{
			evaluationGroup.GET("/list", evaluationHandler.EvaluationList)
			evaluationGroup.GET("/need", evaluationHandler.EvaluationNeeds)
			evaluationGroup.GET("/divisi", evaluationHandler.EvaluationDepartementList)
			evaluationGroup.GET("/:employee_id", evaluationHandler.EvaluationDetail)
			evaluationGroup.GET("approve/:evaluator_id", evaluationHandler.Approve)
			evaluationScoreGroup := evaluationGroup.Group("/score")
			{
				evaluationScoreGroup.GET("/detail", evaluationHandler.ScoreDetail)
				evaluationScoreGroup.POST("", evaluationHandler.Score)
			}
		}
	}
	r.GET("/api/v1/create-token", employeeHandler.CreateToken)

	return r
}
