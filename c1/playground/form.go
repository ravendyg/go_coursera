package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
)

var uploadFormTmpl = []byte(`
	<html>
	<body>
	<form action="/upload" method="post" enctype="multipart/form-data">
	Image: <input type="file" name="my_file">
	<input type="submit" value="Upload">
	</form>
	</body>
	</html>
	`)

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Write(uploadFormTmpl)
}

func uploadPage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(5 * 1024 * 1024)
	file, handler, err := r.FormFile("my_file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Fprintf(w, "handler.Filename %v\n", handler.Filename)
	fmt.Fprintf(w, "handler.Header %#v\n", handler.Header)

	hasher := md5.New()
	io.Copy(hasher, file)
	fmt.Println(w, "md5 %x\n", hasher.Sum(nil))
}

func __main() {
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/upload", uploadPage)

	http.ListenAndServe(":8089", nil)
}
