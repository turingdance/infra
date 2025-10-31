package dbkit

type DomeType string

const (
	Number    DomeType = "number"
	Text      DomeType = "text"
	RichText  DomeType = "richtext"
	Check     DomeType = "check"
	Radio     DomeType = "radio"
	Select    DomeType = "select"
	PickImage DomeType = "pickimage"
	PickFile  DomeType = "pickfile"
	Date      DomeType = "date"
	DateTime  DomeType = "datetime"
)

type Column struct {
	Field  string   `json:"field"`
	Type   string   `json:"type"`
	Size   int      `json:"size"`
	Title  string   `json:"title"`
	Dom    DomeType `json:"dom"`
	Option any      `json:"option"`
}
type Meta struct {
	Table  string   `json:"table"`
	Title  string   `json:"title"`
	Model  string   `json:"model"`
	Column []Column `json:"column"`
}
