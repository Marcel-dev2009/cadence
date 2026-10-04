package models 
import (
 "time"
)
type User struct {
 ID string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`	
 Name string `gorm:"uniqueIndex;not null" json:"name"`
 Email string `gorm:"uniqueIndex;not null" json:"email"`
 Password string `gorm:"uniqueIndex;not null" json:"-"`
 CreatedAt time.Time `json:"created_at"`
 UpdatedAt time.Time `json:"updated_at"`
}