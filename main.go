package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Task struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Status      TaskState `json:"status"`
}

type TaskState int

const (
	Todo TaskState = iota
	InProgress
	Done
)

var stateName = map[TaskState]string{
	Todo:       "todo",
	InProgress: "in-progress",
	Done:       "done",
}

func (ts TaskState) String() string {
	return stateName[ts]
}

func loadTask() ([]Task, error) {
	var tasks []Task

	if _, err := os.Stat("todo.json"); os.IsNotExist(err) {
		return tasks, nil
	}

	data, err := os.ReadFile("todo.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func saveTask(tasks []Task) error {
	data, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile("todo.json", data, 0644)
}

func addTask(description string) error {
	tasks, err := loadTask()
	if err != nil {
		return err
	}

	newTask := Task{
		Id:          len(tasks) + 1,
		Description: description,
		Status:      Todo,
	}

	tasks = append(tasks, newTask)
	return saveTask(tasks)
}

func updateTask(newDescription string, id int) error {

	tasks, err := loadTask()
	if err != nil {
		return err
	}

	var updatedTask Task

	for _, task := range tasks {
		if task.Id == id {
			updatedTask = Task{
				Id:          task.Id,
				Description: newDescription,
				Status:      task.Status,
			}
			break
		}
	}

	newTasks, err := deleteTask(updatedTask.Id)
	if err != nil {
		return err
	}

	newTasks = append(newTasks, updatedTask)
	return saveTask(newTasks)
}

func changeStatus(status TaskState, task Task, id int) error {

	updatedTask := Task{
		Id:          id,
		Description: task.Description,
		Status:      status,
	}

	newTasks, err := deleteTask(id)
	if err != nil {
		return err
	}

	newTasks = append(newTasks, updatedTask)
	return saveTask(newTasks)
}

func deleteTask(id int) ([]Task, error) {
	tasks, err := loadTask()
	if err != nil {
		return nil, err
	}

	for i, task := range tasks {
		if task.Id == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}

	return tasks, nil
}

func main() {
	operation := os.Args[1]

	switch operation {
	case "add":
		todoData := os.Args[2]
		addTask(todoData)

	case "list":
		tasks, err := loadTask()
		if err != nil {
			panic(err)
		}
		for _, v := range tasks {
			fmt.Printf("id: %v, description: %v, status: %v \n", v.Id, v.Description, v.Status)
		}

	case "update":
		taskId := os.Args[2]
		todoData := os.Args[3]
		taskIdNumber, err := strconv.Atoi(taskId)
		if err != nil {
			panic(err)
		}
		updateTask(todoData, taskIdNumber)

	case "mark-in-progress":
		taskId := os.Args[2]
		tasks, err := loadTask()
		if err != nil {
			panic(err)
		}

		taskIdNumber, err := strconv.Atoi(taskId)
		if err != nil {
			panic(err)
		}

		var taskToChange Task
		for _, task := range tasks {
			if task.Id == taskIdNumber {
				taskToChange = task
				break
			}
		}

		newStatus := transition(taskToChange.Status)
		changeStatus(newStatus, taskToChange, taskToChange.Id)

	case "mark-done":
		taskId := os.Args[2]
		tasks, err := loadTask()
		if err != nil {
			panic(err)
		}

		taskIdNumber, err := strconv.Atoi(taskId)
		if err != nil {
			panic(err)
		}

		var taskToChange Task
		for _, task := range tasks {
			if task.Id == taskIdNumber {
				taskToChange = task
				break
			}
		}

		newStatus := transition(taskToChange.Status)
		changeStatus(newStatus, taskToChange, taskToChange.Id)

	default:
		log.Fatal("Wrong operation")
	}

}

func transition(s TaskState) TaskState {
	switch s {
	case Todo:
		return InProgress
	case InProgress:
		return Done
	default:
		panic(fmt.Errorf("unknown state: %s", s))
	}
}
