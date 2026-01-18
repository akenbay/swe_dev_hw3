package storage

import (
	"context"
	"university/internal/model"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetStudentByID(id string) (*model.StudentResponse, error) {
	var err error

	query := `
	SELECT s.id, s.first_name, s.last_name, s.gender, s.birth_date, g.name
	FROM students s
	JOIN groups g ON s.group_id = g.id
	WHERE s.id = $1
	`

	var student model.StudentResponse

	err = r.db.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&student.ID,
		&student.FirstName,
		&student.LastName,
		&student.Gender,
		&student.BirthDate,
		&student.GroupName,
	)

	if err != nil {
		return nil, err
	}

	return &student, nil
}

func (r *Repository) GetAllSchedules() ([]model.ScheduleResponse, error) {
	var err error

	query := `
	SELECT sc.id, f.name, g.name, s.name, sc.class_time
	FROM schedule sc
	JOIN faculties f ON sc.faculty_id = f.id
	JOIN groups g ON sc.group_id = g.id
	JOIN subjects s ON sc.subject_id = s.id
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []model.ScheduleResponse

	for rows.Next() {
		var schedule model.ScheduleResponse
		err := rows.Scan(
			&schedule.ID,
			&schedule.Faculty,
			&schedule.Group,
			&schedule.Subject,
			&schedule.ClassTime,
		)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func (r *Repository) GetGroupSchedule(groupID string) ([]model.ScheduleResponse, error) {
	query := `
	SELECT sc.id, f.name, g.name, s.name, sc.class_time
	FROM schedule sc
	JOIN faculties f ON sc.faculty_id = f.id
	JOIN groups g ON sc.group_id = g.id
	JOIN subjects s ON sc.subject_id = s.id
	WHERE sc.group_id = $1
	`

	rows, err := r.db.Query(context.Background(), query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []model.ScheduleResponse
	for rows.Next() {
		var schedule model.ScheduleResponse
		err := rows.Scan(
			&schedule.ID,
			&schedule.Faculty,
			&schedule.Group,
			&schedule.Subject,
			&schedule.ClassTime,
		)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}

func (r *Repository) CreateAttendanceRecord(record *model.AttendanceRecord) error {
	query := `
	INSERT INTO attendance (student_id, subject_id, visit_day, visited)
	VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		record.StudentID,
		record.SubjectID,
		record.VisitDay,
		record.Visited,
	)

	return err
}

func (r *Repository) GetAttendanceRecordsByStudentID(studentID string) ([]model.AttendanceRecord, error) {
	query := `
	SELECT id, student_id, subject_id, visit_day, visited
	FROM attendance
	WHERE student_id = $1
	LIMIT 5
	`

	rows, err := r.db.Query(context.Background(), query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.AttendanceRecord
	for rows.Next() {
		var record model.AttendanceRecord
		err := rows.Scan(
			&record.ID,
			&record.StudentID,
			&record.SubjectID,
			&record.VisitDay,
			&record.Visited,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func (r *Repository) GetAttendanceRecordsBySubjectID(subjectID string) ([]model.AttendanceRecord, error) {
	query := `
	SELECT id, student_id, subject_id, visit_day, visited
	FROM attendance
	WHERE subject_id = $1
	LIMIT 5
	`

	rows, err := r.db.Query(context.Background(), query, subjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.AttendanceRecord
	for rows.Next() {
		var record model.AttendanceRecord
		err := rows.Scan(
			&record.ID,
			&record.StudentID,
			&record.SubjectID,
			&record.VisitDay,
			&record.Visited,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

// GetUserByEmail retrieves a user by email
func (r *Repository) GetUserByEmail(email string) (*model.User, error) {
	query := `
	SELECT id, email, password_hash, is_active, created_at
	FROM users
	WHERE email = $1
	`

	var user model.User
	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a new user account
func (r *Repository) CreateUser(email, passwordHash string) (*model.User, error) {
	query := `
	INSERT INTO users (email, password_hash)
	VALUES ($1, $2)
	RETURNING id, email, password_hash, is_active, created_at
	`

	var user model.User
	err := r.db.QueryRow(
		context.Background(),
		query,
		email,
		passwordHash,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.IsActive,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserByID retrieves a user by ID with roles
func (r *Repository) GetUserByID(userID string) (*model.UserResponse, error) {
	query := `
	SELECT u.id, u.email, u.is_active, u.created_at
	FROM users u
	WHERE u.id = $1
	`

	var user model.UserResponse
	err := r.db.QueryRow(context.Background(), query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.IsActive,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Get user roles
	rolesQuery := `
	SELECT r.name
	FROM roles r
	JOIN user_roles ur ON r.id = ur.role_id
	WHERE ur.user_id = $1
	`

	rows, err := r.db.Query(context.Background(), rolesQuery, userID)
	if err == nil {
		defer rows.Close()
		user.Roles = []string{}
		for rows.Next() {
			var roleName string
			if err := rows.Scan(&roleName); err == nil {
				user.Roles = append(user.Roles, roleName)
			}
		}
	}

	return &user, nil
}
