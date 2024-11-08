package dbkit

import (
	"net/url"
	"regexp"
	"strings"
)

type DbConfig struct {
	UserName string
	Password string
	Protocal string
	Addrss   string
	Dbname   string
	Params   url.Values
}

// /dbname
// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
func ParseMysql(dsn string) *DbConfig {
	cfg := &DbConfig{}
	re := regexp.MustCompile(`([^:^/]*)\:?([^/]*)@([a-zA-Z]*)\(?([^\(^\)]*)\)?\/([^\?]*)\??(.*)`)
	strs := re.FindStringSubmatch(dsn)
	cfg.UserName = strs[1]
	cfg.Password = strs[2]
	cfg.Protocal = strs[3]
	cfg.Addrss = strs[4]
	cfg.Dbname = strs[5]
	cfg.Params, _ = url.ParseQuery(strs[6])
	return cfg
}

func guessdbtype(dsn string) DBTYPE {
	dbtype := SQLite
	if strings.Contains(dsn, ".db") {
		dbtype = SQLite
	} else if strings.Contains(dsn, ":") {
		dbtype = MySQL
	} else if strings.Contains(dsn, "sql") {
		dbtype = SqlServer
	} else if strings.Contains(dsn, " ") {
		dbtype = PostgreSQL
	}
	return dbtype
}
