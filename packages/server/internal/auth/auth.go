package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

type AuthService struct {
	userRepo    UserRepository
	tokenSvc    TokenService
	passwordSvc PasswordService
}

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash []byte    `json:"-"`
	Salt         []byte    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	LastLoginAt  time.Time `json:"lastLoginAt"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token     string `json:"token"`
	User      *User  `json:"user"`
	ExpiresAt int64  `json:"expiresAt"`
}

type Claims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type UserRepository interface {
	CreateUser(username, email string, passwordHash, salt []byte) (*User, error)
	GetUserByUsername(username string) (*User, error)
	UpdateLastLogin(userID string) error
	UserExists(username string) bool
}

type TokenService interface {
	GenerateToken(user *User) (string, int64, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type PasswordService interface {
	HashPassword(password string) (passwordHash, salt []byte, err error)
	VerifyPassword(password string, passwordHash, salt []byte) bool
}

func NewAuthService(userRepo UserRepository, tokenSvc TokenService, passwordSvc PasswordService) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		tokenSvc:    tokenSvc,
		passwordSvc: passwordSvc,
	}
}

func (a *AuthService) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := a.validateRegistration(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if a.userRepo.UserExists(req.Username) {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	passwordHash, salt, err := a.passwordSvc.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	user, err := a.userRepo.CreateUser(req.Username, req.Email, passwordHash, salt)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	token, expiresAt, err := a.tokenSvc.GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (a *AuthService) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := a.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !a.passwordSvc.VerifyPassword(req.Password, user.PasswordHash, user.Salt) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	a.userRepo.UpdateLastLogin(user.ID)

	token, expiresAt, err := a.tokenSvc.GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := AuthResponse{
		Token:     token,
		User:      user,
		ExpiresAt: expiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (a *AuthService) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Bearer token required", http.StatusUnauthorized)
			return
		}

		claims, err := a.tokenSvc.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-ID", claims.UserID)
		r.Header.Set("X-Username", claims.Username)

		next(w, r)
	}
}

func (a *AuthService) validateRegistration(req RegisterRequest) error {
	if len(req.Username) < 3 || len(req.Username) > 20 {
		return fmt.Errorf("username must be between 3 and 20 characters")
	}
	if len(req.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters")
	}
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}
	return nil
}

type InMemoryUserRepository struct {
	users map[string]*User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*User),
	}
}

func (r *InMemoryUserRepository) CreateUser(username, email string, passwordHash, salt []byte) (*User, error) {
	user := &User{
		ID:           generateUserID(),
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Salt:         salt,
		CreatedAt:    time.Now(),
		LastLoginAt:  time.Now(),
	}
	r.users[username] = user
	return user, nil
}

func (r *InMemoryUserRepository) GetUserByUsername(username string) (*User, error) {
	user, exists := r.users[username]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (r *InMemoryUserRepository) UpdateLastLogin(userID string) error {
	for _, user := range r.users {
		if user.ID == userID {
			user.LastLoginAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("user not found")
}

func (r *InMemoryUserRepository) UserExists(username string) bool {
	_, exists := r.users[username]
	return exists
}

type JWTTokenService struct {
	jwtSecret []byte
}

func NewJWTTokenService(jwtSecret string) *JWTTokenService {
	return &JWTTokenService{
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *JWTTokenService) GenerateToken(user *User) (string, int64, error) {
	expiresAt := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", 0, err
	}

	return tokenString, expiresAt.Unix(), nil
}

func (s *JWTTokenService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

type Argon2PasswordService struct {
	saltLength  uint32
	iterations  uint32
	memory      uint32
	parallelism uint8
	keyLength   uint32
}

func NewArgon2PasswordService() *Argon2PasswordService {
	return &Argon2PasswordService{
		saltLength:  16,
		iterations:  1,
		memory:      64 * 1024,
		parallelism: 4,
		keyLength:   32,
	}
}

func (s *Argon2PasswordService) HashPassword(password string) (passwordHash, salt []byte, err error) {
	salt = make([]byte, s.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, err
	}

	passwordHash = argon2.IDKey([]byte(password), salt, s.iterations, s.memory, s.parallelism, s.keyLength)
	return passwordHash, salt, nil
}

func (s *Argon2PasswordService) VerifyPassword(password string, passwordHash, salt []byte) bool {
	computedHash := argon2.IDKey([]byte(password), salt, s.iterations, s.memory, s.parallelism, s.keyLength)
	return subtle.ConstantTimeCompare(passwordHash, computedHash) == 1
}

func generateUserID() string {
	return fmt.Sprintf("user_%d", time.Now().UnixNano())
}
