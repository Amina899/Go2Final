package internal

import (
	"context"
	"database/sql"
	"newgolang/proto/pb"
	"time"
)

type GradeRepository struct {
	db *sql.DB
}

func NewGradeRepository(db *sql.DB) *GradeRepository {
	return &GradeRepository{
		db: db,
	}
}

func (repo *GradeRepository) CreateGrade(ctx context.Context, req *pb.CreateGradeRequest) (*pb.Grade, error) {
	query := "INSERT INTO grades (user_id, assignment_id, grade, created_at) VALUES ($1, $2, $3, $4) RETURNING id, user_id, assignment_id, grade"
	row := repo.db.QueryRowContext(ctx, query, req.UserId, req.AssignmentId, req.Grade, time.Now())

	var createdGrade pb.Grade
	if err := row.Scan(&createdGrade.Id, &createdGrade.UserId, &createdGrade.AssignmentId, &createdGrade.Grade); err != nil {
		return nil, err
	}

	return &createdGrade, nil
}

func (repo *GradeRepository) UpdateGrade(ctx context.Context, req *pb.UpdateGradeRequest) (*pb.Grade, error) {
	query := "UPDATE grades SET grade = $1 WHERE id = $2 RETURNING id, user_id, assignment_id, grade"
	row := repo.db.QueryRowContext(ctx, query, req.Grade, req.Id)

	var updatedGrade pb.Grade
	if err := row.Scan(&updatedGrade.Id, &updatedGrade.UserId, &updatedGrade.AssignmentId, &updatedGrade.Grade); err != nil {
		return nil, err
	}

	return &updatedGrade, nil
}

func (repo *GradeRepository) GetGrade(ctx context.Context, req *pb.GetGradeRequest) (*pb.Grade, error) {
	query := "SELECT id, user_id, assignment_id, grade FROM grades WHERE id = $1"
	row := repo.db.QueryRowContext(ctx, query, req.Id)

	var grade pb.Grade
	if err := row.Scan(&grade.Id, &grade.UserId, &grade.AssignmentId, &grade.Grade); err != nil {
		return nil, err
	}

	return &grade, nil
}

func (repo *GradeRepository) ListGrades(ctx context.Context, req *pb.ListGradesRequest) ([]*pb.Grade, error) {
	query := "SELECT id, user_id, assignment_id, grade FROM grades LIMIT $1 OFFSET $2"
	rows, err := repo.db.QueryContext(ctx, query, req.PageSize, req.PageNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grades []*pb.Grade
	for rows.Next() {
		var grade pb.Grade
		if err := rows.Scan(&grade.Id, &grade.UserId, &grade.AssignmentId, &grade.Grade); err != nil {
			return nil, err
		}
		grades = append(grades, &grade)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return grades, nil
}

func (repo *GradeRepository) DeleteGrade(ctx context.Context, req *pb.DeleteGradeRequest) (*pb.Grade, error) {
	query := "DELETE FROM grades WHERE id = $1 RETURNING id, user_id, assignment_id, grade"
	row := repo.db.QueryRowContext(ctx, query, req.Id)

	var deletedGrade pb.Grade
	if err := row.Scan(&deletedGrade.Id, &deletedGrade.UserId, &deletedGrade.AssignmentId, &deletedGrade.Grade); err != nil {
		return nil, err
	}

	return &deletedGrade, nil
}
