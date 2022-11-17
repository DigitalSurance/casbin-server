package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/casbin/casbin-server/proto"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type JsonServer struct {
	casbinServer *Server
	router       *chi.Mux
}

func NewJsonServer(casbinServer *Server) *JsonServer {
	router := chi.NewRouter()
	server := &JsonServer{
		casbinServer: casbinServer,
		router:       router,
	}

	// middlewares
	server.router.Use(middleware.Logger)
	server.router.Use(middleware.AllowContentType("application/json"))

	// routes
	server.attachRoutes()

	return server
}

func (s *JsonServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(rw, r)
}

type RelationTuple struct {
	Namespace string `json:"namespace,omitempty"`
	Object    string `json:"object"`
	Relation  string `json:"relation"`
	SubjectId string `json:"subject_id"`
	// SubjectSet SubjectSet 	`json:"subject_set"`
}

func (s *JsonServer) attachRoutes() {
	s.router.Post("/relation-tuple/check", s.handleRelationTupleCheck)
}

func (s *JsonServer) handleRelationTupleCheck(rw http.ResponseWriter, r *http.Request) {
	var relationTuple RelationTuple
	if err := json.NewDecoder(r.Body).Decode(&relationTuple); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	in := proto.EnforceRequest{
		EnforcerHandler: 0,
		Params:          []string{relationTuple.SubjectId, relationTuple.Object, relationTuple.Relation},
	}
	ok, err := s.casbinServer.Enforce(r.Context(), &in)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(map[string]string{"message": "internal error (check logs)"})
		return
	}
	if !ok.Res {
		rw.WriteHeader(http.StatusForbidden)
		json.NewEncoder(rw).Encode(map[string]string{"message": "forbidden"})
		return
	}

	// we gucci
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(map[string]string{"message": "ok"})
}
