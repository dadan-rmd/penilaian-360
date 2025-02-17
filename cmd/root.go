package cmd

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"penilaian-360/config"
	"penilaian-360/internal/app/appcontext"
	"penilaian-360/internal/app/commons"
	"penilaian-360/internal/app/middleware/authMiddleware"
	"penilaian-360/internal/app/repository"
	"penilaian-360/internal/app/repository/departmentRepository"
	"penilaian-360/internal/app/repository/employeeRepository"
	"penilaian-360/internal/app/repository/evaluatedEmployeeRepository"
	"penilaian-360/internal/app/repository/evaluationAnswerRepository"
	"penilaian-360/internal/app/repository/evaluationRepository"
	"penilaian-360/internal/app/repository/evaluatorEmployeeRepository"
	"penilaian-360/internal/app/repository/questionRepository"
	"penilaian-360/internal/app/repository/userRepository"
	"penilaian-360/internal/app/server"
	"penilaian-360/internal/app/service"
	"penilaian-360/internal/app/service/departmentService"
	"penilaian-360/internal/app/service/employeeService"
	"penilaian-360/internal/app/service/evaluationService"
	"penilaian-360/internal/app/service/formHistoryService"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	gologger "github.com/mo-taufiq/go-logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kbrprime-api",
	Short: "A brief description of your application",
	Long:  `A longer description that spans multiple lines and likely contains examples and usage of using your application.`,
	Run: func(cmd *cobra.Command, args []string) {
		loadEnv("")
		start()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
}

func initCommonOptions() (options commons.Options, err error) {
	cfg := config.Config()
	app := appcontext.NewAppContext(cfg)

	logLevel := zerolog.InfoLevel
	logLevelP, err := zerolog.ParseLevel(os.Getenv("APP_LOG_LEVEL"))
	if err == nil {
		logLevel = logLevelP
	}
	zerolog.SetGlobalLevel(logLevel)

	validator := validator.New()

	var mysqlDB *gorm.DB
	if app.GetMysqlOption(appcontext.DBDialectMysql).IsEnable {
		mysqlDB, err = app.GetDBInstance(appcontext.DBDialectMysql)
		if err != nil {
			log.Info().Msgf("Failed to start, error connect to DB MySQL | %v", err)
			return
		}
	}

	ossClient := app.GetClientOss()

	options = commons.Options{
		AppCtx:    app,
		Db:        mysqlDB,
		UUID:      commons.NewUuid(),
		Validator: validator,
		OssClient: ossClient,
	}

	return
}

func loadEnv(envName string) {
	gologger.LogConf.NestedLocationLevel = 2
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:     os.Stderr,
			NoColor: false,
		},
	)

	dotenvPath := "/params/.env"

	if envName == "test" {
		dotenvPath = "/params/.env.test"
	}

	_, file, _, _ := runtime.Caller(0)
	rootPath := path.Join(file, "..", "..") + dotenvPath
	log.Info().Msg("path env =>" + rootPath)
	err := godotenv.Load(rootPath)
	if err != nil {
		log.Error().Msg("Error loading .env file")
	}
}

func start() {
	opt, err := initCommonOptions()
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}

	repo := wiringRepository(repository.Option{
		Options: opt,
	})

	service := wiringService(service.Option{
		Options:      opt,
		Repositories: repo,
	})

	server := server.NewServer(opt, service)

	// run app
	server.StartApp()
}
func wiringRepository(repoOption repository.Option) *repository.Repositories {
	repo := repository.Repositories{
		UserRepository:              userRepository.NewUserRepository(repoOption.Db),
		EvaluationRepository:        evaluationRepository.NewEvaluationRepository(repoOption.Db),
		QuestionRepository:          questionRepository.NewQuestionRepository(repoOption.Db),
		EvaluationAnswerRepository:  evaluationAnswerRepository.NewEvaluationAnswerRepository(repoOption.Db),
		DepartmentRepository:        departmentRepository.NewDepartmentRepository(repoOption.Db),
		EmployeeRepository:          employeeRepository.NewEmployeeRepository(repoOption.Db),
		EvaluatorEmployeeRepository: evaluatorEmployeeRepository.NewEvaluatorEmployeeRepository(repoOption.Db),
		EvaluatedEmployeeRepository: evaluatedEmployeeRepository.NewEvaluatedEmployeeRepository(repoOption.Db),
	}

	return &repo
}

func wiringService(serviceOption service.Option) *service.Services {
	// trx := transaction.NewTransaction(serviceOption.Db)
	svc := service.Services{
		AuthMiddleware:     authMiddleware.NewAuthMiddleware(serviceOption.EmployeeRepository),
		DepartmentService:  departmentService.NewDepartmentService(serviceOption.DepartmentRepository),
		EmployeeService:    employeeService.NewEmployeeService(serviceOption.EmployeeRepository, serviceOption.EvaluatorEmployeeRepository, serviceOption.EvaluatedEmployeeRepository),
		FormHistoryService: formHistoryService.NewFormHistoryService(serviceOption.Db, serviceOption.EvaluationRepository, serviceOption.QuestionRepository, serviceOption.EmployeeRepository, serviceOption.EvaluationAnswerRepository, serviceOption.EvaluatorEmployeeRepository, serviceOption.EvaluatedEmployeeRepository),
		EvaluationService:  evaluationService.NewEvaluationService(serviceOption.Db, serviceOption.EvaluationRepository, serviceOption.QuestionRepository, serviceOption.EmployeeRepository, serviceOption.EvaluationAnswerRepository, serviceOption.EvaluatorEmployeeRepository, serviceOption.EvaluatedEmployeeRepository),
	}
	return &svc
}
