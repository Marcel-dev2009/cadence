package repository

import (
	"time"
	"github.com/Marcel-dev2009/cadence/database/models"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)
type EventsRepository struct {
	DB *gorm.DB
    Cache *cache.Cache
}
func NewEventsRepository(db *gorm.DB, cache *cache.Cache) *EventsRepository{
 return &EventsRepository{DB:db, Cache:cache}
}
type EventInput struct{
    EventName string `json:"event_name"`
    EventTime *time.Time `form:"event_time" time_format:"2006-01-02"`
    UserID string `json:"user_id"`
}

func (r *EventsRepository) CreateEvent(event *models.ScheduledEvents)error {
 if err := r.DB.Create(event).Error; err != nil{
    return  err
 }
 cacheKey := "user_events:"+event.UserID
 r.Cache.Delete(cacheKey)
 println("DEBUG: New Event Created! cache cleared for user:" + event.UserID)
 return  nil
}

func (r *EventsRepository) GetEvents(userId string) ([]models.EventMetaData, error){
   cacheKey := "user_events:" + userId
   if cachedData, found := r.Cache.Get(cacheKey); found {
		// Type assert the cache data back into our slice
		if events, ok := cachedData.([]models.EventMetaData); ok {
			println("DEBUG: Fetching events from RAM Cache! ⚡")
			return events, nil
		}
	}

	// Cache miss! Walk to the back room (Neon Database)
	var events []models.EventMetaData
	err := r.DB.Table("scheduled_events").Where("user_id = ?", userId).Find(&events).Error
	if err != nil {
		return nil, err
	}

	// Save to go-cache for next time (expires in 10 minutes)
	r.Cache.Set(cacheKey, events, 10*time.Minute)
	println("DEBUG: Cache missed. Saved fresh data from Neon to RAM 💾")

	return events, nil
}