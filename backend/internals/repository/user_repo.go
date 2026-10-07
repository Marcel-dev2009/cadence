package repository

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserRepository struct {
 DB *gorm.DB
}
func NewUserRepository(db *gorm.DB) *UserRepository {
 return &UserRepository{DB: db}
}
func (r *UserRepository) GetUserId(c *gin.Context) string{
 UserID := c.GetString("userID")
     if UserID == ""{
       c.JSON(http.StatusUnauthorized, gin.H{"error":"User session not found or is expired"})
       return "User Id not found"
     }
  return UserID
}
func (r *UserRepository) GetgoogleCalenderSync (userId string) (bool, error) {
 var result struct {
 SyncGoogleCalendar bool `gorm:"column:sync_google_calender"`
 }
 err := r.DB.Table("users").Where("id = ?", userId).First(&result).Error
 if err != nil{
  if errors.Is(err, gorm.ErrRecordNotFound){
   return false, err	
  }	
 }
 return  result.SyncGoogleCalendar, nil
}
func (r *UserRepository) SetgoogleSyncStatus(userId string, status bool) error {
 result := r.DB.Table("users").Where("id = ?", userId).Update("sync_google_calender", status)
 if result.Error != nil{
  return result.Error	
 }	
 if result.RowsAffected == 0{
 return errors.New("no user record matched that id allocation string") 	
 }
 return  nil
}

