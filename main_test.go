package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

var a App
var article Article

func TestMain(m *testing.M) {
	conf := "config.toml"
	//port = "3000"

	config := Config{}
	config.Read(conf)
	//config.Server = "localhost"
	config.Database = "articles_test"
	config.Collection = "full"

	a = App{}
	article = Article{
		Title:     "Test Article Title",
		Url:       "http://testurl.test/",
		Excerpt:   "This is a test excerpt for a test article.",
		DateSaved: "2017-01-01T00:00:00.000Z",
		Content:   "This is some test content for my test article.",
		Tags:      "Test,Tags",
	}
	article.ID = bson.NewObjectId()
	//article.ID = bson.ObjectIdHex("599f8e4ef1f58e55330c8b56")
	log.Println("Connecting to database...")
	a.Connect(config)
	ensureCollectionExists()
	log.Printf("Starting server...")
	code := m.Run()
	//clearCollection()
	os.Exit(code)
}

func TestEmptyCollection(t *testing.T) {
	clearCollection("1")
	req, _ := http.NewRequest("GET", "/articles", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	body := response.Body.String()
	if body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentArticle(t *testing.T) {
	addArticle()
	clearCollection("2")

	req, _ := http.NewRequest("GET", "/article/599f7a3ef1f58e72cffe0154", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Article not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Invalid Article ID'. Got '%s'", m["error"])
	}
}

func TestCreateArticle(t *testing.T) {
	clearCollection("3")
	payload := []byte(`{
		"title":"new title",
		"url":"new url",
		"excerpt":"general",
		"date_saved":"2017-01-01T00:00:00.000Z",
		"content":"some sample content"
	}`)

	req, _ := http.NewRequest("POST", "/article", bytes.NewBuffer(payload))
	response := executeRequest(req)

	if response.Code != http.StatusCreated {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusCreated, response.Code)
	}

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["url"] != "new url" {
		t.Errorf("Expected article url to be 'new url'. Got '%v'", m["url"])
	}

	if m["excerpt"] != "general" {
		t.Errorf("Expected article excerpt to be 'general'. Got '%v'", m["category"])
	}
}

func TestGetLinks(t *testing.T) {
	clearCollection("4")
	addArticle()

	req, _ := http.NewRequest("GET", "/articles", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

/*============  Utility Functions ==============================*/

func ensureCollectionExists() {
	if err := db.C(a.config.Collection).Insert(&article); err != nil {
		//if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal("INSERT TEST ARTICLE: ", err)
	}
}

func clearCollection(id string) {
	//a.DB.Exec("DELETE FROM links")
	if _, err := db.C(a.config.Collection).RemoveAll(bson.M{}); err != nil {
		log.Fatal("CLEAR COLLECTION "+id+": ", err)
	}
}

// func dropCollection()
// 	if err := a.dao.DropCollection(); err != nil {
// 		log.Fatal(err)
// 	}
// }

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func addArticle() {
	if err := db.C(a.config.Collection).Insert(article); err != nil {
		log.Fatal("ADD ARTICLE: ", err)
	}
}
