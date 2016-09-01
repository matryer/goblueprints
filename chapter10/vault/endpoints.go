package vault

import (
	"errors"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

// Endpoints represents all endpoints for the vault Service.
type Endpoints struct {
	HashEndpoint     endpoint.Endpoint
	ValidateEndpoint endpoint.Endpoint
}

// Hash uses the HashEndpoint to hash a password.
func (e Endpoints) Hash(ctx context.Context, password string) (string, error) {
	req := hashRequest{Password: password}
	resp, err := e.HashEndpoint(ctx, req)
	if err != nil {
		return "", err
	}
	hashResp := resp.(hashResponse)
	if hashResp.Err != "" {
		return "", errors.New(hashResp.Err)
	}
	return hashResp.Hash, nil
}

// Validate uses the ValidateEndpoint to validate a password and hash pair.
func (e Endpoints) Validate(ctx context.Context, password, hash string) (bool, error) {
	req := validateRequest{Password: password, Hash: hash}
	resp, err := e.ValidateEndpoint(ctx, req)
	if err != nil {
		return false, err
	}
	hashResp := resp.(validateResponse)
	return hashResp.Valid, nil
}
