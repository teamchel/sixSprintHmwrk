package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// Обработчик получения всех задач
func getAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var taskList []Task
	for _, task := range tasks {
		taskList = append(taskList, task)
	}

	if err := json.NewEncoder(w).Encode(taskList); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Обработчик добавления новой задачи
func addTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Генерируем новый ID (увеличиваем на 1)
	newID := strconv.Itoa(len(tasks) + 1)
	newTask.ID = newID
	tasks[newID] = newTask

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newTask); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Обработчик получения задачи по ID
func getTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	task, exists := tasks[id]
	if !exists {
		http.Error(w, "Bad Request: Task not found", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Обработчик удаления задачи по ID
func deleteTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := chi.URLParam(r, "id")
	_, exists := tasks[id]
	if !exists {
		http.Error(w, "Bad Request: Task not found", http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Task deleted successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getAllTasks)
	r.Post("/tasks", addTask)
	r.Get("/tasks/{id}", getTaskByID)
	r.Delete("/tasks/{id}", deleteTaskByID)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
