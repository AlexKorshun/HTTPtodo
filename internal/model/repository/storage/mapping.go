package storage

import "github.com/AlexKorshun/HTTPtodo/internal/model"

type JsonTask struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func (j *JsonTask) convToTask() model.Task {
	return model.Task{ID: j.ID, Text: j.Text, Done: j.Done}
}

func convArrayToTask(arrayJsonTasks []JsonTask) []model.Task {
	var array []model.Task
	for _, value := range arrayJsonTasks {
		array = append(array, value.convToTask())
	}
	return array
}
