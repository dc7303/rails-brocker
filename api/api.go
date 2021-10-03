package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dc7303/rails-brocker/pkg/brocker"
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
	s.brocker.Write("12 36 +p\n")
	fmt.Fprintf(w, "Welcom")
}

func (s *Server) HandleRequests() {
	log.Println("Run brocker")
	s.brocker.Run()
	defer s.brocker.Close()

	log.Println("Run server :10000")
	http.HandleFunc("/", s.calculate)
	err := http.ListenAndServe(":10000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
