package handlers

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"newgolang/assignment-service/repository"
	"newgolang/proto/pb"
)

type AssignmentHandler struct {
	pb.UnimplementedAssignmentServiceServer
	assignmentRepo repository.AssignmentRepository
}

func NewAssignmentHandler(assignmentRepo repository.AssignmentRepository) *AssignmentHandler {
	return &AssignmentHandler{
		assignmentRepo: assignmentRepo,
	}
}

func (h *AssignmentHandler) CreateAssignment(ctx context.Context, req *pb.CreateAssignmentRequest) (*pb.Assignment, error) {
	log.Println("CreateAssignment: Received request")
	defer log.Println("CreateAssignment: Request processed")

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

	createdAssignment, err := h.assignmentRepo.CreateAssignment(ctx, req)
	if err != nil {
		log.Printf("Error creating assignment: %v", err)
		return nil, err
	}
	return createdAssignment, nil
}

func (h *AssignmentHandler) UpdateAssignment(ctx context.Context, req *pb.UpdateAssignmentRequest) (*pb.Assignment, error) {
	log.Println("UpdateAssignment: Received request")
	defer log.Println("UpdateAssignment: Request processed")

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

	updatedAssignment, err := h.assignmentRepo.UpdateAssignment(ctx, req)
	if err != nil {
		log.Printf("Error updating assignment: %v", err)
		return nil, err
	}
	return updatedAssignment, nil
}

func (h *AssignmentHandler) GetAssignment(ctx context.Context, req *pb.GetAssignmentRequest) (*pb.GetAssignmentResponse, error) {
	log.Println("GetAssignment: Received request")
	defer log.Println("GetAssignment: Request processed")

	assignment, err := h.assignmentRepo.GetAssignment(ctx, req)
	if err != nil {
		log.Printf("Error getting assignment: %v", err)
		return nil, err
	}
	return &pb.GetAssignmentResponse{Assignment: assignment}, nil
}

func (h *AssignmentHandler) ListAssignments(ctx context.Context, req *pb.ListAssignmentsRequest) (*pb.ListAssignmentsResponse, error) {
	log.Println("ListAssignments: Received request")
	defer log.Println("ListAssignments: Request processed")

	assignments, err := h.assignmentRepo.ListAssignments(ctx)
	if err != nil {
		log.Printf("Error listing assignments: %v", err)
		return nil, err
	}
	return &pb.ListAssignmentsResponse{Assignments: assignments}, nil
}

func (h *AssignmentHandler) DeleteAssignment(ctx context.Context, req *pb.DeleteAssignmentRequest) (*pb.Assignment, error) {
	log.Println("DeleteAssignment: Received request")
	defer log.Println("DeleteAssignment: Request processed")

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

	deletedAssignment, err := h.assignmentRepo.DeleteAssignment(ctx, req)
	if err != nil {
		log.Printf("Error deleting assignment: %v", err)
		return nil, err
	}
	return deletedAssignment, nil
}
