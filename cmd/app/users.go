package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/NAziz21/TheSystemOfOnlineCourses/cmd/app/middleware"
	"github.com/NAziz21/TheSystemOfOnlineCourses/pkg/users"
)

// Регистрация пользователя
func (s *Server) UserRegistration(writer http.ResponseWriter, request *http.Request) {
	
	var item *users.Registration
	
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	message, err := s.userSvc.Register(request.Context(), item)
	
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println(err) 
	}
}


// Получение запроса от пользователя для получения токена (аутентификация)
func (s *Server) Login(writer http.ResponseWriter, request *http.Request) {
	
	var item *users.Auth
	
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	token, err := s.userSvc.Login(request.Context(), item)
	
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(token)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println(err) 
	}
}


// Запись на курс
func (s *Server) SubscribeToCourse(writer http.ResponseWriter, request *http.Request) {
	
	var courseID *users.CourseSubscribeRequest

	err := json.NewDecoder(request.Body).Decode(&courseID)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Получаем UserID пользователя с контекста
	userID, err := middleware.AuthenticationJWT(request.Context())
	if err != nil {
		log.Printf("Error in getting id from token: %e",err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	response, err := s.userSvc.SubscribeToCourseSvc(request.Context(), courseID.Course_id, userID)
	
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(writer, response)
}


// Получение всех слушателей курса
func (s *Server) AllSubscribersFromCourse(writer http.ResponseWriter, request *http.Request) {
		
	idParam := request.URL.Query().Get("course_id")

	courseID, err := strconv.ParseInt(idParam, 10, 64)
	
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	response, err := s.userSvc.AllSubscribersFromCourseSvc(request.Context(), courseID)
	
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(writer, response)
}


// Получение всех курсов пользователя
func (s *Server) AllUsersCourses(writer http.ResponseWriter, request *http.Request) {
		
	idParam := request.URL.Query().Get("user_id")

	userID, err := strconv.ParseInt(idParam, 10, 64)
	
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	response, err := s.userSvc.AllUsersCoursesSvc(request.Context(), userID)
	
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(writer, response)
}

