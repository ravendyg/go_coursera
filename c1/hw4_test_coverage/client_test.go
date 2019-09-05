package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

type user struct {
	XMLName xml.Name `xml:"row"`
	ID      int      `xml:"id" json:"id"`
	Name    string   `xml:"first_name" json:"name"`
	Age     int      `xml:"age" json:"age"`
	About   string   `xml:"about" json:"about"`
	Gender  string   `xml:"gender" json:"gender"`
}

type root struct {
	XMLName xml.Name `xml:"root"`
	Rows    []user   `xml:"row"`
}

type Handler func(w http.ResponseWriter, r *http.Request)

type TestCase struct {
	request  SearchRequest
	response SearchResponse
	handler  Handler
}

func SearchServer(handler Handler) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handler))
}

func TestGetUser(t *testing.T) {
	request := SearchRequest{
		Limit:      1,
		Offset:     0,
		OrderBy:    OrderByAsIs,
		OrderField: "Age",
		Query:      "",
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		xmlFile, err := os.Open("dataset.xml")
		if err != nil {
			panic(err)
		}
		byteVal, err := ioutil.ReadAll(xmlFile)
		if err != nil {
			panic(err)
		}

		var rows root
		xml.Unmarshal(byteVal, &rows)

		a, err := json.Marshal(rows.Rows[0:1])
		if err != nil {
			panic(err)
		}

		w.Write(a)

		defer xmlFile.Close()
	}

	server := SearchServer(handler)
	defer server.Close()

	client := SearchClient{
		AccessToken: "",
		URL:         server.URL,
	}
	resp, err := client.FindUsers(request)
	if err != nil {
		t.Errorf("unexpected error %#v", err)
	}
	if len(resp.Users) != 1 || resp.NextPage {
		t.Error("should find exactly one user")
	}
	if resp.Users[0].Name != "Boyd" {
		t.Error("found wrong user")
	}
}

func TestGetMaxNumberOfUsers(t *testing.T) {
	request := SearchRequest{
		Limit:      30,
		Offset:     0,
		OrderBy:    OrderByAsIs,
		OrderField: "Age",
		Query:      "",
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		xmlFile, err := os.Open("dataset.xml")
		if err != nil {
			panic(err)
		}
		byteVal, err := ioutil.ReadAll(xmlFile)
		if err != nil {
			panic(err)
		}
		query := r.URL.Query()
		// use default int == 0
		offset, _ := strconv.Atoi(query["offset"][0])
		limit, _ := strconv.Atoi(query["limit"][0])
		last := offset + limit
		if last == 0 {
			last = 1
		}
		var rows root
		xml.Unmarshal(byteVal, &rows)

		a, err := json.Marshal(rows.Rows[offset:last])
		if err != nil {
			panic(err)
		}

		w.Write(a)

		defer xmlFile.Close()
	}

	server := SearchServer(handler)
	defer server.Close()

	client := SearchClient{
		AccessToken: "",
		URL:         server.URL,
	}
	resp, err := client.FindUsers(request)
	if err != nil {
		t.Errorf("unexpected error")
	}
	if len(resp.Users) != 25 {
		t.Error("should find exactly 25 users")
	}
	if !resp.NextPage {
		t.Error("there should be the next page")
	}
}

func TestUnauthorized(t *testing.T) {
	request := SearchRequest{
		Limit:      1,
		Offset:     0,
		OrderBy:    OrderByAsIs,
		OrderField: "Age",
		Query:      "",
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}

	server := SearchServer(handler)
	defer server.Close()

	client := SearchClient{
		AccessToken: "",
		URL:         server.URL,
	}
	_, err := client.FindUsers(request)
	if err.Error() != "Bad AccessToken" {
		t.Errorf("wrong error %#v", err)
	}
}

func TestErrorBadJson(t *testing.T) {
	request := SearchRequest{
		Limit:      1,
		Offset:     0,
		OrderBy:    OrderByAsIs,
		OrderField: "Phone",
		Query:      "",
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{"))
	}

	server := SearchServer(handler)
	defer server.Close()

	client := SearchClient{
		AccessToken: "",
		URL:         server.URL,
	}
	_, err := client.FindUsers(request)
	if err.Error() != "cant unpack error json: unexpected end of JSON input" {
		t.Errorf("wrong error %#v", err)
	}
}

func TestResultBadJson(t *testing.T) {
	request := SearchRequest{
		Limit:      1,
		Offset:     0,
		OrderBy:    OrderByAsIs,
		OrderField: "Phone",
		Query:      "",
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{"))
	}

	server := SearchServer(handler)
	defer server.Close()

	client := SearchClient{
		AccessToken: "",
		URL:         server.URL,
	}
	_, err := client.FindUsers(request)
	if err.Error() != "cant unpack result json: unexpected end of JSON input" {
		t.Errorf("wrong error %#v", err)
	}
}

func TestErrorBadOrderField(t *testing.T) {
	request := SearchRequest{
		Limit:      1,
		Offset:     0,
		OrderBy:    OrderByAsIs,
		OrderField: "Phone",
		Query:      "",
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\": \"ErrorBadOrderField\"}"))
	}

	server := SearchServer(handler)
	defer server.Close()

	client := SearchClient{
		AccessToken: "",
		URL:         server.URL,
	}
	_, err := client.FindUsers(request)
	if err.Error() != "OrderField Phone invalid" {
		t.Errorf("wrong error %#v", err)
	}
}

func TestErrorFatal(t *testing.T) {
	request := SearchRequest{
		Limit:      1,
		Offset:     0,
		OrderBy:    OrderByAsIs,
		OrderField: "Phone",
		Query:      "",
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	server := SearchServer(handler)
	defer server.Close()

	client := SearchClient{
		AccessToken: "",
		URL:         server.URL,
	}
	_, err := client.FindUsers(request)
	if err.Error() != "SearchServer fatal error" {
		t.Errorf("wrong error %#v", err)
	}
}

func TestErrorUnknown(t *testing.T) {
	request := SearchRequest{
		Limit:      1,
		Offset:     0,
		OrderBy:    OrderByAsIs,
		OrderField: "Phone",
		Query:      "",
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"error\": \"UnknownError\"}"))
	}

	server := SearchServer(handler)
	defer server.Close()

	client := SearchClient{
		AccessToken: "",
		URL:         server.URL,
	}
	_, err := client.FindUsers(request)
	fmt.Println(err.Error())
	if err.Error() != "unknown bad request error: UnknownError" {
		t.Errorf("wrong error %#v", err)
	}
}

func TestIncorrectLimit(t *testing.T) {
	request := SearchRequest{
		Limit:      -1,
		Offset:     0,
		OrderBy:    OrderByAsIs,
		OrderField: "Phone",
		Query:      "",
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{"))
	}

	server := SearchServer(handler)
	defer server.Close()

	client := SearchClient{
		AccessToken: "",
		URL:         server.URL,
	}
	_, err := client.FindUsers(request)
	if err.Error() != "limit must be > 0" {
		t.Errorf("wrong error %#v", err)
	}
}

func TestIncorrectOffset(t *testing.T) {
	request := SearchRequest{
		Limit:      1,
		Offset:     -1,
		OrderBy:    OrderByAsIs,
		OrderField: "Phone",
		Query:      "",
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{"))
	}

	server := SearchServer(handler)
	defer server.Close()

	client := SearchClient{
		AccessToken: "",
		URL:         server.URL,
	}
	_, err := client.FindUsers(request)
	if err.Error() != "offset must be > 0" {
		t.Errorf("wrong error %#v", err)
	}
}
