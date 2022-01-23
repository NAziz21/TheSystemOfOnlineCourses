package users

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/NAziz21/TheSystemOfOnlineCourses/helpers/errorsHelper"
	"github.com/NAziz21/TheSystemOfOnlineCourses/pkg/token"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

//Service - is a banner management service.
type Service struct {
	pool *pgxpool.Pool
}


//newService - create a service.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

type Registration struct {
	Name		string		`json:"name"`
	Phone		string		`json:"phone"`
	Login		string		`json:"login"`;	
	Password	string		`json:"password"`
}


type User struct {
	ID			int64		`json:"id"`
	Name 		string 		`json:"name"`
	Phone 		string 		`json:"phone"`
	Login 		string 		`json:"login"`
	Active 		bool		`json:"active"`
	Created 	time.Time 	`json:"created"`
}


type Auth struct {
	Login		string `json:"login"`
	Password	string `json:"password"`	
}



// Регистрация пользователя
func (s *Service) Register(ctx context.Context, registration *Registration) (*User, error) {
	var err error
	item := &User{}

	hash, err := bcrypt.GenerateFromPassword([]byte(registration.Password), bcrypt.DefaultCost)
	
	if err != nil {
		log.Print(err)
	}

	err = s.pool.QueryRow(ctx, 
		`INSERT INTO users(name, phone, login, password) 
		Values($1, $2, $3, $4)
		ON CONFLICT (phone) DO NOTHING RETURNING id, name, phone, login, active, created`, registration.Name, registration.Phone, registration.Login, hash).Scan(
			&item.ID, &item.Name, &item.Phone, &item.Login, &item.Active, &item.Created)

	
	if err == pgx.ErrNoRows {
		log.Print("Error In Method Part: pgx.ErrNoRows")
		return nil, errorsHelper.ErrNoRows
	}
	
	if err != nil {
		log.Print("Error In Method Part")
		return nil, errorsHelper.ErrInternal
	}

	return item, nil
}


// Проверка пароля и токен для пользователя
func (s *Service) Login(ctx context.Context, auth *Auth) (token string, err error) {
	var hash string
	var id int64
	var name string

	err = s.pool.QueryRow(ctx, `SELECT id, password, name from users WHERE login = $1`, auth.Login).Scan(&id, &hash, &name)

	if err == pgx.ErrNoRows {
		log.Print("No Such User")
		return "", errorsHelper.ErrNoSuchUser
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



type CourseSubscribeRequest struct {
	Course_id 		string 	`json:"course_id"`
}


type CourseSubscribeResponse struct {
	Message			string		`json:"Message"`
	ID 				int64 		`json:"ID"`
	CourseName 		string 		`json:"Course Name"`	
	Price			int64 		`json:"Price"`
	Created			time.Time	`json:"Data of subscription"`	
}


// Записаться на курс
func (s *Service) SubscribeToCourseSvc(ctx context.Context, courseID string, userID int64) (*CourseSubscribeResponse, error) {
	courseInfo := &CourseSubscribeResponse{}

	parseCourseID, _ := strconv.ParseInt(courseID, 10, 0)

	err := s.pool.QueryRow(ctx, `SELECT course_name, price FROM courses WHERE id = $1`, parseCourseID).Scan(&courseInfo.CourseName, &courseInfo.Price)
	

	if err == pgx.ErrNoRows {
		log.Print("Error In Method Part: Such course does not exist!")
		return nil, errorsHelper.ErrSubcribeCourse
	}

	if err != nil {
		log.Print("Method Part: Login")
		return nil, errorsHelper.ErrInternal
	}

	// Проверяем существует ли курс в БД (подписан ли уже)
	checkUnique, err := s.checkOfExistingCourse(ctx, parseCourseID, userID)
	if err != nil {
		log.Print("Method Part: Login")
		return nil, errorsHelper.ErrChechUnique
	}

	if checkUnique == true {
		log.Print("Method Part: CheckUnique")
		return nil, errorsHelper.ErrChechUnique
	}

	// Добавляем курс
	err = s.pool.QueryRow(ctx,`INSERT INTO user_courses(user_id, courses_id, course_name, price) Values($1, $2, $3, $4) RETURNING id, created`, 
	userID, parseCourseID, courseInfo.CourseName, courseInfo.Price).Scan(&courseInfo.ID, &courseInfo.Created)

	
	if err == pgx.ErrNoRows {
		log.Print("Error In Method Part: Such course does not exist!")
		return nil, errorsHelper.ErrNoRows
	}
	
	if err != nil {
		log.Print("Error In Method Part")
		return nil, errorsHelper.ErrInternal
	}

	courseInfo.Message = "Subscribed to course!"
	
	return courseInfo, nil
}

// Проверка на повторный subcribe курса
func (s *Service) checkOfExistingCourse(ctx context.Context, coursesID int64, userID int64) (bool, error) {
	var answer bool

	err := s.pool.QueryRow(ctx, `SELECT EXISTS (SELECT course_name FROM user_courses WHERE user_id = $1 AND courses_id = $2)`, userID, coursesID).Scan(&answer)

	if err != nil {
		return true, errorsHelper.ErrInternal
	}

	return answer, nil
}


type SubscribersInfo struct {
	Name 			string	`json:"name of subscriber"`
	CourseName		string	`json:"course name"`
	Active			bool	`json:"active subscribers"`
}


type UserCourses struct {
	Name 			string	`json:"name of user"`
	CourseName		string	`json:"course name"`
	Active			bool	`json:"active course"`
}


// Получение всех слушателей курса
func (s *Service) AllSubscribersFromCourseSvc(ctx context.Context, courseID int64) ([]*SubscribersInfo, error) {
	
	items := make([]*SubscribersInfo, 0)
	query := `SELECT c.course_name,  
			  (SELECT u.name from users u WHERE c.user_id = u.id),
			  c.active
			  FROM user_courses c 
			  WHERE courses_id = $1 AND active = 'TRUE';`
	
	rows, err := s.pool.Query(ctx, query, courseID)
	
	if err != nil {
		log.Print(err)
		return nil, errorsHelper.ErrInternal 
	}

	defer rows.Close()

	for rows.Next() {
		item := &SubscribersInfo{}
		err = rows.Scan(&item.CourseName, &item.Name, &item.Active)
		if err != nil {
			log.Print(err)
			return nil, err
		}		
		items = append(items, item)
	}	

	err = rows.Err()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return items, nil
}


// Получение всех курсов пользователя
func (s *Service) AllUsersCoursesSvc(ctx context.Context, userID int64) ([]*SubscribersInfo, error) {
	
	items := make([]*SubscribersInfo, 0)
	query := `SELECT c.course_name,  
			  (SELECT u.name from users u WHERE c.user_id = u.id),
			  c.active
			  FROM user_courses c 
			  WHERE user_id = $1 AND active = 'TRUE';`
	
	rows, err := s.pool.Query(ctx, query, userID)
	
	if err != nil {
		log.Print(err)
		return nil, errorsHelper.ErrInternal 
	}

	defer rows.Close()

	for rows.Next() {
		item := &SubscribersInfo{}
		err = rows.Scan(&item.CourseName, &item.Name, &item.Active)
		if err != nil {
			log.Print(err)
			return nil, err
		}		
		items = append(items, item)
	}	

	err = rows.Err()
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return items, nil
}







