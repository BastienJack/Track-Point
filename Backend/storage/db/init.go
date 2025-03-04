package db

import (
	"commerce/pkg/viper"
	"commerce/pkg/zap"

	"gorm.io/driver/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"fmt"
	"time"
)

var (
	_db       *gorm.DB
	config    = viper.Init("db")
	zapLogger = zap.InitLogger()
)

func getDsn(driverWithRole string) string {
	// user info
	username := config.Viper.GetString(fmt.Sprintf("%s.username", driverWithRole))
	password := config.Viper.GetString(fmt.Sprintf("%s.password", driverWithRole))

	// TCP info
	host := config.Viper.GetString(fmt.Sprintf("%s.host", driverWithRole))
	port := config.Viper.GetInt(fmt.Sprintf("%s.port", driverWithRole))

	// database
	database := config.Viper.GetString(fmt.Sprintf("%s.database", driverWithRole))

	// charset
	charset := config.Viper.GetString(fmt.Sprintf("%s.charset", driverWithRole))

	// data source name
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		username, password, host, port, database, charset)

	return dsn
}

func init() {
	// source database config
	dsn1 := getDsn("mysql.source")

	var err error
	_db, err = gorm.Open(mysql.Open(dsn1), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic(err.Error())
	}

	// replicas database config
	// dsn2 := getDsn("mysql.replica1")

	// set database resolver
	_db.Use(dbresolver.Register(dbresolver.Config{
		Sources: []gorm.Dialector{mysql.Open(dsn1)},
		// Replicas:          []gorm.Dialector{mysql.Open(dsn2)},
		Policy:            dbresolver.RandomPolicy{},
		TraceResolverMode: false,
	}))

	// other tables will be automatically changed
	if err := _db.AutoMigrate(&User{}, &Event{}); err != nil {
		zapLogger.Fatal(err.Error())
	}

	// open database
	db, err := _db.DB()
	if err != nil {
		zapLogger.Fatal(err.Error())
	}

	// set database options
	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)

	zapLogger.Info("Mysql server connection successful!")
}

func GetDB() *gorm.DB {
	return _db
}
