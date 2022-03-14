package usecase

import (
	"errors"
	"go-midtrans/domain"
	"os"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepository domain.AuthRepository
}

func NewUserRepository(ur domain.AuthRepository) *userUsecase {
	return &userUsecase{userRepository: ur}
}

func (u *userUsecase) Register(name string, email string, password string) (domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return domain.User{}, domain.ErrInternalServer
	}

	existedUser, _ := u.userRepository.GetByEmail(email)
	if existedUser.ID != 0 {
		return domain.User{}, domain.ErrEmailConflict
	}

	user := domain.User{
		Name:     name,
		Email:    email,
		Password: string(hash),
	}

	newUser, err := u.userRepository.Create(user)
	if err != nil {
		return newUser, domain.ErrInternalServer
	}

	return newUser, nil
}

func (u *userUsecase) Login(email string, password string) (domain.User, error) {
	user, err := u.userRepository.GetByEmail(email)
	if err != nil {
		return domain.User{}, domain.ErrInternalServer
	}

	if user.ID == 0 {
		return domain.User{}, domain.ErrNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return domain.User{}, domain.ErrPassNotMatch
	}

	return user, nil
}

func (u *userUsecase) ValidateToken(encodedToken string) (domain.User, error) {
	token, err := jwt.Parse(encodedToken, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("token invalid")
		}
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return domain.User{}, err
	}

	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return domain.User{}, err
	}

	user, err := u.userRepository.GetByEmail(claim["email"].(string))
	if err != nil {
		return domain.User{}, err
	}

	if user.ID == 0 {
		return domain.User{}, domain.ErrNotFound
	}

	return user, nil
}
