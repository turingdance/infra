package dbkit

import (
	"io"
	"os"
	"time"

	ctllogger "github.com/turingdance/infra/logger"
	"gorm.io/gorm/logger"
)

type Option func(*DbContext)

func AutoMigrate(obj ...interface{}) Option {
	return func(ctx *DbContext) {
		if ctx.ModuleMigrates == nil {
			ctx.ModuleMigrates = make([]interface{}, 0)
		}
		ctx.ModuleMigrates = append(ctx.ModuleMigrates, obj...)
	}
}

// SetLogLevel
func SetLogLevel(level ctllogger.LogLevel) Option {
	return func(ctx *DbContext) {
		//fmt.Println("SetLogLevel(level int32),%s ", level, TOString(logger.LogLevel(level)))
		ctx.LogLevel = TOGormLogLevel(level)
	}
}

// 设置最大可用连接数
func SetMaxOpenConns(num int) Option {
	return func(dc *DbContext) {
		dc.MaxOpenConns = num
	}
}

// 设置最大空闲连接数
func SetMaxIdleConns(num int) Option {
	return func(dc *DbContext) {
		dc.MaxIdleConns = num
	}
}

// 设置最大保持时间
func SetConnMaxLifetime(d time.Duration) Option {
	return func(dc *DbContext) {
		dc.ConnMaxLifetime = d
	}
}

// 设置最大空闲时间
func SetConnMaxIdleTime(d time.Duration) Option {
	return func(dc *DbContext) {
		dc.ConnMaxIdleTime = d
	}
}

// SetLogLevel
func WithPrefix(prefix string) Option {
	return func(ctx *DbContext) {
		ctx.TablePrefix = prefix
	}
}

// SetLogLevel
func IgnoreRecordNotFoundError(flag bool) Option {
	return func(ctx *DbContext) {
		ctx.IgnoreRecordNotFoundError = flag
	}
}

// SetLogLevel
func SingularTable(flag bool) Option {
	return func(ctx *DbContext) {
		ctx.SingularTable = flag
	}
}

// SetLogLevel
func WithWriter(writers ...io.Writer) Option {
	return func(ctx *DbContext) {
		ctx.Writer = io.MultiWriter(writers...)
	}
}

// SetLogLevel
func ParameterizedQueries(f bool) Option {
	return func(ctx *DbContext) {
		ctx.ParameterizedQueries = f
	}
}

func Debug(f bool) Option {
	return func(ctx *DbContext) {
		ctx.Debug = f
	}
}

type DbContext struct {
	LogLevel                  logger.LogLevel
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
	Writer                    io.Writer
	TablePrefix               string
	SingularTable             bool
	ModuleMigrates            []interface{}
	Debug                     bool
	MaxIdleConns              int           // 最大空闲连接数
	MaxOpenConns              int           // 最大打开连接数
	ConnMaxLifetime           time.Duration // 连接最大生命周期
	ConnMaxIdleTime           time.Duration // 连接最大空闲时间
}

func NewDbContext() *DbContext {
	return &DbContext{
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      true,
		Writer:                    os.Stdout,
		TablePrefix:               "",
		SingularTable:             true,
		ModuleMigrates:            make([]interface{}, 0),
		Debug:                     false,
		MaxIdleConns:              2,
		MaxOpenConns:              10,
		ConnMaxLifetime:           time.Hour * 1,
		ConnMaxIdleTime:           time.Second * 10,
	}
}
