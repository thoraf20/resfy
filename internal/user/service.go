package user

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *AuthService {
	db.AutoMigrate(&User{})
	return &AuthService{db: db}
}

func (s *AuthService) Register(email, password string) (User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return User{}, err
	}

	user := User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: string(hashed),
	}

	if err := s.db.Create(&user).Error; 
	err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *AuthService) Authenticate(email, password string) (User, error) {
	var user User
	if err := s.db.First(&user, "email = ?", email).Error; 
	err != nil {
		return User{}, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); 
	err != nil {
		return User{}, errors.New("invalid credentials")
	}

	return user, nil
}
