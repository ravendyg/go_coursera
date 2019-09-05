package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"
)

// Post -
type Post struct {
	ID       int
	Text     string
	Author   string
	Comments int
	Time     time.Time
}

func handle(w http.ResponseWriter, req *http.Request) {
	s := ""
	for i := 0; i < 1000; i++ {
		p := &Post{ID: i, Text: "new post"}
		s += fmt.Sprintf("%#v", p)
	}
	w.Write([]byte(s))
}
func main() {
	http.HandleFunc("/", handle)
	fmt.Println("starting server at :8089")
	fmt.Println(http.ListenAndServe(":8089", nil))
}
