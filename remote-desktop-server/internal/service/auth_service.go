// Package service реализует бизнес-логику приложения.
package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/common"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/config"
	"github.com/margar-melkonyan/tic-tac-toe-game/tic-tac-toe.git/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService предоставляет сервис для работы с аутентификацией пользователей.
type AuthService struct {
	repo repository.UserRepository
}

// NewAuthService создает новый экземпляр AuthService.
//
// Параметры:
//   - repoRoom: репозиторий для работы с пользователями
//
// Возвращает:
//   - *AuthService: указатель на созданный сервис
func NewAuthService(repoRoom repository.UserRepository) *AuthService {
	return &AuthService{
		repo: repoRoom,
	}
}

// Claims представляет структуру JWT токена с информацией о пользователе.
type Claims struct {
	Sub struct {
		Email string `json:"email"`
	} `json:"sub"`
	jwt.RegisteredClaims
}

// CheckTokenIsNotExpired проверяет валидность JWT токена.
//
// Параметры:
//   - token: JWT токен (может содержать префикс "Bearer")
//
// Возвращает:
//   - *Claims: данные из токена
//   - error: ошибку, если токен невалиден или просрочен
func CheckTokenIsNotExpired(token string) (*Claims, error) {
	token = strings.TrimSpace(strings.ReplaceAll(token, "Bearer ", ""))

	claims, err := parseToken(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

// SignIn выполняет аутентификацию пользователя.
//
// Параметры:
//   - ctx: контекст
//   - form: данные для входа (email и пароль)
//
// Возвращает:
//   - map[string]string: содержит JWT токен
//   - error: ошибки:
//   - "password is not valid" - неверный пароль
//   - "user not found" - пользователь не найден
//   - ошибки генерации токена
func (service *AuthService) SignIn(ctx context.Context, form common.AuthSignInRequest) (map[string]string, error) {
	currentUser, err := service.repo.FindByEmail(ctx, form.Email)
	if err != nil {
		return nil, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(strings.TrimSpace(form.Password))); err != nil {
		return nil, errors.New("password is not valid")
	}

	accessToken, err := getToken(*currentUser)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"token": accessToken,
	}, nil
}

// SignUp регистрирует нового пользователя.
//
// Параметры:
//   - ctx: контекст
//   - form: данные для регистрации
//
// Возвращает:
//   - error: ошибки:
//   - "user with this email already exists" - пользователь уже существует
//   - ошибки хеширования пароля
//   - ошибки создания пользователя
func (service *AuthService) SignUp(ctx context.Context, form common.AuthSignUpRequest) error {
	if _, err := service.repo.FindByEmail(ctx, form.Email); err == nil {
		return errors.New("user with this email already exists")
	}
	password, err := bcrypt.GenerateFromPassword(
		[]byte(strings.TrimSpace(form.Password)),
		config.ServerConfig.BcryptPower,
	)
	if err != nil {
		return err
	}
	form.Password = string(password)
	return service.repo.Create(ctx, form)
}

// getToken генерирует JWT токен для пользователя.
//
// Параметры:
//   - user: данные пользователя
//
// Возвращает:
//   - string: JWT токен
//   - error: ошибки генерации токена
func getToken(user common.User) (string, error) {
	seconds := config.ServerConfig.JWTConfig.AccessTokenTTL
	duration, err := time.ParseDuration(seconds)

	if err != nil {
		return "", err
	}
	payload := jwt.MapClaims{
		"sub": map[string]interface{}{
			"email": user.Email,
		},
		"exp": time.Now().Add(time.Duration(duration)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	jwtSecret := []byte(config.ServerConfig.JWTConfig.AccessTokenSecret)
	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return t, nil
}

// parseToken разбирает и валидирует JWT токен.
//
// Параметры:
//   - token: JWT токен
//
// Возвращает:
//   - *Claims: данные из токена
//   - error: ошибки:
//   - "token is expired" - токен просрочен
//   - "your token is invalid" - невалидный токен
//   - ошибки парсинга токена
func parseToken(token string) (*Claims, error) {
	var claims Claims
	t, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.ServerConfig.JWTConfig.AccessTokenSecret), nil
	})

	if claims.ExpiresAt != nil && time.Now().Unix() > claims.ExpiresAt.Unix() {
		return nil, errors.New("token is expired")
	}

	if err != nil || !t.Valid {
		return nil, errors.New("your token is invalid")
	}

	return &claims, nil
}
