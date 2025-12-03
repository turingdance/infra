package storage

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/turingdance/infra/slicekit"
)

// 这一层里面实现router
type LocalStorage struct {
	config    []StorageConf
	fileKey   string
	bucketKey string
}

// 处理/upload 逻辑
// /attach/render?key=filekey
func NewLocalStorage(config []StorageConf) *LocalStorage {
	return &LocalStorage{
		config:    config,
		fileKey:   "file",
		bucketKey: "bucket",
	}
}

// 处理/upload 逻辑
func (ls *LocalStorage) UploadV1(r *http.Request) (res Response, err error) {
	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile(ls.fileKey)
	if err != nil {
		return
	}
	bucket := r.FormValue(ls.bucketKey)
	config, ok := slicekit.Find(ls.config, func(item StorageConf, index int, slice []StorageConf) bool {
		return item.Bucket == bucket
	})
	if !ok {
		config, ok = slicekit.Find(ls.config, func(item StorageConf, index int, slice []StorageConf) bool {
			return item.Primary
		})
		if !ok {
			err = errors.New("暂未找到配置配件")
		}
	}
	sufix := filepath.Ext(header.Filename)
	ctx := Context{
		File:       file,
		FileHeader: header,
		Ext:        sufix,
		TeamId:     "",
		UserId:     "",
	}
	filekey, _filepath := config.FileKeyAndPath(ctx)
	dir := filepath.Dir(_filepath)
	if _, e := os.Stat(dir); e != nil {
		if os.IsNotExist(e) {
			os.MkdirAll(dir, 0755)
		} else {
			err = e
			return
		}
	}
	f, err := os.OpenFile(_filepath, os.O_WRONLY|os.O_CREATE, 0666) // 此处假设当前目录下已存在test目录
	if err != nil {
		return
	}
	defer f.Close()

	io.Copy(f, file)
	res.Bucket = config.Bucket
	res.Driver = config.Driver
	res.Key = filekey
	res.Size = uint(header.Size)
	res.SSL = config.Ssl
	res.Name = header.Filename
	return res, nil
}

// 处理/upload 逻辑
func (ls *LocalStorage) Upload(r *http.Request) (result []Response, err error) {
	r.ParseMultipartForm(32 << 20)
	bucket := r.FormValue(ls.bucketKey)
	config, ok := slicekit.Find(ls.config, func(item StorageConf, index int, slice []StorageConf) bool {
		return item.Bucket == bucket
	})
	if !ok {
		config, ok = slicekit.Find(ls.config, func(item StorageConf, index int, slice []StorageConf) bool {
			return item.Primary
		})
		if !ok {
			err = errors.New("暂未找到配置配件")
		}
	}
	fileheaders := r.MultipartForm.File[ls.fileKey]
	result = make([]Response, 0)
	for _, header := range fileheaders {
		sufix := filepath.Ext(header.Filename)
		file, er := header.Open()
		if err != nil {
			err = er
			return
		}

		ctx := Context{
			File:       file,
			FileHeader: header,
			Ext:        sufix,
			TeamId:     "",
			UserId:     "",
		}
		filekey, _filepath := config.FileKeyAndPath(ctx)
		dir := filepath.Dir(_filepath)
		if _, e := os.Stat(dir); e != nil {
			if os.IsNotExist(e) {
				os.MkdirAll(dir, 0755)
			} else {
				err = e
				return
			}
		}
		f, e := os.OpenFile(_filepath, os.O_WRONLY|os.O_CREATE, 0666) // 此处假设当前目录下已存在test目录
		if e != nil {
			err = e
			return
		}
		io.Copy(f, file)
		f.Close()
		file.Close()
		res := Response{}
		res.Bucket = config.Bucket
		res.Driver = config.Driver
		res.Key = filekey
		res.Size = uint(header.Size)
		res.SSL = config.Ssl
		res.Name = header.Filename
		result = append(result, res)
	}
	return result, nil
}
