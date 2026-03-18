package notion

import "time"

// Notion API response types

type DatabaseQueryResponse struct {
	Results    []Page  `json:"results"`
	HasMore    bool    `json:"has_more"`
	NextCursor *string `json:"next_cursor"`
}

type Page struct {
	ID             string              `json:"id"`
	CreatedTime    time.Time           `json:"created_time"`
	LastEditedTime time.Time           `json:"last_edited_time"`
	Properties     map[string]Property `json:"properties"`
	URL            string              `json:"url"`
	Cover          *FileCover          `json:"cover"`
}

type Property struct {
	Type        string        `json:"type"`
	Title       []RichText    `json:"title,omitempty"`
	RichText    []RichText    `json:"rich_text,omitempty"`
	Select      *SelectValue  `json:"select,omitempty"`
	MultiSelect []SelectValue `json:"multi_select,omitempty"`
	Number      *float64      `json:"number,omitempty"`
	Files       []FileObject  `json:"files,omitempty"`
}

type RichText struct {
	PlainText string `json:"plain_text"`
}

type SelectValue struct {
	Name string `json:"name"`
}

type FileObject struct {
	Type     string   `json:"type"`
	File     *FileURL `json:"file,omitempty"`
	External *FileURL `json:"external,omitempty"`
}

type FileURL struct {
	URL string `json:"url"`
}

type FileCover struct {
	Type     string   `json:"type"`
	File     *FileURL `json:"file,omitempty"`
	External *FileURL `json:"external,omitempty"`
}

// Block types

type BlocksResponse struct {
	Results    []Block `json:"results"`
	HasMore    bool    `json:"has_more"`
	NextCursor *string `json:"next_cursor"`
}

type Block struct {
	ID   string `json:"id"`
	Type string `json:"type"`

	Heading2         *HeadingBlock   `json:"heading_2,omitempty"`
	Heading3         *HeadingBlock   `json:"heading_3,omitempty"`
	Paragraph        *ParagraphBlock `json:"paragraph,omitempty"`
	BulletedListItem *ListItemBlock  `json:"bulleted_list_item,omitempty"`
	NumberedListItem *ListItemBlock  `json:"numbered_list_item,omitempty"`
	Callout          *CalloutBlock   `json:"callout,omitempty"`
	Image            *ImageBlock     `json:"image,omitempty"`
}

type HeadingBlock struct {
	RichText []RichText `json:"rich_text"`
}

type ParagraphBlock struct {
	RichText []RichText `json:"rich_text"`
}

type ListItemBlock struct {
	RichText []RichText `json:"rich_text"`
}

type CalloutBlock struct {
	RichText []RichText `json:"rich_text"`
}

type ImageBlock struct {
	Type     string   `json:"type"`
	File     *FileURL `json:"file,omitempty"`
	External *FileURL `json:"external,omitempty"`
}
