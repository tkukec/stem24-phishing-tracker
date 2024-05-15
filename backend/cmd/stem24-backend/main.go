package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/andrezz-b/stem24-phishing-tracker/application"
	"github.com/andrezz-b/stem24-phishing-tracker/domain/services/tenant"
	"github.com/andrezz-b/stem24-phishing-tracker/infrastructure/metrics"
	"github.com/andrezz-b/stem24-phishing-tracker/infrastructure/repositories"
	httpControllers "github.com/andrezz-b/stem24-phishing-tracker/infrastructure/transport/http"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/constants"
	appContext "github.com/andrezz-b/stem24-phishing-tracker/shared/context"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/database"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/health"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/logging"
	"github.com/andrezz-b/stem24-phishing-tracker/shared/middleware"
	runtimebag "github.com/andrezz-b/stem24-phishing-tracker/shared/runtimebag"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"
)

var (
	router   *gin.Engine
	conn     *database.Connection
	writers  logging.LoggerWriters
	amLogger zerolog.Logger

	prometheusRegistry = prometheus.NewRegistry()
	chassisMetrics     = metrics.NewMetrics(prometheusRegistry)
	runtimeBag         *runtimebag.Bag

	// Controllers
	authController    *httpControllers.Auth
	commentController *httpControllers.Comments
	eventController   *httpControllers.Event
	statusController  *httpControllers.Status
)

// @title STEM-24 Git Good Backend service
// @version 1.0.0
// @description This is a service for managing phishing events

// @contact.name Kristijan Fajdetić
// @contact.email Kristijan.Fajdetic@asseco-see.hr

// @host localhost
// @query.collection.format multi
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /api
func main() {
	log.Println("Booting ....")
	log.Println("Loading .env if exists ....")
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file... assuming production ...")
	}

	buildLogs()
	//authentication()
	databaseConnection()
	migrations()
	buildDependencies()
	seeding()

	ctx, cancel := context.WithCancel(context.Background())

	//apply graceful shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", runtimebag.GetEnvString("HTTP_PORT", "8080")),
		Handler: httpRouter(),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// pprof for memory|CPU sampling
	var pprofSrv *http.Server
	if runtimebag.GetEnvBool("PPROF_ENABLED", false) {
		pprofSrv = &http.Server{
			Addr: ":" + runtimebag.GetEnvString("PPROF_PORT", "6061"),
		}
		go func() {
			if err = pprofSrv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
				log.Fatalf("pprof listen: %s\n", err)
			}
		}()
	}

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down http server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	if pprofSrv != nil {
		log.Println("Shutting down pprof http server...")
		if err := pprofSrv.Shutdown(ctx); err != nil {
			log.Fatal("pprof server forced to shutdown:", err)
		}
	}

	cancel()
}

func createLogsDirectory() {
	log.Println("Creating logs directory")
	path := "logs"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.Mkdir(path, 0755); err != nil { // 0755 sets permissions for the directory
			log.Panicf("failed craeting log path %s", path)
		}
	}
}

func authentication() {
	log.Println("Checking authentication status")
	if _, err := os.Stat(middleware.PubKeyPath); os.IsNotExist(err) {
		log.Println("Missing public key, fetching ...")
		middleware.FetchKey()
	}
}

func databaseConnection() {
	log.Println("Opening database connection")
	var err error
	conn, err = database.NewConnection(
		runtimebag.GetEnvString(constants.DatabaseDriver, ""),
		runtimebag.GetEnvString(constants.DatabaseUser, ""),
		runtimebag.GetEnvString(constants.DatabasePass, ""),
		runtimebag.GetEnvString(constants.DatabaseHost, ""),
		runtimebag.GetEnvString(constants.DatabasePort, ""),
		runtimebag.GetEnvString(constants.DatabaseName, ""),
		runtimebag.GetEnvBool(constants.QueryLog, false),
		writers,
	)
	if err != nil {
		log.Panic(err.Error())
	}

	maxIdleConns64 := runtimebag.GetEnvInt(constants.MaxIdleConns, 0)
	if maxIdleConns64 > math.MaxInt32 {
		log.Fatalf("MaxIdleConns value is too large: %d", maxIdleConns64)
	}
	maxIdleConns := int(maxIdleConns64)

	maxOpenConns64 := runtimebag.GetEnvInt(constants.MaxOpenConns, 0)
	if maxOpenConns64 > math.MaxInt32 {
		log.Fatalf("MaxOpenConns value is too large: %d", maxOpenConns64)
	}
	maxOpenConns := int(maxOpenConns64)

	connMaxLifetime64 := runtimebag.GetEnvInt(constants.ConnMaxLifetime, 0)
	if connMaxLifetime64 > math.MaxInt32 {
		log.Fatalf("ConnMaxLifetime value is too large: %d", connMaxLifetime64)
	}
	connMaxLifetime := int(connMaxLifetime64)

	err = conn.ConfigConnectionPooling(maxIdleConns, maxOpenConns, connMaxLifetime)
	if err != nil {
		log.Panic(err.Error())
	}
}

