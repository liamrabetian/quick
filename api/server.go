package api

import (
	"github.com/mohammadrabetian/quick/docs"
	"github.com/mohammadrabetian/quick/handlers"
	"github.com/mohammadrabetian/quick/middleware"
	"github.com/mohammadrabetian/quick/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	Store  *Store
	config util.Config
	router *gin.Engine
}

// creates an HTTP server
func NewServer(config util.Config) *Server {
	store := NewStore(config)
	server := &Server{config: config, Store: store}
	server.setupRouter()
	return server

}

func (s *Server) setupRouter() {
	router := gin.Default()

	// Auth endpoints
	versionOne := router.Group("v1/auth")
	versionOne.POST("/login", handlers.Login)

	// Register the request logger middleware
	router.Use(s.SetupRequestLogger())

	walletGroup := router.Group("/api/v1/wallets")

	// Apply authentication middleware to walletGroup routes
	walletGroup.Use(middleware.Auth(&s.Store.userSvc))

	err := router.SetTrustedProxies([]string{"192.168.1.2"})
	if err != nil {
		logrus.Fatalf("failed to set trusted proxies")
	}

	// server controllers
	router.GET("/ping", s.Ping)

	// set swagger info
	docs.SwaggerInfo.Title = "Quick Swagger API"
	docs.SwaggerInfo.Description = "Interact with the APIs here"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = s.config.HTTPServer.Address
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// business logic controllers
	{
		walletGroup.GET("/:wallet_id/balance", handlers.GetBalance)
		walletGroup.POST("/:wallet_id/credit", handlers.CreditWallet)
		walletGroup.POST("/:wallet_id/debit", handlers.DebitWallet)
	}

	s.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func (s *Server) SetupRequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logrus.WithFields(logrus.Fields{
			"request_method": c.Request.Method,
			"request_path":   c.Request.URL.Path,
		})

		logger.Info("Started handling request")
		c.Next()
		logger.Info("Completed handling request")
	}
}
