package storage

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/turingdance/infra/alikit/osskit"
	"github.com/turingdance/infra/signer"
	"github.com/turingdance/infra/slicekit"
)

var ossClient *oss.Client = nil

func SignUrl(storageConfs []StorageConf, enpryt signer.Enpryt, ossConf osskit.OssConf, reskey string) (url string, err error) {

	conf, ok := slicekit.Find(storageConfs, func(item StorageConf, index int, slice []StorageConf) bool {
		return strings.HasPrefix(reskey, item.Bucket+"/")
	})
	if !ok {
		err = errors.New("暂时不支持")
		return
	}
	du := time.Second * 60 * 60
	if conf.Driver == DriverOSS {
		return SignOssUrl(conf, ossConf, reskey, du)
	}
	if conf.Driver == DriverLocal {
		return SignLocalUrl(conf, enpryt, reskey, du)
	}

	return url, err
}

func SignLocalUrl(conf StorageConf, enpryt signer.Enpryt, reskey string, du time.Duration) (url string, err error) {

	signer := signer.New(enpryt.Method, enpryt.Secret)

	// 如果是本地路径
	url = ""

	if conf.Ssl {
		url += "https://"
	} else {
		url += "http://"
	}
	url += conf.Host + "/"
	url += reskey
	if conf.AuthRequired {
		// 签名
		data := map[string]string{
			"requestURI": reskey,
		}
		expireat := time.Now().Add(du).Unix()
		sign, _ := signer.Sign(data, du)
		url += fmt.Sprintf("?sign=%s&expireAt=%d", sign, expireat)
	}

	return url, err
}
func SignOssUrl(conf StorageConf, ossConf osskit.OssConf, reskey string, du time.Duration) (url string, err error) {

	// 判断是否是oss

	if ossClient == nil {
		_ossClient, _err := oss.New(ossConf.Endpoint, ossConf.AccessKeyId, ossConf.AccessKeySecret)
		if _err != nil {
			err = _err
			return
		}
		ossClient = _ossClient
	}
	bucket, _err := ossClient.Bucket(conf.Bucket)
	if err != nil {
		err = _err
		return
	}
	url, err = bucket.SignURL(reskey, oss.HTTPGet, du)

	return url, err
}
