package dbkit

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type DbConfig struct {
	DbType      DBTYPE
	UserName    string
	Password    string
	Protocal    string
	Addrss      string
	Host        string // 主机名
	Port        string // 端口
	Dbname      string // 数据库名
	Params      url.Values
	FilePath    string // SQLite文件路径
	OriginalDSN string // 原始DSN字符串
	StandardDSN string
}

// Parse 解析各种数据库类型的 DSN 字符串
func Parse(dsn string) (*DbConfig, error) {
	dsn = strings.TrimSpace(dsn)
	if dsn == "" {
		return nil, errors.New("DSN 字符串不能为空")
	}
	// 检测数据库类型并调用相应的解析器
	switch {
	case isMySQLDSN(dsn):
		return ParseMysql(dsn)
	case isPostgreSQLDSN(dsn):
		return ParsePostgreSQL(dsn)
	case isSQLiteDSN(dsn):
		return ParseSQLite(dsn)
	case isSQLServerDSN(dsn):
		return ParseSQLServer(dsn)
	case isTiDBDSN(dsn):
		return ParseTiDB(dsn)
	case isClickHouseDSN(dsn):
		return ParseClickHouse(dsn)
	default:
		return nil, fmt.Errorf("无法识别的 DSN 格式: %s", dsn)
	}
}

// isMySQLDSN 检测是否为 MySQL DSN
func isMySQLDSN(dsn string) bool {
	// 检查 URL 格式 mysql://user:pass@host/dbname
	if strings.HasPrefix(dsn, "mysql://") {
		return true
	}
	// 检查标准格式 user:pass@protocol(addr)/dbname
	if strings.Contains(dsn, "@") && strings.Contains(dsn, ")") {
		return true
	}

	return false
}

// isPostgreSQLDSN 检测是否为 PostgreSQL DSN
func isPostgreSQLDSN(dsn string) bool {
	// 检查 URL 格式 postgres://user:pass@host/dbname
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		return true
	}
	// 检查关键字/值格式 host=... user=...
	if strings.Contains(dsn, "host=") && strings.Contains(dsn, "user=") {
		return true
	}
	return false
}

// isSQLiteDSN 检测是否为 SQLite DSN
func isSQLiteDSN(dsn string) bool {
	// 检查 URL 格式 sqlite:///path/to/db
	if strings.HasPrefix(dsn, "sqlite://") {
		return true
	}
	// 检查文件路径格式 (假设不以关键字/值对格式出现)
	if !strings.Contains(dsn, "=") && (strings.Contains(dsn, ".db") || strings.Contains(dsn, "/")) {
		return true
	}
	return false
}

// isSQLServerDSN 检测是否为 SQL Server DSN
func isSQLServerDSN(dsn string) bool {
	// 检查 URL 格式 sqlserver://user:pass@host/dbname
	if strings.HasPrefix(dsn, "sqlserver://") {
		return true
	}
	// 检查连接字符串格式 Server=...;Database=...
	if strings.Contains(dsn, "Server=") && strings.Contains(dsn, "Database=") {
		return true
	}
	return false
}

// isTiDBDSN 检测是否为 TiDB DSN
func isTiDBDSN(dsn string) bool {
	// TiDB 兼容 MySQL 格式，检查是否显式指定 tidb:// 协议
	if strings.HasPrefix(dsn, "tidb://") {
		return true
	}
	// 否则按 MySQL 处理 (TiDB 兼容 MySQL 协议)
	return isMySQLDSN(dsn)
}

// isClickHouseDSN 检测是否为 ClickHouse DSN
func isClickHouseDSN(dsn string) bool {
	// 检查 URL 格式 clickhouse://user:pass@host/dbname
	return strings.HasPrefix(dsn, "clickhouse://")

}

// ParsePostgreSQL 解析 PostgreSQL DSN
func ParsePostgreSQL(dsn string) (*DbConfig, error) {
	config := &DbConfig{
		OriginalDSN: dsn,
		DbType:      PostgreSQL,
		Params:      make(url.Values),
	}

	// 处理 URL 格式
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		u, err := url.Parse(dsn)
		if err != nil {
			return nil, fmt.Errorf("解析 PostgreSQL URL 失败: %v", err)
		}

		config.Protocal = u.Scheme
		if u.User != nil {
			config.UserName = u.User.Username()
			config.Password, _ = u.User.Password()
		}

		host, port := splitHostPort(u.Host)
		config.Host = host
		config.Port = port
		config.Addrss = u.Host

		config.Dbname = strings.TrimPrefix(u.Path, "/")
		config.Params = u.Query()

		return config, nil
	}

	// 处理关键字/值格式
	parts := strings.Fields(dsn)
	for _, part := range parts {
		if !strings.Contains(part, "=") {
			continue
		}

		key, value, _ := strings.Cut(part, "=")
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		switch strings.ToLower(key) {
		case "host":
			config.Host = value
			config.Addrss = value
		case "port":
			config.Port = value
			if config.Addrss != "" {
				config.Addrss = fmt.Sprintf("%s:%s", config.Addrss, value)
			}
		case "user", "username":
			config.UserName = value
		case "password":
			config.Password = value
		case "dbname":
			config.Dbname = value
		case "protocol":
			config.Protocal = value
		default:
			config.Params.Add(key, value)
		}
	}

	return config, nil
}

