package main

import (
	"fmt"
	"net/http"
	"text/template"
)

// USER -
type User struct {
	ID     int
	Name   string
	Active bool
}

func (u *User) PrintActive() string {
	if !u.Active {
		return ""
	}
	return "method says user " + u.Name + " active"
}

func IsUserOdd(u *User) bool {
	return u.ID%2 != 0
}

func main() {
	tmplFuncs := template.FuncMap{
		"OddUser": IsUserOdd,
	}
	tmpl, err := template.
		New("").
		Funcs(tmplFuncs).
		ParseFiles("users.html")
	// tmpl := template.Must(template.ParseFiles("users.html"))
	if err != nil {
		panic(err)
	}
	users := []User{
		User{1, "Vasily", true},
		User{2, "<i>Ivan</i>", false},
		User{3, "Dmitry", true},
	}

	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			tmpl.ExecuteTemplate(w, "users.html", struct {
				Users []User
			}{
				users,
			})
		})

	fmt.Println("starting server at :8089")
	http.ListenAndServe(":8089", nil)
}
