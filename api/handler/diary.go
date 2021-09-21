package handler

import (
	"encoding/json"
	"github.com/baryon-m/dear-diary/domain/entity"
	"github.com/baryon-m/dear-diary/domain/entity/diary"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func createEntry(service diary.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errMessage := "Error creating diary entry"

		var input struct{
			Author string
			Content string
		}
		err := json.NewDecoder(r.Body).Decode(&input)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMessage))
			return
		}

		e := &diary.Entry{
			ID: entity.NewID(),
			Author: input.Author,
			Content: input.Content,
			CreatedAt: time.Now().Format(time.RFC3339),
		}

		e.ID, err = service.Create(e)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMessage))
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(e); err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMessage))
			return
		}
	})
}

func deleteEntry(service diary.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errMessage := "Error when deleting a diary entry"
		vars := mux.Vars(r)

		id, err := entity.StringToID(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMessage))
			return
		}

		err = service.Delete(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMessage))
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	})
}

func fetchOneEntry(service diary.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errMessage := "Error fetching diary entry"
		vars := mux.Vars(r)

		id, err := entity.StringToID(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMessage))
			return
		}

		data, err := service.FetchOne(id)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMessage))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errMessage))
			return
		}

		toJSON := diary.Entry{
			ID: data.ID,
			Author: data.Author,
			Content: data.Content,
			CreatedAt: data.CreatedAt,
		}
		if err := json.NewEncoder(w).Encode(toJSON); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMessage))
		}
	})
}

func fetchAllEntries(service diary.UseCase) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		errMessage := "Error fetching a list of diary entries"

		data, err := service.FetchAll()
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMessage + "- error Content-Type set"))
			return
		}

		if data == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errMessage + "- data is nil"))
			return
		}

		var toJSON = struct {
			Entries []*diary.Entry
		}{}
		for _, e := range data {

			toJSON.Entries = append(toJSON.Entries, e)
		}

		if err := json.NewEncoder(w).Encode(toJSON); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMessage + "- JSON Encoding error"))
		}
	})
}

func MakeDiaryHandlers(r *mux.Router, service diary.UseCase) {
	r.Handle("/diary", createEntry(service)).Methods("POST").Name("createEntry")
	r.Handle("/diary/{id}", fetchOneEntry(service)).Methods("GET").Name("fetchOneEntry")
	r.Handle("/diary", fetchAllEntries(service)).Methods("GET").Name("fetchAllEntries")
	r.Handle("/diary/{id}", deleteEntry(service)).Methods("DELETE").Name("deleteEntry")
}
