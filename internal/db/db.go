package db

import (
	"sync"

	"github.com/mohdjishin/SplitWise/config"
	"github.com/mohdjishin/SplitWise/internal/models"
	log "github.com/mohdjishin/SplitWise/logger"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBManagerInterface interface {
	Connect()
	GetDB() *gorm.DB
}

type DBManager struct {
	db *gorm.DB
}

var (
	dbManagerInstance DBManagerInterface
	once              sync.Once
)

func SetDbManager(manager DBManagerInterface) {
	dbManagerInstance = manager
}

func GetDbManagerInstance() DBManagerInterface {
	once.Do(func() {
		if dbManagerInstance == nil {
			dbManagerInstance = &DBManager{}
			dbManagerInstance.Connect()
		}
	})
	return dbManagerInstance
}

func init() {
	_ = GetDbManagerInstance()
}
func (m *DBManager) Connect() {
	log.Info("Connecting to database")
	var err error
	m.db, err = gorm.Open(postgres.Open(config.GetConfig().DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database", zapcore.Field{Key: "error", Type: zapcore.ErrorType, Interface: err})
	}
	log.Info("Connected to database")
	log.Info("Migrating database")
	err = m.db.AutoMigrate(&models.User{}, &models.Group{}, models.BillHistory{}, &models.Bill{}, &models.GroupMember{})
	if err != nil {
		log.Fatal("failed to migrate database", zapcore.Field{Key: "error", Type: zapcore.ErrorType, Interface: err})
	}
	log.Info("Database migration successful")
}

func (m *DBManager) GetDB() *gorm.DB {
	return m.db
}

func GetDb() *gorm.DB {
	return GetDbManagerInstance().GetDB()
}
