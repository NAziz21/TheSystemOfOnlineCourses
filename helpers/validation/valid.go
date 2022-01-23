package validation

import (
	"github.com/NAziz21/TheSystemOfOnlineCourses/helpers/errorsHelper"
	"github.com/NAziz21/TheSystemOfOnlineCourses/pkg/managers"
)

//Валидация полученных данных Manager

// Login
func ValidateLogin(manager *managers.AuthManager) (err error) {
	
	if len (manager.Login) < 4 || manager.Login == "" {
		return errorsHelper.ErrValidationManagerLengthLogin
	}

	
	if len (manager.Password) < 4 || manager.Password == "" {
		return errorsHelper.ErrValidationManagerLengthPassword
	}

	return  nil

}

// Register
func ValidateRegister(manager *managers.ManagerRegistration) (err error) {
	
	if len (manager.Name) <= 2|| manager.Name == "" {
		return errorsHelper.ErrValidationManagerRegLengthName
	}

	if len (manager.Login) < 4 || manager.Login == "" {
		return errorsHelper.ErrValidationManagerLengthLogin
	}

	if len (manager.Phone) != 13 {
		return errorsHelper.ErrValidationManagerRegLengthPhone
	}

	if len (manager.Password) < 4 || manager.Password == "" {
		return errorsHelper.ErrValidationManagerLengthPassword
	}

	return  nil

}