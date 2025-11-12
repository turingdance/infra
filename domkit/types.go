package domkit

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
