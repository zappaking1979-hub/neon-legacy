package application

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	"github.com/neonlegacy/server/internal/domain/player"
)

var (
	ErrEmailTaken       = errors.New("email already in use")
	ErrUsernameTaken    = errors.New("username already in use")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrInvalidUsername  = errors.New("username must be 3-20 characters, letters and numbers only")
	ErrInvalidPassword  = errors.New("password must be at least 8 characters")
	ErrInvalidGender    = errors.New("gender must be male, female, or other")
	ErrInvalidCreds    = errors.New("invalid email or password")
	ErrNotAuthenticated = errors.New("not authenticated")
)

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
)

type Session struct {
	PlayerID  uuid.UUID
	Token     string
	ExpiresAt time.Time
}

type AuthService struct {
	playerRepo player.Repository
	redis      *redis.Client
	sessionTTL time.Duration
	bcryptCost int
}

func NewAuthService(pr player.Repository, rdb *redis.Client, sessionTTL time.Duration, bcryptCost int) *AuthService {
	return &AuthService{
		playerRepo: pr,
		redis:      rdb,
		sessionTTL: sessionTTL,
		bcryptCost: bcryptCost,
	}
}

func (s *AuthService) Register(ctx context.Context, email, username, password string, gender player.Gender) (*player.Player, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	username = strings.TrimSpace(username)

	if !emailRegex.MatchString(email) {
		return nil, ErrInvalidEmail
	}
	if !usernameRegex.MatchString(username) {
		return nil, ErrInvalidUsername
	}
	if len(password) < 8 {
		return nil, ErrInvalidPassword
	}
	switch gender {
	case player.GenderMale, player.GenderFemale, player.GenderOther:
	default:
		return nil, ErrInvalidGender
	}

	existing, err := s.playerRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("check email: %w", err)
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}

	existing, err = s.playerRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("check username: %w", err)
	}
	if existing != nil {
		return nil, ErrUsernameTaken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	p := player.NewPlayer(email, username, string(hash), gender)
	if err := s.playerRepo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("create player: %w", err)
	}

	return p, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*Session, *player.Player, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	p, err := s.playerRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, fmt.Errorf("lookup player: %w", err)
	}
	if p == nil {
		return nil, nil, ErrInvalidCreds
	}

	if err := bcrypt.CompareHashAndPassword([]byte(p.PasswordHash), []byte(password)); err != nil {
		return nil, nil, ErrInvalidCreds
	}

	session, err := s.createSession(ctx, p.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("create session: %w", err)
	}

	return session, p, nil
}

func (s *AuthService) ValidateSession(ctx context.Context, token string) (*player.Player, error) {
	key := "session:" + token
	playerIDStr, err := s.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, ErrNotAuthenticated
	}
	if err != nil {
		return nil, fmt.Errorf("redis get: %w", err)
	}

	if err := s.redis.Expire(ctx, key, s.sessionTTL).Err(); err != nil {
		return nil, fmt.Errorf("redis expire: %w", err)
	}

	playerID, err := uuid.Parse(playerIDStr)
	if err != nil {
		return nil, fmt.Errorf("parse player id: %w", err)
	}

	p, err := s.playerRepo.GetByID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("get player: %w", err)
	}
	if p == nil {
		return nil, ErrNotAuthenticated
	}

	return p, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.redis.Del(ctx, "session:"+token).Err()
}

func (s *AuthService) createSession(ctx context.Context, playerID uuid.UUID) (*Session, error) {
	token := uuid.New().String()
	now := time.Now()

	key := "session:" + token
	if err := s.redis.Set(ctx, key, playerID.String(), s.sessionTTL).Err(); err != nil {
		return nil, fmt.Errorf("redis set: %w", err)
	}

	return &Session{
		PlayerID:  playerID,
		Token:     token,
		ExpiresAt: now.Add(s.sessionTTL),
	}, nil
}
