package models

type User struct {
	Id      int
	Name    string `json:"name"`
	Age     string `json:"age"`
	Friends []User `json:"friends"`
}
type Frds struct {
	SourceID string `json:"source_id"`
	TargetID string `json:"target_id"`
}
type Age struct {
	Age string `json:"new_age"`
}
