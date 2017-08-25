package main

import (
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type SumArticle struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Title     string        `bson:"title" json:"title"`
	Url       string        `bson:"url" json:"url"`
	Excerpt   string        `bson:"excerpt" json:"excerpt"`
	DateSaved string        `bson:"date_saved" json:"date_saved"`
	Tags      string        `bson:"tags" json:"tags"`
}

type Article struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Title     string        `bson:"title" json:"title"`
	Url       string        `bson:"url" json:"url"`
	Excerpt   string        `bson:"excerpt" json:"excerpt"`
	DateSaved string        `bson:"date_saved" json:"date_saved"`
	Content   string        `bson:"content" json:"content"`
	Tags      string        `bson:"tags" json:"tags"`
}

type ArticlesDAO struct {
	Server     string
	Database   string
	Collection string
}

var db *mgo.Database

func (a *ArticlesDAO) Connect() {
	session, err := mgo.Dial(a.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(a.Database)
}

func (a *ArticlesDAO) FindAll() ([]SumArticle, error) {
	var articles []SumArticle
	err := db.C(a.Collection).Find(bson.M{}).All(&articles)
	return articles, err
}

func (a *ArticlesDAO) FindById(id string) (Article, error) {
	var article Article
	err := db.C(a.Collection).FindId(bson.ObjectIdHex(id)).One(&article)
	return article, err
}

func (a *ArticlesDAO) Insert(article Article) error {
	err := db.C(a.Collection).Insert(&article)
	return err
}

func (a *ArticlesDAO) Delete(article Article) error {
	err := db.C(a.Collection).Remove(&article)
	return err
}

func (a *ArticlesDAO) Update(article Article) error {
	err := db.C(a.Collection).UpdateId(article.ID, &article)
	return err
}

func (a *ArticlesDAO) RemoveAll() (*mgo.ChangeInfo, error) {
	info, err := db.C(a.Collection).RemoveAll(bson.M{})
	return info, err
}

func (a *ArticlesDAO) DropCollection() error {
	err := db.C(a.Collection).DropCollection()
	return err
}
