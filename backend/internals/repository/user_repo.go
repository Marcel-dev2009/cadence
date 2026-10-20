package repository

import (
	"errors"
	"time"
  "github.com/Marcel-dev2009/cadence/database/models"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

type UserRepository struct {
 DB *gorm.DB
 Cache *cache.Cache
}
func NewUserRepository(db *gorm.DB, cache *cache.Cache) *UserRepository {
 return &UserRepository{DB: db, Cache: cache}
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
func (r *UserRepository) SetSpotifySyncStatus(userId string, status bool) error {
 result := r.DB.Table("users").Where("id = ?", userId).Update("spotify_sync", status)
 if result.Error != nil{
  return result.Error	
 }	
 if result.RowsAffected == 0{
 return errors.New("no user record matched that id allocation string") 	
 }
 return  nil
}
func (r *UserRepository) GetUserData (userId string) ([]models.UserMetaData, error){
 var response []models.UserMetaData
 cacheKey := "user_data:" + userId
 if cachedData, found := r.Cache.Get(cacheKey); found{
  if data, ok := cachedData.([]models.UserMetaData); ok {
		println("DEBUG: Fetching user data from RAM Cache! ⚡")
		return data, nil
  }
 }
 err := r.DB.Table("users").Where("id = ?", userId).Find(&response).Error
 if err != nil{
  if errors.Is(err, gorm.ErrRecordNotFound){
    return nil, err
  }
 }
 	r.Cache.Set(cacheKey, response, 10*time.Minute)
	println("DEBUG: Cache missed. Saved fresh data from Neon to RAM 💾")
  return response, nil
}