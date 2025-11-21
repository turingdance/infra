package cond

const (
	Asc  = "asc"
	Desc = "desc"
)

type Order struct {
	Field  string `json:"field"`
	Method string `json:"method"`
}

func (o Order) Desc() string {
	return o.Field + " desc "
}
func (o Order) Asc() string {
	return o.Field + " asc "
}
func (o Order) Build() (string, error) {
	switch o.Method {
	case Asc:
		return o.Asc(), nil
	case Desc:
		return o.Desc(), nil
	default:
		return o.Field, nil
	}
}
