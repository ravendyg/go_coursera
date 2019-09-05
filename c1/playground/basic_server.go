package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

func handleConnection(conn net.Conn) {
	name := conn.RemoteAddr().String()
	fmt.Printf("%+v connected\n", name)
	conn.Write([]byte("Hello, " + name + "\n\r"))

	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "Exit" {
			conn.Write([]byte("Bye\n\r"))
			fmt.Println(name, "disconnected")
			break
		} else if text != "" {
			fmt.Println(name, "enters", text)
			conn.Write([]byte("You entered " + text + "\n\r"))
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
	w.Write([]byte("!!!"))
}

func _main() {
	// listener, err := net.Listen("tcp", ":8080")
	// if err != nil {
	// 	panic(err)
	// }

	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	go handleConnection(conn)
	// }

	http.HandleFunc("/page",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Single page:", r.URL.String())
		})

	http.HandleFunc("/pages/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Multiple page:", r.URL.String())
		})

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
