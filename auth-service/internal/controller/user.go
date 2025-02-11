package controller

import (
	"context"
	"database/sql"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"newgolang/auth-service/internal/repository"
	"newgolang/auth-service/pkg/jwtc"
	"newgolang/auth-service/pkg/utils"
	"newgolang/proto/pb"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	log.Println("CreateUser: Received request")
	defer log.Println("CreateUser: Request processed")

	hashedPassword, _ := utils.HashPassword(req.Password)

	user := &pb.User{
		Name:     req.Name,
		Surname:  req.Surname,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := h.userRepo.Save(ctx, user); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	log.Println("UpdateUser: Received request")
	defer log.Println("UpdateUser: Request processed")

	user := &pb.User{
		Name:    req.Name,
		Surname: req.Surname,
	}

	if err := h.userRepo.UpdateByEmail(ctx, req.Email, user); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

func (h *UserHandler) GetUserByEmail(ctx context.Context, req *pb.GetUserRequestByEmail) (*pb.User, error) {
	log.Println("GetUserByEmail: Received request")
	defer log.Println("GetUserByEmail: Request processed")

	user, err := h.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return user, nil
}

func (h *UserHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	log.Println("ListUsers: Received request")
	defer log.Println("ListUsers: Request processed")

	users, err := h.userRepo.GetAll(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "users not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ListUsersResponse{Users: users}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.User, error) {
	log.Println("DeleteUser: Received request")
	defer log.Println("DeleteUser: Request processed")

	if err := h.userRepo.DeleteByID(ctx, req.Id); err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Println("Login: Received request")
	defer log.Println("Login: Request processed")

	user, err := h.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, _ := jwtc.GenerateToken(user.Email)

	return &pb.LoginResponse{Jwt: token}, nil
}

func (h *UserHandler) DecryptJwt(ctx context.Context, req *pb.DecryptJwtRequest) (*pb.DecryptJwtResponse, error) {
	log.Println("DecryptJwt: Received request")
	defer log.Println("DecryptJwt: Request processed")

	email, err := jwtc.ParseToken(req.Jwt)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DecryptJwtResponse{Email: email}, nil
}
