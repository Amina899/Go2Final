package repository

import (
	"context"
	"database/sql"
	"newgolang/proto/pb"
)

type AssignmentRepository struct {
	db *sql.DB
}

func NewAssignmentRepository(db *sql.DB) *AssignmentRepository {
	return &AssignmentRepository{
		db: db,
	}
}

func (repo *AssignmentRepository) CreateAssignment(ctx context.Context, req *pb.CreateAssignmentRequest) (*pb.Assignment, error) {
	query := "INSERT INTO assignment (subject_name, assignment_name) VALUES ($1, $2) RETURNING id, subject_name, assignment_name"
	row := repo.db.QueryRowContext(ctx, query, req.SubjectName, req.AssignmentName)

	var createdAssignment pb.Assignment
	if err := row.Scan(&createdAssignment.Id, &createdAssignment.SubjectName, &createdAssignment.AssignmentName); err != nil {
		return nil, err
	}

	return &createdAssignment, nil
}

func (repo *AssignmentRepository) UpdateAssignment(ctx context.Context, req *pb.UpdateAssignmentRequest) (*pb.Assignment, error) {
	query := "UPDATE assignment SET subject_name = $1, assignment_name = $2 WHERE id = $3 RETURNING id, subject_name, assignment_name"
	row := repo.db.QueryRowContext(ctx, query, req.SubjectName, req.AssignmentName, req.Id)

	var updatedAssignment pb.Assignment
	if err := row.Scan(&updatedAssignment.Id, &updatedAssignment.SubjectName, &updatedAssignment.AssignmentName); err != nil {
		return nil, err
	}

	return &updatedAssignment, nil
}

func (repo *AssignmentRepository) GetAssignment(ctx context.Context, req *pb.GetAssignmentRequest) (*pb.Assignment, error) {
	query := "SELECT id, subject_name, assignment_name, created_at FROM assignment WHERE id = $1"
	row := repo.db.QueryRowContext(ctx, query, req.Id)

	var assignment pb.Assignment
	if err := row.Scan(&assignment.Id, &assignment.SubjectName, &assignment.AssignmentName, &assignment.CreatedAt); err != nil {
		return nil, err
	}

	return &assignment, nil
}

func (repo *AssignmentRepository) ListAssignments(ctx context.Context) ([]*pb.Assignment, error) {
	query := "SELECT id, subject_name, assignment_name FROM assignment"
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []*pb.Assignment
	for rows.Next() {
		var assignment pb.Assignment
		if err := rows.Scan(&assignment.Id, &assignment.SubjectName, &assignment.AssignmentName); err != nil {
			return nil, err
		}
		assignments = append(assignments, &assignment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return assignments, nil
}

func (repo *AssignmentRepository) DeleteAssignment(ctx context.Context, req *pb.DeleteAssignmentRequest) (*pb.Assignment, error) {
	query := "DELETE FROM assignment WHERE id = $1 RETURNING id, subject_name, assignment_name"
	row := repo.db.QueryRowContext(ctx, query, req.Id)

	var deletedAssignment pb.Assignment
	if err := row.Scan(&deletedAssignment.Id, &deletedAssignment.SubjectName, &deletedAssignment.AssignmentName); err != nil {
		return nil, err
	}

	return &deletedAssignment, nil
}
