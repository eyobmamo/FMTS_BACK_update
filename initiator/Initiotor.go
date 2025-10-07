package initiator

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// "FMTS/internal/adapter/inbound/http/responseutil"
	// "FMTS/pkg/utils"

	"FMTS/utils"
	util "FMTS/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Initiator() {
	logger := util.NewLogger()

	logger.Infof("Initializing MongoDB client...")
	mongoClient := InitMongo(logger)
	logger.Infof("MongoDB client initialized")

	// logger.Infof("Initializing Minio client...")
	// minioClient := InitMinio(cfg.MinioEndPoint, cfg.MinioAccessKey, cfg.MinioSecretKey, logger)
	// logger.Infof("Minio client initialized")

	logger.Infof("Initializing persistence...")
	persistence := InitPersistence(mongoClient, "FMTS", logger)
	logger.Infof("Persistence initialized")

	logger.Infof("Initializing JWT manager...")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	key := os.Getenv("KEY")
	iv := os.Getenv("IV")

	if jwtSecretKey == "" {
		logger.Fatalf("JWT_SECRET_KEY environment variable is not set")
	}
	if key == "" {
		logger.Fatalf("KEY environment variable is not set")
	}
	if iv == "" {
		logger.Fatalf("IV environment variable is not set")
	}

	jwtManager := utils.NewJWTManager(jwtSecretKey, 15*time.Minute, 7*24*time.Hour, key, iv) // 15 min access, 7 days refresh
	logger.Infof("JWT manager initialized")

	logger.Infof("Initializing domain services...")
	domain := InitDomain(persistence, logger, jwtManager)
	logger.Infof("Domain services initialized")

	logger.Infof("Initializing application services...")
	application := InitApplication(domain, logger)
	logger.Infof("Application services initialized")

	logger.Infof("Initializing adapter services...")
	adapter := InitAdapter(application, logger)
	logger.Infof("Adapter services initialized")

	logger.Infof("Initializing Chi router...")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))
	r.NotFound(NotFoundHandler)
	logger.Infof("Chi router initialized")

	logger.Infof("Initializing routes...")
	InitRoutes(r, adapter, jwtSecretKey, key, iv, logger)
	logger.Infof("Routes initialized")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // for local development
	}
	println(port)
	server := http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Infof("🚀 Server started")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Server stopped with error: %v", err)
		}
	}()

	sig := <-quit
	logger.Infof("Server shutting down with signal: %v", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Failed to shutdown gracefully: %v", err)
	}

	logger.Infof("Server shutdown successfully")
}

// NotFoundHandler returns a 404 response in JSON format.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error":"resource not found"}`))
}
