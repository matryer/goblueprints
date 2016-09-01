package main

import (
	"github.com/go-kit/kit/endpoint"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

// HasherService provides password hashing capabilities.
type HasherService interface {
	Hash(password string) (string, error)
	Valid(password, hash string) bool
}

type hasherService struct{}

func (hasherService) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (hasherService) Valid(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return false
	}
	return true
}

type hashRequest struct {
	P string `json:"p"`
}
type hashResponse struct {
	H   string `json:"h"`
	Err string `json:"err,omitempty"`
}
type validateRequest struct {
	P string `json:"p"`
	H string `json:"h"`
}
type validateResponse struct {
	V bool `json:"v"`
}

func makeHashEndpoint(srv HasherService) endpoint.Endpoint {
	return endpoint.Endpoint(func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(hashRequest)
		v, err := srv.Hash(req.P)
		if err != nil {
			return hashResponse{v, err.Error()}, nil
		}
		return hashResponse{v, ""}, nil
	})
}

func makeValidateEndpoint(srv HasherService) endpoint.Endpoint {
	return endpoint.Endpoint(func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(validateRequest)
		v := srv.Valid(req.P, req.H)
		return validateResponse{v}, nil
	})
}
