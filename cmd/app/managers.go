package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/NAziz21/TheSystemOfOnlineCourses/cmd/app/middleware"
	"github.com/NAziz21/TheSystemOfOnlineCourses/helpers/validation"
	"github.com/NAziz21/TheSystemOfOnlineCourses/pkg/managers"
)

// Регистрация менеджера
func (s *Server) ManagerRegistration(writer http.ResponseWriter, request *http.Request) {
	
	var item *managers.ManagerRegistration
	
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := middleware.AuthenticationJWT(request.Context())
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}


	// Проверяем на Менеджер-Админ
	isAdmin, err := s.managerSvc.CheckIsAdmin(request.Context(), id)

	if err != nil {
		log.Printf("Manager is not an Admin! Error: %s", err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if !isAdmin {
		log.Printf("Manager is not an Admin! Error: %s", err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Валидация
	err = validation.ValidateRegister(item)

	if err != nil {
		log.Printf("Validation Failed! Error: %s", err)
		http.Error(writer, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	// Далее регистрируем менеджера 
	message, err := s.managerSvc.RegisterManager(request.Context(), item)
	
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(writer, message)
}


// Логин Менеджера
func (s *Server) ManagerLogin(writer http.ResponseWriter, request *http.Request) {
	
	var item *managers.AuthManager
	
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}	

	
	err = validation.ValidateLogin(item)

	if err != nil {
		log.Printf("Validation Failed! Error: %s", err)
		http.Error(writer, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}


	token, err := s.managerSvc.MLogin(request.Context(), item)
	
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(writer, token)
}


// Добавление курса
func (s *Server) AddCourse(writer http.ResponseWriter, request *http.Request) {
	
	name := &managers.CourseRequest{}

	err := json.NewDecoder(request.Body).Decode(&name)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Получаем ID с контекста
	managerID, err := middleware.AuthenticationJWT(request.Context())
	if err != nil {
		log.Printf("(Addcourse) Error in getting id from token: %e",err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}


	course, err := s.managerSvc.AddCourseSvc(request.Context(), name.Name, managerID)
	
	if err != nil {
		log.Printf("json problem: %s", err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respondJSON(writer, course)
}




// Регистрация пользователя
func (s *Server) AddAdmin(writer http.ResponseWriter, request *http.Request) {
	
	var item *managers.ManagerRegistration
	
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	message, err := s.managerSvc.RegisterManager(request.Context(), item)
	
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


// respondJSON - response from JSON.
func respondJSON(writer http.ResponseWriter, item interface{}) {

	data, err := json.Marshal(item)
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