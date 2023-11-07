package types

type Task struct {
	ID      int    `json:"ID"`
	Name    string `json:"Name"`
	Content string `json:"Content"`
	IsComplete bool `json:"IsComplete"`
}

type AllTasks []Task