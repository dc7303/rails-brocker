package api

import (
	"log"
	"net/http"

	"github.com/dc7303/rails-brocker/brocker"
)

type Server struct {
	brocker *brocker.Brocker
}

func New() *Server {
	dir := "/Users/scott/Workspace/dc7303/rails-brocker/rails-study/blog"
	return &Server{
		brocker: brocker.New(dir),
	}
}

func (s *Server) calculate(w http.ResponseWriter, r *http.Request) {
	s.brocker.Write("Article\n")
}

func (s *Server) HandleRequests() {
	log.Println("Run brocker")
	if err := s.brocker.Run(); err != nil {
		log.Fatal(err)
	}

	log.Println("Run server :10000")
	http.HandleFunc("/", s.calculate)

	err := http.ListenAndServe(":10000", nil)
	if err != nil {
		log.Fatal(err)
	}

	if err = s.brocker.Close(); err != nil {
		log.Fatal(err)
	}
}
