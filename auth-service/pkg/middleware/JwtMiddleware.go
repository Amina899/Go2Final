package middleware

import (
	"auth-service/auth-service/pkg/jwtc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

func JWTMiddleware(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	authorization, ok := md["authorization"]
	if !ok || len(authorization) == 0 {
		return "", status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	authHeader := authorization[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", status.Errorf(codes.InvalidArgument, "invalid authorization header format")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", status.Errorf(codes.Unauthenticated, "token is not provided")
	}

	email, err := jwtc.ParseToken(token)
	if err != nil {
		return "", status.Error(codes.Unauthenticated, err.Error())
	}

	return email, nil
}
