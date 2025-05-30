package dbkit

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func TOString(l logger.LogLevel) string {
	if l == logger.Error {
		return "error"
	} else if l == logger.Info {
		return "INFO"
	} else if l == logger.Silent {
		return "SILENY"
	} else if l == logger.Warn {
		return "warn"
	} else {
		return fmt.Sprintf("%d", l)
	}
}

// 切换loglevel
func TOGormLogLevel(level string) (r logger.LogLevel) {
	_level := strings.ToLower(level)
	switch _level {
	case "error":
		r = logger.Error
	case "warn":
		r = logger.Warn
	case "debug":
	case "info":
	case "trace":
		r = logger.Info
	case "fatal":
	case "panic":
		r = logger.Error
	default:
		r = logger.Info
	}
	return r
}

// gorm.Dialector
func OpenDb(linkstr string, opts ...Option) (db *gorm.DB, err error) {
	var dialector gorm.Dialector
	dbinfo, err := Parse(linkstr)
	if err != nil {
		return nil, err
	}
	dbtype := dbinfo.DbType
	dsn := linkstr
	switch dbtype {
	case MySQL, Tidb:
		dialector = mysql.Open(dbinfo.StandardDSN)
	case SQLite:
		dialector = sqlite.Open(dsn)
	case SqlServer:
		dialector = sqlserver.Open(dsn)
	case PostgreSQL:
		dialector = postgres.Open(dsn)
	default:

	}
	ctx := NewDbContext()
	for _, opt := range opts {
		opt(ctx)
	}
	newLogger := logger.New(
		log.New(ctx.Writer, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,                   // Slow SQL threshold
			LogLevel:                  ctx.LogLevel,                  // Log level
			IgnoreRecordNotFoundError: ctx.IgnoreRecordNotFoundError, // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      ctx.ParameterizedQueries,      // Don't include params in the SQL log
			Colorful:                  false,                         // Disable color
		},
	)
	//fmt.Println("level%d=%s", ctx.LogLevel, TOString(ctx.LogLevel))
	db, err = gorm.Open(dialector, &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: ctx.SingularTable,
			TablePrefix:   ctx.TablePrefix,
		},
	})

	if err != nil {
		return nil, err
	}
	if ctx.Debug {
		db = db.Debug()
	}
	if len(ctx.ModuleMigrates) > 0 {
		err = db.AutoMigrate(ctx.ModuleMigrates...)
	}
	if err != nil {
		return
	}
	// 设置 连接池
	sqlDb, _ := db.DB()
	sqlDb.SetMaxIdleConns(ctx.MaxIdleConns)
	sqlDb.SetMaxOpenConns(ctx.MaxOpenConns)
	sqlDb.SetConnMaxIdleTime(ctx.ConnMaxIdleTime)
	sqlDb.SetConnMaxLifetime(ctx.ConnMaxLifetime)
	return db, err
}
