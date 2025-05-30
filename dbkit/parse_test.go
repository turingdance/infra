package dbkit

import "testing"

type DATA struct {
	Dsn    string
	DBTYPE DBTYPE
}

func TestDSN(t *testing.T) {

	datas := []DATA{
		{
			Dsn:    "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
			DBTYPE: MySQL,
		},
		{
			Dsn:    "mysql://user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
			DBTYPE: MySQL,
		},
		{
			Dsn:    "sqlserver://user:pass@localhost:1433?database=dbname",
			DBTYPE: SqlServer,
		},
		{
			Dsn:    "/mnt/data/test.db",
			DBTYPE: SQLite,
		},
		{
			Dsn:    "postgres://user:pass@localhost:5432/dbname?sslmode=disable",
			DBTYPE: PostgreSQL,
		},
	}
	for _, v := range datas {
		info, err := Parse(v.Dsn)
		if err != nil {
			t.Errorf("%s,%s", v.Dsn, err.Error())
		}
		if info.DbType != v.DBTYPE {
			t.Errorf("%s,expected %s while get %s", v.Dsn, v.DBTYPE, info.DbType)
		}
	}

}
