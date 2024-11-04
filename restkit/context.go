package restkit

import (
	"net/http"
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/turingdance/infra/restkit/validatekit"
	"github.com/turingdance/infra/wraper"
)

type (
	Context interface {
		// Request returns `*http.Request`.
		Request() *http.Request

		// SetRequest sets `*http.Request`.
		SetRequest(r *http.Request)

		// Response returns `*Response`.
		Writer() http.ResponseWriter

		SetWriter(w http.ResponseWriter)

		Params() *url.Values

		// Bind binds the request body into provided type `i`. The default binder
		// does it based on Content-Type header.
		Bind(ptr interface{}) error
		// Validate validates provided `i`. It is usually called after `Context#Bind()`.
		Validate(ptr interface{}) error
		// 校验
		RegisterValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool)
	}
	context struct {
		request  *http.Request
		writer   http.ResponseWriter
		params   *url.Values
		validate *validator.Validate
	}
)

func NewContext(w http.ResponseWriter, r *http.Request) *context {
	return &context{
		writer:   w,
		request:  r,
		params:   &url.Values{},
		validate: validator.New(),
	}
}

// 获得request
func (c *context) RegisterValidation(tag string, fn validator.Func, callValidationEvenIfNull ...bool) {
	c.validate.RegisterValidation(tag, fn, callValidationEvenIfNull...)
}

// 获得request
func (c *context) Request() *http.Request {
	return c.request
}

// 获得param
func (c *context) Params() *url.Values {
	return c.params
}

// 获得param
func (c *context) AddParam(key, val string) {
	c.params.Add(key, val)
}

// 设置request
func (c *context) SetRequest(r *http.Request) {
	c.request = r
}

// 获得Writer
func (c *context) Writer() http.ResponseWriter {
	return c.writer
}

// 初始化writer
func (c *context) SetWriter(w http.ResponseWriter) {
	c.writer = w
}

// 格式校验+自由绑定
func (c *context) Bind(ptrdata any) error {
	// 绑定并校验
	if err := wraper.Bind(c.request, ptrdata); err != nil {
		return err
	}
	return c.Validate(ptrdata)
}

func (c *context) Validate(ptrstru any) (err error) {
	if c.validate != nil {
		err = c.validate.Struct(ptrstru)
		err = validatekit.ProcessError(ptrstru, err) //处理错误信息
	}
	return
}
