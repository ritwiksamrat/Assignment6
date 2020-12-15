package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

var Articles []Article

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnSingleArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	ret, err := db.Query("select * from sampledb where id=(?)", key)
	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode(ret)

}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var article Article
	json.Unmarshal(reqBody, &article)

	//	Articles = append(Articles, article)

	result, err := db.Prepare("INSERT INTO sampledb(Id,title,desc,content) VALUES(?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	_, err = result.Exec(article.Id, article.Title, article.Desc, article.Content)
	if err != nil {
		panic(err.Error())
	}

	fmt.Fprintf(w, "Article Has Been Created!!!")
}

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	res, err := db.Query("delete from sampledb where id=(?)", id)
	if err != nil {
		panic(err.Error())
	}
	json.NewEncoder(w).Encode(res)

}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/article", createNewArticle).Methods("POST")
	myRouter.HandleFunc("/article/{id}", deleteArticle).Methods("DELETE")
	myRouter.HandleFunc("/article/{id}", returnSingleArticle).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

var db *sql.DB
var err error

func main() {

	db, err = sql.Open("mysql", "root:12345@tcp(127.0.0.1:3306)/sample")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("DataBase is Connected!!!")
	defer db.Close()

	Articles = []Article{
		Article{Id: "1", Title: "First Article", Desc: "This is my first article", Content: "first Content"},
		Article{Id: "2", Title: "Second Article", Desc: "This is my second article", Content: "second Content"},
	}
	handleRequests()
}
