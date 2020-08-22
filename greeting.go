package hellopostgres

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
	"net/http"
)

type greeting struct {
	Id      int64  `json:"id" gorm:"primary_key"`
	Name    string `json:"name"`
	Prefix  string `json:"prefix"`
	Postfix string `json:"postfix"`
}

type dbConfig struct {
	Type    string `json:"type"`
	Host    string `json:"host"`
	Port    int32  `json:"port"`
	Db      string `json:"db"`
	User    string `json:"user"`
	Pass    string `json:"pass"`
	MinConn int32  `json:"minConn"`
	MaxConn int32  `json:"maxConn"`
}

type httpConfig struct {
	Port string `json:"port"`
}

type greetingResource struct {
	dbCfg   dbConfig
	httpCfg httpConfig
	db      *gorm.DB // could have made a service for this
}

type GreetingResource interface {
	Init() error
	Close()
}

func NewGreetingResource(dbCfg dbConfig, httpCfg httpConfig) GreetingResource {
	return &greetingResource{dbCfg: dbCfg, httpCfg: httpCfg}
}

func (g *greetingResource) Init() error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", g.dbCfg.Host, g.dbCfg.Port, g.dbCfg.User, g.dbCfg.Db, g.dbCfg.Pass)
	db, err := gorm.Open(g.dbCfg.Type, connStr)
	if err != nil {
		return errors.Wrap(err, "unable to parse db config url provided")
	}
	db.DB().SetMaxIdleConns(int(g.dbCfg.MinConn))
	db.DB().SetMaxOpenConns(int(g.dbCfg.MaxConn))
	g.db = db

	router := mux.NewRouter()
	router.HandleFunc("/hello", g.helloGet)
	router.HandleFunc("/hello/greeting", g.greetingPost)
	router.HandleFunc("/hello/greeting/{name}", g.greetingNameGet)
	return http.ListenAndServe(":"+g.httpCfg.Port, router)
}

func (g *greetingResource) Close() {
	if g.db != nil {
		_ = g.db.Close()
	}
}

func (g *greetingResource) helloGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	if _, err := w.Write([]byte("hello")); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
}

func (g *greetingResource) greetingNameGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	name := mux.Vars(r)["name"]

	greet := greeting{}
	g.db.Table("greeting").Where("name = ?", name).First(&greet)

	if len(greet.Prefix) <= 0 {
		greet.Prefix = "Hello"
	}

	if greet.Name == "Leo" {
		greet.Name = "Master Leo"
	}

	if _, err := w.Write([]byte(greet.Prefix + " " + greet.Name)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
}

func (g *greetingResource) greetingPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	greet := greeting{}
	if err := json.NewDecoder(r.Body).Decode(&greet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	g.db.Table("greeting").Create(&greet)

	w.Header().Set("Content-Type", "text/plain")
}
