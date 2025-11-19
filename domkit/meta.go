package domkit

type Column struct {
	Field  string   `json:"field"`
	Type   string   `json:"type"`
	Size   int      `json:"size"`
	Title  string   `json:"title"`
	Dom    DomeType `json:"dom"`
	Option any      `json:"option"`
}
type Meta struct {
	Table      string   `json:"table"`
	Title      string   `json:"title"`
	Model      string   `json:"model"`
	PrimaryKey string   `json:"primaryKey"`
	Column     []Column `json:"column"`
}
