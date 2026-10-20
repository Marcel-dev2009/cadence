package models 
import (
 "time"
)
type User struct {
 ID string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`	
 Name string `gorm:"uniqueIndex;not null" json:"name"`
 Email string `gorm:"uniqueIndex;not null" json:"email"`
 Password string `gorm:"uniqueIndex;not null" json:"-"`
 Image   string  `gorm:"uniqueIndex" json:"image"`
 SyncGoogleCalender bool `gorm:"default:false" json:"sync_google_calender"`
 SpotifyTrackUrl string  `gorm:"default:https://open.spotify.com/track/2a1o6ZejUi8U3wzzOtCOYw?si=bbc071489c834353" json:"spotify_track_url"`
 SyncSpotify bool   `gorm:"default:false" json:"sync_spotify"`
 CreatedAt time.Time `json:"created_at"`
 UpdatedAt time.Time `json:"updated_at"`
}

type ScheduledEvents struct {
EventID    string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"event_id"`
EventName  string    `gorm:"type:varchar(255);not null" json:"event_name"`
EventTime *time.Time `gorm:"type:timestamp with time zone;not null" json:"event_time"`
EventStatus string    `gorm:"type:varchar(50);default:'pending'" json:"event_status"`
UserID string `gorm:"type:varchar(255);index;not null" json:"user_id"`
User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
type EventMetaData struct {
EventID    string    `gorm:"serializer:json;type:jsonb"`
EventName  string    `gorm:"serializer:json;type:jsonb"`
EventTime *time.Time `gorm:"serializer:json;type:jsonb"`
EventStatus string   `gorm:"serializer:json;type:jsonb"`
CreatedAt   time.Time  `gorm:"serializer:json;type:jsonb"`
UpdatedAt   time.Time `gorm:"serializer:json;type:jsonb"`
}
type UserMetaData struct {
 ID string  `gorm:"serializer:json;type:jsonb"`	
 Name string  `gorm:"serializer:json;type:jsonb"`
 Email string  `gorm:"serializer:json;type:jsonb"`
 Password string  `gorm:"serializer:json;type:jsonb"`
 Image   string   `gorm:"serializer:json;type:jsonb"`	
}