package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"database/sql"
	"time"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/go-playground/validator.v9"
)

const DB_PATH string = "./foo.db"

var validate *validator.Validate

type FormUserInfo struct {
	Name    string     `validate:"required"`
	Email   string     `validate:"required,email"`
	Message string
	Errors  []string
}

type DbUserInfo struct {
	Name    sql.NullString
	Email   sql.NullString
	Message sql.NullString
	Created time.Time
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			// TODO Treat it better
			fmt.Fprintf(w, "Something went wrong when saving. Error message: ", r)
		}
	}()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("form.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// TODO extend the validation to limit the size of the input
		validate = validator.New()
		userInfo, err := validateFormUserInfo(r)
		if err != nil {
			userInfo.Errors = getErrorMessages(err)
			t, _ := template.ParseFiles("form.gtpl")
			t.Execute(w, userInfo)
		} else {
			id := saveToDb(userInfo.Name, userInfo.Email, userInfo.Message)
			fmt.Println(w, "Wrote to the database: %s, %s. With ID: %d", r.Form["name"], r.Form["message"], id)
			handleList(w, r)
		}
	}
}

func validateFormUserInfo(r *http.Request) (*FormUserInfo, error) {
	validate = validator.New()
	userInfo := &FormUserInfo{
		Name: r.Form["name"][0],
		Email: r.Form["email"][0],
		Message:  r.Form["message"][0],
	}
	err := validate.Struct(userInfo)
	return userInfo, err
}

func getErrorMessages(err error) []string {
	all := []string{}
	for _, err := range err.(validator.ValidationErrors) {
		//TODO Display more friendly/informative messages
		all = append(all, err.Field())
	}
	return all
}

func saveToDb(name string, email string, message string) int64 {
	db, err := sql.Open("sqlite3", DB_PATH)
	fmt.Println(name, email, message)
	PanicIf(err)
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO form_info(name, email, message) values(?,?,?)")
	PanicIf(err)
	res, err := stmt.Exec(name, email, message)
	PanicIf(err)
	id, err := res.LastInsertId()
	PanicIf(err)
	return id
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func handleList(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			// TODO Treat it better
			fmt.Fprintf(w, "Something went wrong when listing. Error message: ", r)
		}
	}()
	info := listFromDb()
	t, _ := template.ParseFiles("list.gtpl")
	t.Execute(w, info)
}

func listFromDb() []DbUserInfo {
	result := []DbUserInfo{}
	db, _ := sql.Open("sqlite3", DB_PATH)
	defer db.Close()
	rows, err := db.Query("SELECT name, email, message FROM form_info")
	PanicIf(err)
	for rows.Next() {
		var name sql.NullString
		var email sql.NullString
		var message sql.NullString
		err = rows.Scan(&name, &email, &message)
		PanicIf(err)
		result = append(result, DbUserInfo{Name: name, Email: email, Message: message})
	}
	fmt.Println(result)
	return result
}

func createTableIfNeeded() {
	db, err := sql.Open("sqlite3", DB_PATH)
	PanicIf(err)
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS form_info (" +
		"id INTEGER PRIMARY KEY AUTOINCREMENT, " +
		"name TEXT NOT NULL, " +
		"email TEXT NOT NULL, " +
		"message TEXT, " +
		"created TIMESTAMP DEFAULT CURRENT_TIMESTAMP)")
	PanicIf(err)
	statement.Exec()
	PanicIf(err)
}

func main() {
	createTableIfNeeded()
	http.HandleFunc("/", handleForm)
	http.HandleFunc("/save", handleForm)
	http.HandleFunc("/list", handleList)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}