package notion

import (
	"encoding/json"
	"time"
)

type Parent struct {
	Type       string `json:"type"`
	DatabaseID string `json:"database_id"`
	PageID     string `json:"page_id"`
}

type Properties interface {
	Title() string
}

type PageTypeProperties struct {
	TitleValue TitlePropertyValue `json:"title"`
}

func (p PageTypeProperties) Title() string {
	return p.TitleValue.Title[0].PlainText
}

type WorkSpaceTypeProperties struct {
	TitleValue TitlePropertyValue `json:"title"`
}

func (w WorkSpaceTypeProperties) Title() string {
	return w.TitleValue.Title[0].PlainText
}

//type DatabaseTypeProperties map[string]Properties

type TitlePropertyValue struct {
	Id    string     `json:"id"`
	Type  string     `json:"type"`
	Title []RichText `json:"title"`
}

type Page struct {
	Object         string    `json:"object"`
	ID             string    `json:"id"`
	CreatedTime    time.Time `json:"created_time"`
	LastEditedTime time.Time `json:"last_edited_time"`
	Parent         Parent    `json:"parent"`
	Archived       bool      `json:"archived"`
	Properties     Properties
}

func (p *Page) UnmarshalJSON(data []byte) error {
	type Alias Page
	a := &struct {
		Properties json.RawMessage `json:"properties"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	switch p.Parent.Type {
	case "workspace":
		var properties WorkSpaceTypeProperties
		if err := json.Unmarshal(a.Properties, &properties); err != nil {
			return err
		}
		p.Properties = &properties
	case "page_id":
		var properties PageTypeProperties
		if err := json.Unmarshal(a.Properties, &properties); err != nil {
			return err
		}
		p.Properties = &properties
	}

	return nil
}
