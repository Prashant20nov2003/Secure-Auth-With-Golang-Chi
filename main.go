package main

import (
	. "betamart/controller"
	"betamart/internal/database"
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Port is not found in environment")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not found in environment")
	}

	connection, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}

	queries := database.New(connection)

	apiCfg := ApiConfig{
		Query: queries,
		DB:    connection,
	}

	router := chi.NewRouter()

	corsOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   corsOrigins,
		AllowedMethods:   []string{"POST", "GET", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	apiRouter := chi.NewRouter()

	apiRouter.Use(httprate.LimitByIP(120, 1*time.Minute))

	// User Auth
	apiRouter.Post("/register", apiCfg.RegisterUser)
	apiRouter.With(httprate.LimitByIP(5, 1*time.Minute)).Post("/login", apiCfg.LoginUser)
	apiRouter.Get("/getusername", apiCfg.UserMiddleware(apiCfg.GetUsername))

	// Email Verificaiton
	apiRouter.With(httprate.LimitByIP(5, 1*time.Minute)).Post("/generate_email_verification", apiCfg.UserMiddleware(apiCfg.GenerateEmailVerification))
	apiRouter.Post("/resend_email_code/{id}", apiCfg.ResendEmailVerificationCode)
	apiRouter.Get("/email_verification/{id}", apiCfg.FetchEmailVerification)
	apiRouter.With(httprate.LimitByIP(10, 1*time.Minute)).Post("/email_verification/{id}", apiCfg.VerifyEmailVerification)

	// Update Security
	apiRouter.With(httprate.LimitByIP(5, 1*time.Minute)).Post("/forgot_password/{username}", apiCfg.ForgotPassword(apiCfg.GenerateEmailVerification))
	apiRouter.Post("/forgot_password/{username}/{id}", apiCfg.ForgotPassword(apiCfg.ChangePassword))
	apiRouter.With(httprate.LimitByIP(5, 1*time.Minute)).Post("/change_password", apiCfg.UserMiddleware(apiCfg.GenerateEmailVerification))
	apiRouter.Post("/change_password/{id}", apiCfg.UserMiddleware(apiCfg.ChangePassword))

	// Product
	apiRouter.Post("/product", apiCfg.UserMiddleware(apiCfg.PostProduct))
	apiRouter.Get("/product", apiCfg.UserMiddleware(apiCfg.GetProduct))
	apiRouter.Get("/product?isPrivate=true", apiCfg.UserMiddleware(apiCfg.GetProduct))

	router.Mount("/api", apiRouter)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	go StorageProxy()

	log.Printf("Server is running: http://localhost:%s", portString)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Couldn't connect to server:", err)
	}
}
