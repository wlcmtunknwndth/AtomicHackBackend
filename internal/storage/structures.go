package storage

import (
	"github.com/lib/pq"
)

type Model struct {
	ID uint `gorm:"primaryKey" json:"id"`
	//CreatedAt time.Time      `json:"created_at"`
	//UpdatedAt time.Time      `json:"updated_at"`
	//DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type User struct {
	Model
	Username string `json:"username" gorm:"check:length(username) < 32;check length(username) > 6"`
	Password string `json:"password" gorm:"check:length(password) < 32;check length(password) > 6"`
	isAdmin  string
}

type Request struct {
	Model
	Message string `json:"message" gorm:"type:varchar(1024)"`
}

type Response struct {
	Model
	Message string        `json:"message" gorm:"type:varchar(2048)"`
	Similar pq.Int64Array `json:"similar" gorm:"type:bigint[]"`
}

type Solved struct {
	Model
	RequestId  uint64 `json:"request_id" gorm:"check: request_id > 0"`
	SolutionId uint64 `json:"response_id" gorm:"check: solution_id > 0"`
}
