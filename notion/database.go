package notion

import "time"

type Database struct {
	Object         string     `json:"object"`
	ID             string     `json:"id"`
	CreatedTime    time.Time  `json:"created_time"`
	LastEditedTime time.Time  `json:"last_edited_time"`
	Title          []RichText `json:"title"`
	//Properties     Properties
}

func (d Database) DatabaseTitle() string {
	return d.Title[0].PlainText
}
