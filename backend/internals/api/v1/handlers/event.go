package handlers

import (
	"net/http"
	"strings"
	"time"
  "github.com/Marcel-dev2009/cadence/database/config"
	"github.com/Marcel-dev2009/cadence/database/models"
	"github.com/Marcel-dev2009/cadence/repository"
	"github.com/gin-gonic/gin"
)

type EventInput struct{
	EventName string `json:"event_name"`
    EventTime *time.Time `form:"event_time" time_format:"2006-01-02"`
}
 
func CreateEventHandler(c *gin.Context) {
    var input EventInput
      db := config.DB	
      dataCache := config.DataListCache
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
      userRepo := repository.NewUserRepository(db, dataCache)
      eventRepo := repository.NewEventsRepository(db, dataCache)
     userID := c.GetString("userID")
     if userID == ""{
       c.JSON(http.StatusUnauthorized, gin.H{"error":"User session not found or is expired"}) 
       return
     }
      isSynced, err := userRepo.GetgoogleCalenderSync(userID); 
   if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking sync status"})
      return
    } 
    if !isSynced {
      c.JSON(http.StatusForbidden, gin.H{"error": "Please sync Google Calendar first"})
      return
    }
    if strings.TrimSpace(input.EventName) == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You must give the event a name"})
        return
    }
    if err := c.ShouldBindQuery(&input); err != nil{
        c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid Time format"})
    }
   if input.EventTime == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You must set a specific time for the event via ?event_time=YYYY-MM-DD"})
		return
	}
    if input.EventTime.Before(time.Now()) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot schedule an event in the past!"})
        return
    }
    // create the event for them
   newEvent := models.ScheduledEvents{
    EventName: input.EventName,
    EventTime: input.EventTime,
    UserID: userID,
   } 
    if err := eventRepo.CreateEvent(&newEvent); err != nil{
      c.JSON(http.StatusBadRequest, gin.H{"error":"failed to create event"})
      return
    }
    c.JSON(http.StatusCreated, gin.H{
      "message":"event created succcesfully",
      "data": newEvent,
    })
   }

  func GetUserEventsHandler(c *gin.Context) {
     db := config.DB	
      dataCache := config.DataListCache
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	eventRepo := repository.NewEventsRepository(db, dataCache)

	events, err := eventRepo.GetEvents(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to look up events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": events})
}

func GetDashboardHighlightHandler(c *gin.Context) {
    db := config.DB	
      dataCache := config.DataListCache
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	eventRepo := repository.NewEventsRepository(db, dataCache)

	// Fire our optimized query method
	upcomingEvent, err := eventRepo.GetMostUpcomingEvent(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load dashboard data"})
		return
	}

	// If upcomingEvent is nil, the user simply hasn't scheduled anything yet
	if upcomingEvent == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "No upcoming events scheduled",
			"data":    nil,
		})
		return
	}

	// Return the single event object cleanly
	c.JSON(http.StatusOK, gin.H{
		"message": "Highlighted event loaded",
		"data":    upcomingEvent,
	})
}
func UpdateEventStatusHandler (c *gin.Context){
    db := config.DB	
    dataCache := config.DataListCache
 eventID := c.Param("id")
 userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
 if eventID == ""{
  c.JSON(http.StatusBadRequest, gin.H{"error":"missing event ID from request"})
  c.Abort()
  return
 }
 eventRepo := repository.NewEventsRepository(db, dataCache)
 if err := eventRepo.UpdateEventStatus(userID, eventID); err != nil{
  c.JSON(http.StatusBadRequest, gin.H{"error":"failed to update event status"})
  return
 }
 c.JSON(http.StatusOK, gin.H{"message":"event status updated succesfully"})
}