// ParseSQLite 解析 SQLite DSN
func ParseSQLite(dsn string) (*DbConfig, error) {
	config := &DbConfig{
		OriginalDSN: dsn,
		DbType:      SQLite,
		Params:      make(url.Values),
	}

	// 处理 URL 格式
	if strings.HasPrefix(dsn, "sqlite://") {
		u, err := url.Parse(dsn)
		if err != nil {
			return nil, fmt.Errorf("解析 SQLite URL 失败: %v", err)
		}

		config.Protocal = u.Scheme
		config.FilePath = strings.TrimPrefix(u.Path, "/")
		config.Params = u.Query()

		return config, nil
	}

	// 处理文件路径格式
	if strings.Contains(dsn, "?") {
		path, paramsStr, _ := strings.Cut(dsn, "?")
		config.FilePath = path

		params, err := url.ParseQuery(paramsStr)
		if err != nil {
			return nil, fmt.Errorf("解析 SQLite 参数失败: %v", err)
		}
		config.Params = params
	} else {
		config.FilePath = dsn
	}

	return config, nil
}

// ParseSQLServer 解析 SQL Server DSN
func ParseSQLServer(dsn string) (*DbConfig, error) {
	config := &DbConfig{
		OriginalDSN: dsn,
		DbType:      SqlServer,
		Params:      make(url.Values),
	}

	// 处理 URL 格式
	if strings.HasPrefix(dsn, "sqlserver://") {
		u, err := url.Parse(dsn)
		if err != nil {
			return nil, fmt.Errorf("解析 SQL Server URL 失败: %v", err)
		}

		config.Protocal = u.Scheme
		if u.User != nil {
			config.UserName = u.User.Username()
			config.Password, _ = u.User.Password()
		}

		host, port := splitHostPort(u.Host)
		config.Host = host
		config.Port = port
		config.Addrss = u.Host

		config.Dbname = strings.TrimPrefix(u.Path, "/")
		config.Params = u.Query()

		return config, nil
	}

	// 处理连接字符串格式
	parts := strings.Split(dsn, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if !strings.Contains(part, "=") {
			continue
		}

		key, value, _ := strings.Cut(part, "=")
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		switch strings.ToLower(key) {
		case "server":
			host, port := splitHostPort(value)
			config.Host = host
			config.Port = port
			config.Addrss = value
		case "user id", "user", "username":
			config.UserName = value
		case "password", "pwd":
			config.Password = value
		case "database", "initial catalog":
			config.Dbname = value
		default:
			config.Params.Add(key, value)
		}
	}

	return config, nil
}

// ParseMysql 解析 MySQL 格式的 DSN 字符串
func ParseMysql(dsn string) (*DbConfig, error) {
	config := &DbConfig{
		OriginalDSN: dsn,
		DbType:      MySQL,
		Params:      make(url.Values),
		StandardDSN: dsn,
	}

	// 处理空 DSN
	if strings.TrimSpace(dsn) == "" {
		return nil, errors.New("MySQL DSN 不能为空")
	}

	// 检查是否包含协议前缀 (如 mysql://)
	if strings.Contains(dsn, "://") {
		//mysql：// [标准协议]
		if strings.Contains(dsn, ")") {
			arr := strings.Split(dsn, "://")
			config.StandardDSN = arr[1]
			return parseMysqlStandard(config.StandardDSN, config)
		} else {
			return parseMysqlURL(dsn, config)
		}
	}

	// 标准格式: user:password@protocol(address)/dbname?params
	// 支持自定義的
	return parseMysqlStandard(dsn, config)
}

