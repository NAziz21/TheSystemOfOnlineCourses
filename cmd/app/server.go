package app

import (
	"net/http"

	"github.com/NAziz21/TheSystemOfOnlineCourses/cmd/app/middleware"
	"github.com/NAziz21/TheSystemOfOnlineCourses/pkg/managers"
	"github.com/NAziz21/TheSystemOfOnlineCourses/pkg/users"
	"github.com/gorilla/mux"
)

const (
	GET = "GET"
	POST = "POST"
	DELETE = "DELETE"
)

type Server struct {
	mux 			*mux.Router
	userSvc			*users.Service
	managerSvc		*managers.Service
}


func NewServer(mux *mux.Router, userSvc *users.Service, managerSvc *managers.Service) *Server {
	return &Server{mux: mux, userSvc: userSvc, managerSvc: managerSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) Init() {
	s.mux.Use(middleware.Logger)
	s.mux.HandleFunc("/api/v1/register/users", s.UserRegistration).Methods(POST)
	s.mux.HandleFunc("/api/v1/login/users", s.Login).Methods(POST)
	
	customerAuthenticateMd := middleware.Auth()
	customersSubrouter := s.mux.PathPrefix("/auth/users").Subrouter()
	customersSubrouter.Use(customerAuthenticateMd)
	customersSubrouter.HandleFunc("/api/v1/subscribe", s.SubscribeToCourse).Methods(POST)
	customersSubrouter.HandleFunc("/api/v1/course/subscribers", s.AllSubscribersFromCourse).Methods(GET)
	customersSubrouter.HandleFunc("/api/v1/course/user", s.AllUsersCourses).Methods(GET)
	
	s.mux.HandleFunc("/managers/addADmin", s.AddAdmin).Methods(POST) // Для добавления главного Менеджера Админа, который будет добавлять других менеджеров
	s.mux.HandleFunc("/api/v1/login/managers", s.ManagerLogin).Methods(POST)
	managerMD := middleware.Auth()
	managersSubrouter := s.mux.PathPrefix("/auth/manager/").Subrouter()
	managersSubrouter.Use(managerMD)
	managersSubrouter.HandleFunc("/api/v1/register", s.ManagerRegistration).Methods(POST)
	managersSubrouter.HandleFunc("/api/v1/course/create", s.AddCourse).Methods(POST)
	
}


















