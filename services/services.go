package services

import (
	"encoding/json"
	"os"
	"github.com/joho/godotenv"
	"log"
	"io/ioutil"
	"net/http"
	"strconv"
	"todoapi/Types"
	"github.com/gorilla/mux"
	"database/sql"
 	_ "github.com/lib/pq" 
)
var db *sql.DB
func init() {
	// Cargar las variables de entorno desde el archivo .env
	
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error: Unable to load .env file: %v", err)
	}

	// Obtener los valores de las variables de entorno
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	// Crear la cadena de conexión a la base de datos utilizando las variables de entorno
	connStr := "user=" + dbUser + " password=" + dbPass + " host=" + dbHost + " dbname=TodoDb sslmode=disable"

	var errDB error
	db, errDB = sql.Open("postgres", connStr)
	if errDB != nil {
		log.Fatalf("Error: Unable to connect to the database: %v", errDB)
	}
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rowsTareas, err := db.Query(`SELECT * FROM public."TodoItems"`)

	if err != nil {
		log.Fatalf("Error: Unable to execute the query: %v", err)
	}

	defer rowsTareas.Close()

	var tasks []types.Task
	for rowsTareas.Next() {
		var task types.Task
		err := rowsTareas.Scan(&task.ID, &task.Name, &task.IsComplete, &task.Content,)
		if err != nil {
			log.Fatalf("Error: Unable to scan row: %v", err)
			continue
		}
		tasks = append(tasks, task)
	}

	// Log the results to the console
	for _, task := range tasks {
		log.Printf("Task ID: %d, Name: %s, Content: %s\n", task.ID, task.Name, task.Content)
	}

	json.NewEncoder(w).Encode(tasks)
}


func CreateTask(w http.ResponseWriter, r *http.Request) {
    var newTask types.Task

    reqBody, err := ioutil.ReadAll(r.Body)

    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }

    json.Unmarshal(reqBody, &newTask)

    // Ejecuta una consulta de inserción en la base de datos
    _, err = db.Exec(`INSERT INTO public."TodoItems" ("Name", "Content", "IsComplete") VALUES ($1, $2, $3)`, newTask.Name, newTask.Content, newTask.IsComplete)

    if err != nil {
        log.Fatalf("Error: Unable to insert task into the database: %v", err)
        http.Error(w, "Failed to insert task into the database", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(newTask)
}
func UpdateTask(w http.ResponseWriter, r *http.Request) {
    // Extraer el ID de la URL o los parámetros de la solicitud
    vars := mux.Vars(r)
    taskID := vars["id"]

    var updatedTask types.Task

    reqBody, err := ioutil.ReadAll(r.Body)

    if err != nil {
        http.Error(w, "Failed to read request body", http.StatusBadRequest)
        return
    }

    err = json.Unmarshal(reqBody, &updatedTask)

    if err != nil {
        http.Error(w, "Failed to unmarshal JSON", http.StatusBadRequest)
        return
    }

    // Ejecuta una consulta de actualización en la base de datos
    _, err = db.Exec(`UPDATE public."TodoItems" SET "Name" = $1, "Content" = $2, "IsComplete" = $3 WHERE "Id" = $4`,
        updatedTask.Name, updatedTask.Content, updatedTask.IsComplete, taskID)

    if err != nil {
        log.Fatalf("Error: Unable to update task in the database: %v", err)
        http.Error(w, "Failed to update task in the database", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(updatedTask)
}


func DeleteTask(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    tasksID, err := strconv.Atoi(vars["id"])
    if err != nil {
        log.Println("Invalid id")
        return
    }

    // Construye la consulta SQL para eliminar una tarea por su ID
    query := `DELETE FROM public."TodoItems" WHERE public."TodoItems"."Id" = $1`

    // Ejecuta la consulta SQL para eliminar la tarea
    _, err = db.Exec(query, tasksID)
    if err != nil {
        log.Println("Error: Unable to delete task from the database:", err)
    } else {
        log.Printf("Task with ID %d has been deleted from the database\n", tasksID)
    }
}



func GetOneTask(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	tasksID, err := strconv.Atoi(vars["id"])
	if err != nil{
		log.Println(w, "invalid id")
	}
	query := `SELECT * FROM public."TodoItems" WHERE public."TodoItems"."Id" = $1`
	row := db.QueryRow(query, tasksID)
	var task types.Task
	err =row.Scan(&task.ID, &task.Name, &task.IsComplete, &task.Content,)
	if err != nil {
        log.Println("Error: Unable to retrieve task data from the database:", err)
    } else {
        // Imprime los datos en la consola
	w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(task)
	}
}
