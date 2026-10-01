package main

import (
	"log"
	"net/http"
	"os"
	"restApiGo/internal/database"
	"restApiGo/internal/handlers"
	"restApiGo/internal/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		databaseUrl = "postgres://taskuser:taskpass@localhost:5432/tasksdb?sslmode=disable"
	}
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}
	log.Printf("Starting server on port %s", serverPort)
	db, err := database.Connect(databaseUrl)
	createFiber(serverPort)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Printf("Database connection established")

	taskStore := database.NewTaskStore(db)
	handler := handlers.NewHandler(taskStore)
	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", methodHandler(handler.GetAllTasks, http.MethodGet))
	mux.HandleFunc("/tasks/create", methodHandler(handler.CreateTask, http.MethodPost))
	mux.HandleFunc("/tasks/", taskIdHandler(handler))

	loggedMux := loggingMiddleware(mux)
	serverAddr := ":" + serverPort
	err = http.ListenAndServe(serverAddr, loggedMux)
}

func loggingMiddleware(mux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		mux.ServeHTTP(w, r)
	})
}

func methodHandler(handlerFunc http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != allowedMethod {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
		handlerFunc(w, r)
	}
}
func taskIdHandler(handler *handlers.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetTaskById(w, r)
		case http.MethodPost:
			handler.CreateTask(w, r)
		case http.MethodPut:
			handler.UpdateTaskById(w, r)
		case http.MethodDelete:
			handler.DeleteTaskById(w, r)
		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	}
}
func createFiber(port string) {
	app := fiber.New()
	routes.Setup(app)
	app.Listen(":" + port)

}
