package api

import (
	"github.com/go-redis/redis/v8"
	"github.com/mohammadrabetian/quick/handlers"
	"github.com/mohammadrabetian/quick/pkg/mysql"
	"github.com/mohammadrabetian/quick/repository"
	"github.com/mohammadrabetian/quick/service"
	"github.com/mohammadrabetian/quick/util"
	"gorm.io/gorm"
)

// All stores e.g. nosql,sql
type Store struct {
	SQL       *gorm.DB
	Cache     *redis.Client
	walletSvc service.WalletService
	userSvc   service.UserService
}

func NewStore(config util.Config) *Store {
	db := mysql.NewDatabase(config)
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host,
		Password: config.Redis.Password,
		DB:       config.Redis.DBName,
	})

	// initialize the repo
	walletRepo := repository.NewWalletMySQLRepository(db.DB, rdb)
	userRepo := repository.NewUserMySQLRepository(db.DB)

	// initialize the service and handlers
	walletSvc := service.NewWalletService(walletRepo)
	userSvc := service.NewUserService(userRepo)
	handlers.InitWalletHandlers(walletSvc)
	handlers.InitUserHandlers(userSvc)

	return &Store{
		SQL:       db.DB,
		Cache:     rdb,
		walletSvc: *walletSvc,
		userSvc:   *userSvc,
	}
}
