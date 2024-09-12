package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"unicode/utf8"

	pb "github.com/kTowkA/GophKeeper/grpc"
	"github.com/kTowkA/GophKeeper/internal/model"
	"github.com/kTowkA/GophKeeper/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	validate_password_min_length = 10
	validate_password_lower_case = true
	validate_password_upper_case = true
	validate_password_numbers    = true
	validate_password_symbols    = true
)

var (
	symbols                      = []string{"!", "@", "#", "$", "%", "\\^", "&", "?", "*", "\\(", "\\)", "\\[", "\\]", "_", "\\-", "+", "=", "|", "\\."}
	ErrValidatePasswordMinLength = fmt.Errorf("пароль должен быть не менее %d символов длиной", validate_password_min_length)
	ErrValidatePasswordLowerCase = errors.New("пароль должен содержать символы в нижнем регистре")
	ErrValidatePasswordUpperCase = errors.New("пароль должен содержать символы в верхнем регистре")
	ErrValidatePasswordNumbers   = errors.New("пароль должен содержать числа")
	ErrValidatePasswordSymbols   = fmt.Errorf("пароль должен содержать специальные символы \"%s\"", strings.Join(symbols, ""))
)

func (s *Server) Register(ctx context.Context, r *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := validatePassword(r.Password)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("получение хеша пароля", slog.String("логин", r.Login), slog.String("ошибка", err.Error()))
		return nil, err
	}
	_, err = s.db.Register(
		ctx,
		model.StorageRegisterRequest{
			Login:    r.Login,
			Password: string(hash),
		},
	)
	switch {
	case errors.Is(err, storage.ErrLoginIsAlreadyOccupied):
		s.log.Debug("попытка регистрации с существующим логином", slog.String("логин", r.Login))
		return nil, status.Error(codes.AlreadyExists, err.Error())
	case err != nil:
		s.log.Error("регистрация пользователя", slog.String("логин", r.Login), slog.String("ошибка", err.Error()))
		return nil, err
	}
	return &pb.RegisterResponse{
		RegisterStatus:  true,
		RegisterMessage: "ok",
	}, nil
}

func validatePassword(password string) error {
	errs := make([]error, 0)
	if utf8.RuneCountInString(password) < validate_password_min_length {
		errs = append(errs, ErrValidatePasswordMinLength)
	}
	if validate_password_lower_case {
		rg := regexp.MustCompile(`[a-z]`)
		if rg.FindString(password) == "" {
			errs = append(errs, ErrValidatePasswordLowerCase)
		}
	}
	if validate_password_upper_case {
		rg := regexp.MustCompile(`[A-Z]`)
		if rg.FindString(password) == "" {
			errs = append(errs, ErrValidatePasswordUpperCase)
		}
	}
	if validate_password_numbers {
		rg := regexp.MustCompile(`[0-9]`)
		if rg.FindString(password) == "" {
			errs = append(errs, ErrValidatePasswordNumbers)
		}
	}
	if validate_password_symbols {
		rg := regexp.MustCompile(`[` + strings.Join(symbols, "") + `]`)
		if rg.FindString(password) == "" {
			errs = append(errs, ErrValidatePasswordSymbols)
		}
	}
	return errors.Join(errs...)
}
