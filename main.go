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
	ID     int    `json:"ID"`
	Title  string `json:"title"`
	Status string `json:"status"`
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

func main() {
	r := gin.Default()
	r.GET("api/todos", getTodosHandler)

	r.Run(":1234")
}
