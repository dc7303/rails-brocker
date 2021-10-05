package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/dc7303/rails-brocker/brocker"
)

type writeCodeBody struct {
	Code string `json:code`
}

type Server struct {
	brocker *brocker.Brocker
}

func New() *Server {
	dir := "/Users/scott/Workspace/dc7303/rails-brocker/rails-study/blog"
	return &Server{
		brocker: brocker.New(dir),
	}
}

func (s *Server) writeCode(w http.ResponseWriter, r *http.Request) {
	var wcb writeCodeBody
	err := json.NewDecoder(r.Body).Decode(&wcb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("[input code] ", wcb.Code)

	s.brocker.Write(fmt.Sprintf("%s\n", wcb.Code))
}

func (s *Server) HandleRequests() {
	log.Println("Run brocker")
	if err := s.brocker.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println("Run server :10000")

	router := mux.NewRouter()
	router.HandleFunc("/", s.writeCode).Methods("POST")

	srv := &http.Server{
		Addr:         "0.0.0.0:10000",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	if err = s.brocker.Close(); err != nil {
		log.Fatal(err)
	}
}
