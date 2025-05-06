package main

import (
	"time"
)

type Course struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
	Instructor  string `json:"instructor"`
	Category    string `json:"category"`
	CoverImg    string `json:"coverImg"`
	IsPublished bool   `json:"isPublished"`
	Sections    []Section
	CreateAt    time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
type Section struct {
	Title     string `json:"title" bson:"title"`
	SubTopics []SubTopic
}

type SubTopic struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
