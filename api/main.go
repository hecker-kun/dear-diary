package main

import (
	"database/sql"
	"fmt"
	"github.com/baryon-m/dear-diary/api/handler"
	"github.com/baryon-m/dear-diary/config"
	"github.com/baryon-m/dear-diary/domain/entity/diary"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

func main() {
	psql := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DbHost,
		config.DbPort,
		config.DbUser,
		config.DbPassword,
		config.DbDatabase,
	)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	diaryRepo := diary.NewPostgreSQLRepoRepository(db)
	diaryService := diary.NewService(diaryRepo)

	r := mux.NewRouter()
	handler.MakeDiaryHandlers(r, diaryService)

	http.Handle("/", r)

	srv := &http.Server{
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr: "localhost:8080",
		Handler: context.ClearHandler(http.DefaultServeMux),
	}

	log.Println("Server started on :8080")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
