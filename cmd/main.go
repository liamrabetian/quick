//	@securitydefinitions.apikey	ApiKeyAuth
//
//	@in							header
//	@name						Authorization
package main

import (
	"fmt"
	"os"

	"github.com/mohammadrabetian/quick/api"
	"github.com/mohammadrabetian/quick/domain"
	"github.com/mohammadrabetian/quick/util"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		logrus.WithError(err).
			Fatal("cannot load config")
	}

	// Set up Logger
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// More readable logs for development env
	if config.Environment == "development" {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			DisableColors: false,
			FullTimestamp: true,
		})
		logrus.SetOutput(os.Stderr)
	}

	runGinServer(config)
}

func runGinServer(config util.Config) {
	server := api.NewServer(config)

	// Auto-migrate the database schema
	migrateDatabase(server.Store.SQL)
	// Seed the database
	createdUsers := seedUsersDatabase(server.Store.SQL)
	seedWalletsDatabase(server.Store.SQL, createdUsers)

	err := server.Start(config.HTTPServer.Address)
	if err != nil {
		logrus.WithError(err).
			Fatal("cannot run server")
	}
}

func seedWalletsDatabase(db *gorm.DB, users []domain.User) {
	var walletCount int64
	db.Model(&domain.Wallet{}).Count(&walletCount)

	if walletCount == 0 {
		initialWallets := []domain.Wallet{
			{ID: 1, Balance: decimal.NewFromInt(100), UserID: users[0].Username},
			{ID: 2, Balance: decimal.NewFromInt(200), UserID: users[1].Username},
			{ID: 3, Balance: decimal.NewFromInt(300), UserID: users[2].Username},
		}

		for _, wallet := range initialWallets {
			err := db.Create(&wallet).Error
			if err != nil {
				logrus.Fatalf("Failed to seed database: %v", err)
			}
		}

		logrus.Info("Wallet Database seeding completed")
	}
}

func seedUsersDatabase(db *gorm.DB) []domain.User {
	var userCount int64
	db.Model(&domain.User{}).Count(&userCount)
	var createdUsers []domain.User
	if userCount == 0 {
		initialUsers := []domain.User{
			{ID: 1, Username: "user1", Password: "password1"},
			{ID: 2, Username: "user2", Password: "password2"},
			{ID: 3, Username: "user3", Password: "password3"},
		}

		for _, user := range initialUsers {
			token := fmt.Sprintf("token_%d", user.ID)
			user.Token = token
			err := db.Create(&user).Error
			if err != nil {
				logrus.Fatalf("Failed to seed database: %v", err)
			}
			createdUsers = append(createdUsers, user)
		}

		logrus.Info("User database seeding completed")
	}
	return createdUsers
}

func migrateDatabase(db *gorm.DB) {
	err := db.AutoMigrate(&domain.Wallet{}, &domain.User{})
	if err != nil {
		logrus.Fatalf("Failed to migrate database: %v", err)
	}
	logrus.Info("Database migration completed")
}
