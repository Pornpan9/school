package main

import (
	"net/http"
	// "fmt"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
	"os"	
	"database/sql"
	"log"
)

func pingHandler(c *gin.Context) {

	// c.JSON(200, gin.H{
	// 	"message": 1,
	// })

	response := gin.H{
		"message": "This is Ping GET",
	}

	c.JSON(http.StatusOK, response)

}

func pingPostHandler(c *gin.Context) {

	response := gin.H{
		"message": "This is Ping POST",
	}

	c.JSON(http.StatusOK, response)

}

type Student struct {
	Name string `json:"name"`
	ID   int `json:"student_id"`
}

var students = map[int]Student{
	4749320: Student{Name: "Pornpan", ID: 4749320},
	4749321: Student{Name: "Noom", ID: 4749321},
}

func getStudentHandler(c *gin.Context) {

	student := []Student{}
	for _, s := range students{
		student = append(student, s)
	}
	c.JSON(http.StatusOK, student)//return slice
}

func postStudentHandler(c *gin.Context) {

	s:= Student{}

	// fmt.Printf("before bind % #v\n", s)
	if err := c.ShouldBindJSON(&s); err != nil{
		c.JSON(http.StatusBadRequest,err)
		return
	}
	// fmt.Printf("after bind % #v\n", s)

	id := len(students)
	id++
	s.ID = id
	students[id] = s

	c.JSON(http.StatusOK, students)//return map
}

type Todo struct{
	ID int
	Title string
	Status string
}

func getTodosHandler(c *gin.Context){

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil{
		log.Fatal("faltal : ", err.Error())
	}

	stmt, err := db.Prepare("SELECT id, title, status FROM todos")
	if err != nil{
		log.Fatal("Prepare error : ", err.Error())
	}

	rows, err := stmt.Query()
	if err != nil{
		log.Fatal("error : ", err.Error())
	}

	todos := []Todo{}
	for rows.Next(){
		t := Todo{}
		err := rows.Scan(&t.ID, &t.Title, &t.Status)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error : ": err.Error()})
			return
		}
		todos = append(todos, t)		
	}

	c.JSON(http.StatusOK, todos)
}

func main() {

	r := gin.Default()

	r.GET("/ping", pingHandler)
	r.POST("/ping", pingPostHandler)

	r.GET("/students", getStudentHandler)
	r.POST("/students", postStudentHandler)

	r.GET("/api/todos", getTodosHandler)

	r.Run(":1234")
	// r.Run() // listen and serve on 0.0.0.0:8080

}
