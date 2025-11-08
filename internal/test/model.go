package test

import (
	"app/internal/tag"
)

type TaskOption struct {
	Option string `json:"option"`
	IsRight bool `json:"is_right"`
}

type OpenQuestion struct {
	Description string `json:"description"`
}

type CreateTask struct {
	TestID int  `json:"test_id"`
	Type string `json:"type"`
	Options string `json:"options"` //[]TaskOption or OpenQuestion, 
}

type Task struct {
	CreateTask
	TaskID int  `json:"task_id"`

}

type CreateTest struct {
	Title string `json:"title"`
	TimeLimit int `json:"time_limit"`
	Type string `json:"type"`
	Tasks []CreateTask
}

type Test struct {
	CreateTest
	Tags []tag.Tag
	TestID int `json:"test_id"`
}


type CreateTestRequest struct {
	Title     string       `json:"title"`
	TimeLimit int          `json:"time_limit"`
	Type      string       `json:"type"`
	Tasks     []CreateTask `json:"tasks"`
	TagIDs    []int        `json:"tag_ids"`
}