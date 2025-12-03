package storage

import "strings"

type DriverType string

const (
	DriverLocal DriverType = "local"
	DriverOSS   DriverType = "oss"
)

type StorageConf struct {
	Driver       DriverType `json:"driver"`    //驱动
	Bucket       string     `json:"bucket"`    // 使用的bucket  mnt
	Datapath     string     `json:"-"`         // 映射路径
	Ssl          bool       `json:"ssl"`       //是否使用https
	Host         string     `json:"host"`      //对外域名 ios.turingdance.com
	FunSubDir    []string   `json:"funSubDir"` //yyyy,mm,dd,yymmdd,userId,teamId,
	FunName      string     `json:"funName"`   //filename,uuid,timestamp 文件名称生成策略
	AuthRequired bool       `json:"authRequired"`
	Primary      bool       `json:"primary"`
}

func (c *StorageConf) FileKeyAndPath(ctx Context) (key string, path string) {
	subdirs := Invoke(ctx, c.FunSubDir...)
	filekey := Invoke(ctx, c.FunName)
	if c.Driver == DriverLocal {
		key = strings.Join([]string{c.Bucket, subdirs, filekey + ctx.Ext}, "/")
		path = strings.Join([]string{c.Datapath, subdirs, filekey + ctx.Ext}, "/")
	} else {
		key = strings.Join([]string{subdirs, filekey + ctx.Ext}, "/")
		path = strings.Join([]string{subdirs, filekey + ctx.Ext}, "/")
	}
	return
}
