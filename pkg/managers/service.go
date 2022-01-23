package managers

import (
	"context"
	"log"

	"github.com/NAziz21/TheSystemOfOnlineCourses/helpers/errorsHelper"
	tokenHelpers "github.com/NAziz21/TheSystemOfOnlineCourses/pkg/token"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)


type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}


type Manager struct {
	ID		int64 	`json:"id"`
	Name	string	`json:"name"`
	Phone	string	`json:"phone"`
	Login	string	`json:"login"`
	Active	string	`json:"active"`  
}

type ManagerRegistration struct{
	Name		string	`json:"name"`
	Phone		string	`json:"phone"`
	Login		string	`json:"login"`
	Password	string	`json:"password"`
}

type Message struct{
	Name	string	`json:"name"`
	Phone	string	`json:"phone"`
	Login	string	`json:"login"`
}

type AuthManager struct {
	Login 		string	`json:"login"`
	Password	string	`json:"password"`
}


// Регистрация пользователя
func (s *Service) RegisterManager(ctx context.Context, registration *ManagerRegistration) (*Message, error) {
	var err error
	item := &Message{}

	hash, err := bcrypt.GenerateFromPassword([]byte(registration.Password), bcrypt.DefaultCost)
	
	if err != nil {
		log.Print(err)
	}

	err = s.pool.QueryRow(ctx, 
		`INSERT INTO managers(name, phone, login, password) 
		Values($1, $2, $3, $4)
		ON CONFLICT (phone) DO NOTHING RETURNING name, phone, login`, registration.Name, registration.Phone, registration.Login, hash).Scan(
			&item.Name, &item.Phone, &item.Login)

		
	if err == pgx.ErrNoRows {
		log.Print("Error In Method Part (Register Manager): pgx.ErrNoRows")
		return nil, errorsHelper.ErrNoRows
	}
	
	if err != nil {
		log.Print("Error In Method Part")
		return nil, errorsHelper.ErrInternal
	}

	return item, nil
}


// Проверка пароля и токен для пользователя
func (s *Service) MLogin(ctx context.Context, auth *AuthManager) (token string, err error) {
	var hash string
	var id int64
	var name string
	

	err = s.pool.QueryRow(ctx, `SELECT id, password, name from managers WHERE login = $1`, auth.Login).Scan(&id, &hash, &name)
	
	if err == pgx.ErrNoRows {
		log.Print("No Such Manager")
		return "", errorsHelper.ErrNoSuchManager
	}

	if err != nil {
		log.Print("Method Part: Login")
		return "", errorsHelper.ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(auth.Password))
	if err != nil {
		log.Print("Method Part: bcrypt Password")
		return "", errorsHelper.ErrInvalidPassword
	}
	
	token, err = tokenHelpers.GenerateToken(name, id) 
	if err != nil {
		log.Print("Method Part: token wrong")
		return "", tokenHelpers.ErrInvalidToken
	}
	
	return token, nil
}


func (s *Service) CheckIsAdmin(ctx context.Context, id int64) (answer bool, err error) {
	
	err = s.pool.QueryRow(ctx, `SELECT is_admin from managers WHERE id = $1`, id).Scan(&answer)

	if err == pgx.ErrNoRows {
		log.Print("Method Manager Part (Check Is Admin): not Admin")
		return false, errorsHelper.ErrNoAdmin
	}

	if err != nil {
		log.Print("Method Manager Part (Check Is Admin): Internal Error")
		return false, errorsHelper.ErrInternal
	}

	return answer, nil
}

type CourseRequest struct {
	Name 		string	`json:"name"`
}


type CourseResponse struct {
	Message				string	`json:"message"`
	NameOfCourse 		string	`json:"name of course"`
}


func (s *Service) AddCourseSvc(ctx context.Context, name string, id int64) (*CourseResponse, error) {
	
	message := &CourseResponse{}

	err := s.pool.QueryRow(ctx, 
		`INSERT INTO courses(manager_id, course_name) 
		Values($1, $2)
		ON CONFLICT (course_name) DO NOTHING RETURNING course_name`, id, name).Scan(
			&message.NameOfCourse)

		
	if err == pgx.ErrNoRows {
		log.Print("Error In Method Part (AddCourseSvc): pgx.ErrNoRows")
		return nil, errorsHelper.ErrNoRows
	}
	
	if err != nil {
		log.Print("Error In Method Part(AddCourseSvc)")
		return nil, errorsHelper.ErrInternal
	}	

	message.Message = "Course was successfuly added!"
	return message, nil
}















