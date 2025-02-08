package models

type TemplateData struct {
	Post            *Post
	Posts           *[]Post
	Categories      []string
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
	User            *User
	NumberOfPage    int
	CurrentPage     int
	Limit           int
	Category        string
	Category_id     int
	URL             string
	LimitVariation  []int
	Quote           string
}
