package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type Response struct {
	Rid     string `json:"id"`
	Rtitle  string `json:"title"`
	Rstatus string `json:"status"`
}

func getTodosHandler(c *gin.Context) {
	db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	stmt, err := db.Prepare("SELECT  id, title, status FROM todos")
	if err != nil {
		log.Fatal("stmt error ", err.Error())
	}
	defer db.Close()

	rows, _ := stmt.Query()
	todos := []Todo{}

	for rows.Next() {
		t := Todo{}
		err := rows.Scan(&t.ID, &t.Title, &t.Status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error ": err.Error()})
			return
		}
		todos = append(todos, t)
	}
	fmt.Println("YEHHHHHHH !!! connect to database already !!!!")
	fmt.Println(todos)
	c.JSON(200, todos)
	return
}

func postTodosHandler(c *gin.Context) {
	url := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal("faltal", err.Error())
	}
	defer db.Close()

	t := Todo{}
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	fmt.Println(t)

	title := t.Title
	status := t.Status

	query := `
	INSERT INTO todos (title, status) VALUES ($1, $2) RETURNING id
	`
	var id int

	row := db.QueryRow(query, title, status)
	err = row.Scan(&id)

	if err != nil {
		log.Fatal("Can't scan id", err.Error())
	}
	fmt.Println("insert sucess id : ", id)
	t.ID = id
	c.JSON(201, t)
	//	c.JSON(http.StatusOK, t)
	return
}

func getTodosByIdHandler(c *gin.Context) {
	idInput := c.Param("id")

	db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	stmt, _ := db.Prepare("SELECT id, title, status FROM todos WHERE id=$1")

	row := stmt.QueryRow(idInput)
	t := Todo{}
	r := Response{}

	err := row.Scan(&t.ID, &t.Title, &t.Status)
	if err != nil {
		log.Fatal("error", err.Error())
	}

	title := t.Title
	status := t.Status
	r.Rid = idInput
	r.Rstatus = t.Status
	r.Rtitle = t.Title

	fmt.Println("one row ", idInput, title, status)

	fmt.Println("Select by ID !!!!")
	fmt.Println(t)
	fmt.Println(r)
	c.JSON(200, t)
	return
}

func putTodosByIdHandler(c *gin.Context) {
	idInput := c.Param("id")

	db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	stmt, err := db.Prepare("UPDATE todos SET status=$2 WHERE id=$1;")

	if err != nil {
		log.Fatal("Can't scan id", err.Error())
	}
	if _, err3 := stmt.Exec(idInput, "inactive"); err3 != nil {
		log.Fatal("ex error", err3.Error())
	}

	stmt2, _ := db.Prepare("SELECT id, title, status FROM todos WHERE id=$1")

	row := stmt2.QueryRow(idInput)
	t := Todo{}

	err2 := row.Scan(&t.ID, &t.Title, &t.Status)
	if err2 != nil {
		log.Fatal("error", err2.Error())
	}

	title := t.Title
	status := t.Status

	fmt.Println("one row ", idInput, title, status)

	fmt.Println("update sucess id : ", idInput)
	fmt.Println("Select by ID !!!!", idInput)
	c.JSON(200, t)
	return
}

func deleteTodosByIdHandler(c *gin.Context) {
	idInput := c.Param("id")

	db, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	//	stmt2, err4 := db.Prepare("SELECT id, title, status FROM todos WHERE id=$1")

	//	if err4 != nil {
	//		log.Fatal("Can't scan id", err4.Error())
	//	}

	//	row := stmt2.QueryRow(idInput)
	//	t := Todo{}

	//	err2 := row.Scan(&t.ID, &t.Title, &t.Status)
	//	if err2 != nil {
	//		log.Fatal("error", err2.Error())
	//	}

	//	title := t.Title
	//	status := t.Status

	t := Todo{}

	fmt.Println("row for delete ", idInput)

	stmt, err := db.Prepare("DELETE FROM todos  WHERE id=$1 RETURNING Id,Status,Title;")

	if err != nil {
		log.Fatal("Can't scan id==> ", err.Error())
	}

	fmt.Println("one row ", idInput)

	fmt.Println("Select by ID !!!!", idInput, stmt, t)
	c.JSON(200, gin.H{"status": "success"})
	//c.JSON(http.StatusOK, "s")
	return
}

func main() {
	r := gin.Default()
	r.GET("api/todos", getTodosHandler)
	r.POST("api/todos", postTodosHandler)
	r.GET("api/todos/:id", getTodosByIdHandler)
	r.PUT("api/todos/:id", putTodosByIdHandler)
	r.DELETE("api/todos/:id", deleteTodosByIdHandler)

	r.Run(":1234")
}
