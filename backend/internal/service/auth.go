package service

import (
	"errors"

	"github.com/acatchai/catdiary/backend/internal/model"
	"github.com/acatchai/catdiary/backend/internal/repository"
	"github.com/acatchai/catdiary/backend/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists    = errors.New("用户名已存在")
	ErrUserNotFound  = errors.New("用户不存在")
	ErrWrongPassword = errors.New("密码错误")
)

// Register 注册
func Register(username, password, email string) (*model.User, error) {
	// 1. 检查用户名是否冲突
	existingUser, err := repository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// 2， 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Email:        email,
	}

	// 3. 落库
	err = repository.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Login 登录
func Login(username, password string) (string, *model.User, error) {
	// 1. 查找用户
	user, err := repository.GetUserByUsername(username)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, ErrUserNotFound
	}
	// 2. 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", nil, ErrWrongPassword
	}
	// 3. 生成 Token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func GetMe(userID uint) (*model.User, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
