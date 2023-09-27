package gapi

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/corey888773/golang-course/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
)

func (server *Server) authorizeUser(ctx *context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(*ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	values := md.Get(authorizationHeaderKey)
	if len(values) == 0 {
		return nil, errors.New("missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Split(authHeader, " ")

	if len(fields) < 2 {
		return nil, errors.New("invalid authorization header")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationTypeBearer {
		return nil, errors.New("unsuported authorization type")
	}

	accesToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accesToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token %s", err)
	}

	return payload, nil
}
