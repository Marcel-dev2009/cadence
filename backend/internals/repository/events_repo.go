package repository

import (
	"errors"
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
 cacheKey := "user_events:" + event.UserID
 r.Cache.Delete(cacheKey)
 r.Cache.Delete("user_upcoming_event:" + event.UserID)
 println("DEBUG: New Event Created! cache cleared for user:" + event.UserID)
 return  nil
}
func (r *EventsRepository) UpdateEventStatus(userId, eventID string) error {
  var cacheKey = "user_events:" + userId 
 result := r.DB.Table("scheduled_events").Where("event_id = ?", eventID).Update("event_status", "finished")
 if result.Error != nil {
  return result.Error
  }	
 if result.RowsAffected == 0 {
 return errors.New("failed to update event status: event record not found")
   }
 r.Cache.Delete(cacheKey)
 r.Cache.Delete("user_upcoming_event:" + userId) 
 return nil
}
func (r *EventsRepository) GetMostUpcomingEvent(userID string) (*models.ScheduledEvents, error) {
	cacheKey := "user_upcoming_event:" + userID

	// 1. Try Cache First ⚡
	if cachedData, found := r.Cache.Get(cacheKey); found {
		if event, ok := cachedData.(*models.ScheduledEvents); ok {
			println("DEBUG: Most upcoming event loaded from RAM Cache! ⚡")
			return event, nil
		}
	}

	// 2. Cache Miss: Query Neon with specific ordering and a Limit of 1
	var event models.ScheduledEvents
	err := r.DB.Table("scheduled_events").
		Where("user_id = ?", userID).
		Where("event_time >= ?", time.Now()). // Only future events
		Order("event_time ASC").               // Earliest time first (closest to now) [1]
		Limit(1).                             // Look up exactly 1 item [2]
		First(&event).                        // Pour data into our single struct box
		Error

	// 3. Handle Empty Results
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If they have no future events scheduled, return nil safely without an error
			return nil, nil 
		}
		return nil, err
	}

	// 4. Save to Cache for 5 minutes (Short expiration since time.Now() moves)
	r.Cache.Set(cacheKey, &event, 5*time.Minute)
	println("DEBUG: Cache missed. Saved highlighted event from Neon 💾")

	return &event, nil
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