// parseMysqlURL 解析带协议前缀的 MySQL DSN (如 mysql://user:pass@host:port/dbname?param=value)
func parseMysqlURL(dsn string, config *DbConfig) (*DbConfig, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("解析 MySQL URL 失败: %v", err)
	}

	// 验证协议
	if u.Scheme != "mysql" {
		return nil, fmt.Errorf("无效的 MySQL 协议: %s", u.Scheme)
	}

	// 解析用户信息
	if u.User != nil {
		config.UserName = u.User.Username()
		config.Password, _ = u.User.Password()
	}

	// 解析主机和端口
	host, port := splitHostPort(u.Host)
	config.Host = host
	config.Port = port
	config.Addrss = u.Host // 兼容原始Addrss字段

	// 解析数据库名
	config.Dbname = strings.TrimPrefix(u.Path, "/")

	// 解析参数
	config.Params = u.Query()

	return config, nil
}

// parseMysqlStandard 解析标准格式的 MySQL DSN (如 user:password@protocol(address)/dbname?params)
func parseMysqlStandard(dsn string, config *DbConfig) (*DbConfig, error) {
	// 分割用户信息和地址部分
	userInfoPart, hostPart, found := strings.Cut(dsn, "@")
	if !found {
		return nil, errors.New("无效的 MySQL DSN 格式: 缺少 @ 符号")
	}

	// 解析用户名和密码
	user, pass, found := strings.Cut(userInfoPart, ":")
	if found {
		config.UserName = user
		config.Password = pass
	} else {
		config.UserName = userInfoPart
	}

	// 分割地址部分和数据库名
	hostPart, dbPart, found := strings.Cut(hostPart, "/")
	if !found {
		return nil, errors.New("无效的 MySQL DSN 格式: 缺少 / 符号")
	}

	// 解析协议和地址
	protocol, address, found := strings.Cut(hostPart, "(")
	if found {
		config.Protocal = protocol
		// 移除右括号
		if strings.HasSuffix(address, ")") {
			address = strings.TrimSuffix(address, ")")
		}
		config.Addrss = address

		// 解析主机和端口
		host, port := splitHostPort(address)
		config.Host = host
		config.Port = port
	} else {
		// 如果没有协议，默认为 tcp
		config.Protocal = "tcp"
		config.Addrss = hostPart

		// 解析主机和端口
		host, port := splitHostPort(hostPart)
		config.Host = host
		config.Port = port
	}

	// 解析数据库名和参数
	dbName, paramsStr, found := strings.Cut(dbPart, "?")
	if found {
		config.Dbname = dbName
		params, err := url.ParseQuery(paramsStr)
		if err != nil {
			return nil, fmt.Errorf("解析 MySQL 参数失败: %v", err)
		}
		config.Params = params
	} else {
		config.Dbname = dbPart
	}

	return config, nil
}

// splitHostPort 分割主机和端口
func splitHostPort(hostPort string) (host, port string) {
	// 处理 IPv6 地址
	if strings.HasPrefix(hostPort, "[") && strings.Contains(hostPort, "]") {
		endIndex := strings.Index(hostPort, "]")
		if endIndex > 0 && endIndex < len(hostPort)-1 && hostPort[endIndex+1] == ':' {
			return hostPort[1:endIndex], hostPort[endIndex+2:]
		}
		return hostPort[1:endIndex], ""
	}

	// 处理普通主机:端口格式
	host, port, found := strings.Cut(hostPort, ":")
	if found {
		return host, port
	}
	return hostPort, ""
}

// ParseTiDB 解析 TiDB DSN (兼容 MySQL)
func ParseTiDB(dsn string) (info *DbConfig, err error) {
	// TiDB 兼容 MySQL 协议，直接使用 MySQL 解析器
	info, err = ParseMysql(dsn)
	info.DbType = Tidb
	return
}

// ParseClickHouse 解析 ClickHouse DSN
func ParseClickHouse(dsn string) (*DbConfig, error) {
	config := &DbConfig{
		OriginalDSN: dsn,
		DbType:      ClickHouse,
		Params:      make(url.Values),
	}

	// 处理 URL 格式
	if strings.HasPrefix(dsn, "clickhouse://") {
		u, err := url.Parse(dsn)
		if err != nil {
			return nil, fmt.Errorf("解析 ClickHouse URL 失败: %v", err)
		}

		config.Protocal = u.Scheme
		if u.User != nil {
			config.UserName = u.User.Username()
			config.Password, _ = u.User.Password()
		}

		host, port := splitHostPort(u.Host)
		config.Host = host
		config.Port = port
		config.Addrss = u.Host

		config.Dbname = strings.TrimPrefix(u.Path, "/")
		config.Params = u.Query()

		return config, nil
	}

	// 其他格式暂不支持
	return nil, fmt.Errorf("不支持的 ClickHouse DSN 格式: %s", dsn)
}
