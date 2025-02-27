package initialize

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"time"

	global "gatewayhandler/pkg/global"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 数据库配置
type DbConfig struct {
	Host          string
	Port          int
	DbName        string
	Username      string
	Password      string
	TimeZone      string
	LogLevel      int
	SlowThreshold int
	IdleConns     int
	OpenConns     int
}

func PgInit(cfgPath string) (*gorm.DB, error) {
	// 初始化配置
	config, err := LoadDbConfig()
	if err != nil {
		logrus.Errorf("加载数据库配置失败: %v", err)
		return nil, err
	}

	// 初始化数据库
	db, err := PgConnect(config)
	if err != nil {
		logrus.Error("连接数据库失败:", err)
		return nil, err
	}
	global.DB = db
	return db, nil
}

// LoadDbConfig 从配置文件加载数据库配置
func LoadDbConfig() (*DbConfig, error) {
	config := &DbConfig{
		Host:          viper.GetString("db.psql.host"),
		Port:          viper.GetInt("db.psql.port"),
		DbName:        viper.GetString("db.psql.dbname"),
		Username:      viper.GetString("db.psql.username"),
		Password:      viper.GetString("db.psql.password"),
		TimeZone:      viper.GetString("db.psql.time_zone"),
		LogLevel:      viper.GetInt("db.psql.log_level"),
		SlowThreshold: viper.GetInt("db.psql.slow_threshold"),
		IdleConns:     viper.GetInt("db.psql.idle_conns"),
		OpenConns:     viper.GetInt("db.psql.open_conns"),
	}

	// 设置默认值
	if config.Host == "" {
		config.Host = "localhost"
	}
	if config.Port == 0 {
		config.Port = 5432
	}
	if config.TimeZone == "" {
		config.TimeZone = "Asia/Shanghai"
	}
	if config.LogLevel == 0 {
		config.LogLevel = 1
	}
	if config.SlowThreshold == 0 {
		config.SlowThreshold = 200
	}
	if config.IdleConns == 0 {
		config.IdleConns = 10
	}
	if config.OpenConns == 0 {
		config.OpenConns = 50
	}

	// 检查必要的配置
	if config.DbName == "" || config.Username == "" || config.Password == "" {
		return nil, fmt.Errorf("database configuration is incomplete")
	}

	return config, nil
}

// GormLogger 自定义 GORM 日志器
type GormLogger struct {
	LogLevel logger.LogLevel
}

// LogMode 设置日志级别
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = level
	return l
}

// Info 输出 Info 级别日志
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		logrus.WithContext(ctx).Infof(msg, data...)
	}
}

// Warn 输出 Warn 级别日志
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		logrus.WithContext(ctx).Warnf(msg, data...)
	}
}

// Error 输出 Error 级别日志
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		logrus.WithContext(ctx).Errorf(msg, data...)
	}
}

// Trace 输出 Trace 级别日志
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}
	elapsed := time.Since(begin)
	sql, rows := fc()

	// 获取当前调用栈信息（穿透 4 层）
	pc := make([]uintptr, 1)
	runtime.Callers(7, pc)
	frame, _ := runtime.CallersFrames(pc).Next()
	dir := filepath.Dir(frame.File)
	logEntry := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"elapsed": elapsed,
		"rows":    rows,
		"gorm_caller": fmt.Sprintf("%s/%s:%d",
			filepath.Base(dir),
			filepath.Base(frame.File),
			frame.Line), // 直接记录 SQL 触发点
	})
	filepath.Base(frame.File)
	switch {
	case err != nil && l.LogLevel >= logger.Error:
		logEntry.Errorf("SQL: %s, Error: %v", sql, err)
	case l.LogLevel >= logger.Info:
		logEntry.Infof("SQL: %s", sql)
	}
}

// PgInit 初始化数据库连接
func PgConnect(config *DbConfig) (*gorm.DB, error) {
	dataSource := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable TimeZone=%s",
		config.Host, config.Port, config.DbName, config.Username, config.Password, config.TimeZone)

	var err error
	db, err := gorm.Open(postgres.Open(dataSource), &gorm.Config{
		Logger:                 &GormLogger{LogLevel: logger.Info},
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: false, // use singular table name, table for `User` would be `user` with this option enabled
		},
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取原生数据库连接失败: %v", err)
	}

	sqlDB.SetMaxIdleConns(config.IdleConns)
	sqlDB.SetMaxOpenConns(config.OpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logrus.Println("连接数据库完成...")

	return db, nil
}
