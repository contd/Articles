package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

type Page struct {
	Title    string
	Author   string
	Articles *[]SumArticle
}

type App struct {
	Router *mux.Router
	dao    ArticlesDAO
	config Config
}

func (a *App) Connect(config Config) {
	a.config = config //.Read(conf)
	a.dao.Server = a.config.Server
	a.dao.Database = a.config.Database
	a.dao.Collection = a.config.Collection
	a.dao.Connect()
	a.Router = mux.NewRouter()
	a.InitializeRoutes()
}

func (a *App) InitializeRoutes() {
	a.Router.HandleFunc("/", a.HtmlFrontEndPoint).Methods("GET")
	a.Router.HandleFunc("/articles", a.AllArticlesEndPoint).Methods("GET")
	a.Router.HandleFunc("/article", a.CreateArticleEndPoint).Methods("POST")
	a.Router.HandleFunc("/article", a.UpdateArticleEndPoint).Methods("PUT")
	a.Router.HandleFunc("/article", a.DeleteArticleEndPoint).Methods("DELETE")
	a.Router.HandleFunc("/article/{id}", a.FindArticleEndpoint).Methods("GET")
}

func (a *App) Run(port string) {
	log.Printf("Server running on http://localhost:%s/", port)
	corsObj := handlers.AllowedOrigins([]string{"*"})
	log.Println(http.ListenAndServe(":"+port, handlers.CORS(corsObj)(a.Router)))
	//if err := http.ListenAndServe(":"+port, handlers.CORS(corsObj)(a.Router)); err != nil {
	//	log.Fatal(err)
	//}
}

func (a *App) HtmlFrontEndPoint(w http.ResponseWriter, r *http.Request) {
	articles, err := a.dao.FindAll()
	if err != nil {
		log.Fatal("Error FindAll: ", err)
	}
	page := &Page{Title: "Articles Saved", Author: "Jason Kumpf", Articles: &articles}
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, page)
}

func (a *App) AllArticlesEndPoint(w http.ResponseWriter, r *http.Request) {
	articles, err := a.dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(articles) <= 0 {
		articles = []SumArticle{}
	}
	respondWithJson(w, http.StatusOK, articles)
}

func (a *App) FindArticleEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	article, err := a.dao.FindById(params["id"])

	if err.Error() == "not found" {
		respondWithError(w, http.StatusNotFound, "Article not found")
		return
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Article ID")
		return
	}

	respondWithJson(w, http.StatusOK, article)
}

func (a *App) CreateArticleEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var article Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	article.ID = bson.NewObjectId()
	if err := a.dao.Insert(article); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, article)
}

func (a *App) UpdateArticleEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var article Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := a.dao.Update(article); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) DeleteArticleEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var article Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := a.dao.Delete(article); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