func migrations() {
	log.Println("Running migrations ...")
	if runtimebag.GetEnvBool(constants.ShouldRunMigrations, constants.ShouldRunMigrationsDefault) {
		if err := repositories.Migrate(conn); err != nil {
			log.Panic(err.Error())
		}
	}
}

func seeding() {
	if runtimebag.GetEnvBool(constants.DatabaseSeed, true) {
		log.Println("Seeding database ...")
		tenantApp := application.NewTenant(
			tenant.NewNewTenant(
				repositories.NewTenant(conn),
				amLogger,
			),
			repositories.NewTenant(conn),
			amLogger,
		)

		_, appErr := tenantApp.Create(appContext.NewRequestContext("SEED PROCESS", "SEED PROCESS", nil, nil),
			&application.CreateTenantRequest{
				Name: constants.DefaultTenant,
			})
		if appErr != nil {
			log.Fatal(appErr.ToDto())
		}
	}
}

func buildLogs() {
	createLogsDirectory()
	log.Println("Building logs ...")
	writers = logging.GetLogWriters("logs/agent-management.log", 15, 1024, 7, constants.ServiceName)
	amLogger = zerolog.New(io.MultiWriter(writers.Writers()...)).Level(logLevel())
}

func buildDependencies() {
	log.Println("Building dependencies ...")
	runtimeBag = runtimebag.NewBagWithPreloadedEnvs()
	prometheusRegistry = prometheus.NewRegistry()
	chassisMetrics = metrics.NewMetrics(prometheusRegistry)

	baseController := httpControllers.NewController(repositories.NewTenant(conn))

	authApp := application.NewAuth(
		repositories.NewUser(conn),
		amLogger,
	)
	authController = httpControllers.NewAuth(
		authApp,
		baseController,
	)

	commentApp := application.NewComment(
		repositories.NewComment(conn),
		amLogger,
	)
	commentController = httpControllers.NewComments(
		commentApp,
		baseController,
	)

	statusApp := application.NewStatus(
		repositories.NewStatus(conn),
		amLogger,
	)
	statusController = httpControllers.NewStatus(
		statusApp,
		baseController,
	)

	eventApp := application.NewEvent(
		repositories.NewEvent(conn),
		amLogger,
	)
	eventController = httpControllers.NewEvent(
		eventApp,
		baseController,
	)
}

func httpRouter() *gin.Engine {
	if router != nil {
		return router
	}
	gin.DefaultWriter = io.MultiWriter(os.Stdout)
	router = gin.New()
	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowAllOrigins:  true,
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
	router.Use(cors.New(config))
	router.Use(gin.RecoveryWithWriter(io.MultiWriter(writers.Writers()...)))
	//router.Use(middleware.XCorrelate())

	router.GET("/api/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"live": "true"})
	})

	router.GET("/api/health/ready", func(c *gin.Context) {
		status := health.Validate()
		c.JSON(status.StatusCode(), status)
	})

	router.GET("/metrics", gin.WrapH(promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{Registry: prometheusRegistry})))

	router.Use(middleware.Authenticate())
	router.Use(middleware.Tenant())
	if runtimebag.GetEnvBool(constants.RequestLog, false) {
		router.Use(middleware.Log(amLogger))
	}
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	router.POST("/api/register", authController.SignUpUser)
	router.POST("/api/login", authController.LoginUser)
	router.POST("/api/otp/generate", authController.GenerateOTP)
	router.POST("/api/otp/verify", authController.VerifyOTP)
	router.POST("/api/otp/validate", authController.ValidateOTP)
	router.POST("/api/otp/disable", authController.DisableOTP)

	router.GET("/api/events", eventController.GetAll)
	router.GET("/api/events/:id", eventController.Get)
	router.POST("/api/events", eventController.Create)
	router.PUT("/api/events/:id", eventController.Update)
	router.DELETE("/api/events/:id", eventController.Delete)

	router.GET("/api/status", statusController.GetAll)
	router.GET("/api/status/:id", statusController.Get)
	router.POST("/api/status", statusController.Create)
	router.PUT("/api/status/:id", statusController.Update)
	router.DELETE("/api/status/:id", statusController.Delete)

	router.GET("/api/comments", commentController.GetAll)
	router.GET("/api/comments/:id", commentController.Get)
	router.POST("/api/comments", commentController.Create)
	router.PUT("/api/comments/:id", commentController.Update)
	router.DELETE("/api/comments/:id", commentController.Delete)

	return router
}

func logLevel() zerolog.Level {
	switch strings.ToUpper(runtimebag.GetEnvString(constants.LogLevel, "ERROR")) {
	case "DEBUG":
		return zerolog.DebugLevel
	case "WARN":
		return zerolog.WarnLevel
	case "INFO":
		return zerolog.InfoLevel
	case "ERROR":
		return zerolog.ErrorLevel
	}
	return zerolog.DebugLevel
}
