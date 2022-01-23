package errorsHelper


import "errors"



var ErrNotFound = errors.New("item not found")
var ErrTableIsEmpty = errors.New("table is empty")
var ErrInternal = errors.New("internal error")
var ErrNoRows = errors.New("no rows were returned")


var ErrInvalidPassword = errors.New("invalid password")


var ErrNoSuchUser = errors.New("no such user")
var ErrNoSuchManager = errors.New("no such manager")




var ErrNoAdmin = errors.New("not an admin")

// Validation
var ErrValidationManager = errors.New("validation failed")

// Validation: Login
var ErrValidationManagerLengthLogin = errors.New("length of Login is too short or empty")

var ErrValidationManagerLengthPassword = errors.New("length of Password is too short or empty")

// Validation: Register
var ErrValidationManagerRegLengthName = errors.New("length of Name is to short or empty")
var ErrValidationManagerRegLengthPhone = errors.New("length of Phone is not to be short less 12")
var ErrValidationManagerRegLengthLogin = errors.New("length of Login is not to be short less 4 and long 30 or empty")
var ErrValidationManagerRegLengthPassword = errors.New("length of Password is too short or empty")


// AddCourse 
var ErrAddCourseInternal = errors.New("")

// Subscribe Course
var ErrSubcribeCourse = errors.New("no such course")
var ErrChechUnique = errors.New("already subcribed")
