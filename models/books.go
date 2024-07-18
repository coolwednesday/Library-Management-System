package models

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Isbn   int    `json:"isbn"`
}

type LendingRecord struct {
	UserID int `json:"userid"`
	ISBN   int `json:"isbn"`
}
