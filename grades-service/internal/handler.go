package internal

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"newgolang/proto/pb"
)

type GradesHandler struct {
	pb.UnimplementedGradeServiceServer
	gradeRepo GradeRepository
}

func NewGradeHandler(gradeRepo GradeRepository) *GradesHandler {
	return &GradesHandler{
		gradeRepo: gradeRepo,
	}
}

func (h *GradesHandler) CreateGrade(ctx context.Context, req *pb.CreateGradeRequest) (*pb.Grade, error) {
	log.Println("CreateGrade: Received request")
	defer log.Println("CreateGrade: Request processed")

	userConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer userConn.Close()

	userClient := pb.NewUserServiceClient(userConn)

	// Call the DecryptJwt method on the gRPC client
	decryptReq := &pb.DecryptJwtRequest{Jwt: req.Jwt}
	decryptRes, err := userClient.DecryptJwt(ctx, decryptReq)
	if err != nil {
		return nil, err
	}

	userEmail := decryptRes.Email

	userDataReq := &pb.GetUserRequestByEmail{Email: userEmail}
	userDataRes, err := userClient.GetUserByEmail(ctx, userDataReq)

	if userDataRes.Role != "TEACHER" {
		log.Printf("CreateGrade: User does not have permission: %v", userEmail)
		return nil, errors.New("you don't have permisson")
	}

	createdGrade, err := h.gradeRepo.CreateGrade(ctx, req)
	if err != nil {
		return nil, err
	}

	log.Println("CreateGrade: Grade created successfully")
	return createdGrade, nil
}

func (h *GradesHandler) UpdateGrade(ctx context.Context, req *pb.UpdateGradeRequest) (*pb.Grade, error) {
	log.Println("UpdateGrade: Received request")
	defer log.Println("UpdateGrade: Request processed")

	userConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer userConn.Close()

	userClient := pb.NewUserServiceClient(userConn)

	// Call the DecryptJwt method on the gRPC client
	decryptReq := &pb.DecryptJwtRequest{Jwt: req.Jwt}
	decryptRes, err := userClient.DecryptJwt(ctx, decryptReq)
	if err != nil {
		return nil, err
	}

	userEmail := decryptRes.Email

	userDataReq := &pb.GetUserRequestByEmail{Email: userEmail}
	userDataRes, err := userClient.GetUserByEmail(ctx, userDataReq)

	if userDataRes.Role != "TEACHER" {
		log.Printf("You don't have permisson!: %v", err)
		return nil, errors.New("you don't have permisson")
	}

	updatedGrade, err := h.gradeRepo.UpdateGrade(ctx, req)
	if err != nil {
		return nil, err
	}
	return updatedGrade, nil
}

func (h *GradesHandler) GetGrade(ctx context.Context, req *pb.GetGradeRequest) (*pb.GetGradeResponse, error) {
	log.Println("GetGrade: Received request")
	defer log.Println("GetGrade: Request processed")

	grade, err := h.gradeRepo.GetGrade(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.GetGradeResponse{Grade: grade}, nil
}

func (h *GradesHandler) ListGrades(ctx context.Context, req *pb.ListGradesRequest) (*pb.ListGradesResponse, error) {
	log.Println("ListGrades: Received request")
	defer log.Println("ListGrades: Request processed")

	grades, err := h.gradeRepo.ListGrades(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.ListGradesResponse{Grades: grades}, nil
}

func (h *GradesHandler) DeleteGrade(ctx context.Context, req *pb.DeleteGradeRequest) (*pb.Grade, error) {
	log.Println("DeleteGrade: Received request")
	defer log.Println("DeleteGrade: Request processed")

	userConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer userConn.Close()

	userClient := pb.NewUserServiceClient(userConn)

	// Call the DecryptJwt method on the gRPC client
	decryptReq := &pb.DecryptJwtRequest{Jwt: req.Jwt}
	decryptRes, err := userClient.DecryptJwt(ctx, decryptReq)
	if err != nil {
		return nil, err
	}

	userEmail := decryptRes.Email

	userDataReq := &pb.GetUserRequestByEmail{Email: userEmail}
	userDataRes, err := userClient.GetUserByEmail(ctx, userDataReq)

	if userDataRes.Role != "TEACHER" {
		log.Printf("You don't have permisson!: %v", err)
		return nil, errors.New("you don't have permisson")
	}

	deletedGrade, err := h.gradeRepo.DeleteGrade(ctx, req)
	if err != nil {
		return nil, err
	}
	return deletedGrade, nil
}
