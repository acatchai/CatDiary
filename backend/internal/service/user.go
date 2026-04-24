package service

import (
	"errors"
	"strings"

	"github.com/acatchai/catdiary/backend/internal/model"
	"github.com/acatchai/catdiary/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailExists     = errors.New("邮箱已存在")
	ErrInvalidUsername = errors.New("用户名不能为空")
	ErrInvalidEmail    = errors.New("邮箱不能为空")
)

func UpdateMe(userID uint, username, email, avatar *string) (*model.User, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	updates := make(map[string]any)

	if username != nil {
		newUsername := strings.TrimSpace(*username)
		if newUsername == "" {
			return nil, ErrInvalidUsername
		}
		if newUsername != user.Username {
			existing, err := repository.GetUserByUsername(newUsername)
			if err != nil {
				return nil, err
			}
			if existing != nil && existing.ID != userID {
				return nil, ErrUserExists
			}
			updates["username"] = newUsername
		}
	}
	if email != nil {
		newEmail := strings.TrimSpace(*email)
		if newEmail == "" {
			return nil, ErrInvalidEmail
		}
		if newEmail != user.Email {
			existing, err := repository.GetUserByEmail(newEmail)
			if err != nil {
				return nil, err
			}
			if existing != nil && existing.ID != userID {
				return nil, ErrEmailExists
			}
			updates["email"] = newEmail
		}
	}

	if avatar != nil {
		updates["avatar"] = strings.TrimSpace(*avatar)
	}

	if err := repository.UpdateUserByID(userID, updates); err != nil {
		return nil, err
	}
	return repository.GetUserByID(userID)
}

func ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrWrongPassword
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return repository.UpdateUserPasswordHash(userID, string(hashedPassword))
}
