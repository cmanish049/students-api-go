package storage

import "github.com/cmanish049/students-api/internal/types"

// create interface
type Storage interface {
	// define methods for storage operations
	CreateStudent(name, email string, age int) (int64, error)

	GetStudentById(id int64) (types.Student, error)
	GetStudentList() ([]types.Student, error)
	UpdateStudent(id int64, name, email string, age int) error

	DeleteStudent(id int64) error
}
