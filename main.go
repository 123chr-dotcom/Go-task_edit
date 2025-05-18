package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sync"
)

type Task struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

var (
	tasks     []Task
	tasksLock sync.Mutex
	nextID    = 1
	dataFile  = "tasks.json"
)

func main() {
	loadTasks()
	
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/toggle", toggleHandler)
	
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	fmt.Println("任务管理服务启动在 :1145")
	http.ListenAndServe(":1145", nil)
}

func loadTasks() {
	file, err := os.ReadFile(dataFile)
	if err == nil {
		json.Unmarshal(file, &tasks)
		if len(tasks) > 0 {
			nextID = tasks[len(tasks)-1].ID + 1
		}
	}
}

func saveTasks() {
	data, _ := json.Marshal(tasks)
	os.WriteFile(dataFile, data, 0644)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tasksLock.Lock()
	defer tasksLock.Unlock()
	tmpl.Execute(w, tasks)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content := r.FormValue("content")
		if content != "" {
			tasksLock.Lock()
			task := Task{ID: nextID, Content: content, Done: false}
			tasks = append(tasks, task)
			nextID++
			saveTasks()
			tasksLock.Unlock()
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func toggleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		if id != "" {
			tasksLock.Lock()
			for i := range tasks {
				if fmt.Sprint(tasks[i].ID) == id {
					tasks[i].Done = !tasks[i].Done
					break
				}
			}
			saveTasks()
			tasksLock.Unlock()
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
