package test

import (
	"app/internal/tag"
	"encoding/json"
)

type QuizeData struct {
	Text  string `json:"text"`
	IsRight bool   `json:"is_right"`
}

type QuizeTaskData struct {	
	Question string `json:"question"`
	Options []QuizeData `json:"options"`
}

type OpenQuestion struct {
	Description string `json:"description"`
}

type LifeCodeQuestion struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	StartTeplate string `json:"start_teplate"`
}

/*
Types:
	"quize"
	"oepn"
	"liveCode"
*/

type CreateTask struct {
	TestID int    `json:"test_id"`
	Type   string `json:"type"`
	Data   json.RawMessage `json:"data"` //[]TaskOption or OpenQuestion or LifeCodeQuestion
}

type Task struct {
	CreateTask
	TaskID int `json:"task_id"`
}

type CreateTest struct {
	Title     string `json:"title"`
	TimeLimit int    `json:"time_limit"`
	Type      string `json:"type"`
	Tasks     []CreateTask
}

type Test struct {
	CreateTest
	Tags   []tag.Tag
	TestID int `json:"test_id"`
}

type CreateTestRequest struct {
	Title     string       `json:"title"`
	TimeLimit int          `json:"time_limit"`
	Type      string       `json:"type"`
	Tasks     []CreateTask `json:"tasks"`
	TagIDs    []int        `json:"tag_ids"`
}
