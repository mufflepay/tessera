package models

import (
	"time"

	"github.com/google/uuid"
)

type Stat struct {
	StatsID             uuid.UUID `json:"stats_id"`
	ActiveAgents        int       `json:"active_users"`
	TotalTickets        int       `json:"total_tickets"`
	TotalOpen           int       `json:"total_open"`
	TotalClosed         int       `json:"total_closed"`
	TotalInProgress     int       `json:"total_in_progress"`
	TotalResolved       int       `json:"total_resolved"`
	TotalUnresolved     int       `json:"total_unresolved"`
	TotalReopened       int       `json:"total_reopened"`
	TotalHighPriority   int       `json:"total_high_priority"`
	TotalMediumPriority int       `json:"total_medium_priority"`
	TotalLowPriority    int       `json:"total_low_priority"`
	TotalEscalated      int       `json:"total_escalated"`
	TotalUnassigned     int       `json:"total_unassigned"`
	TotalAssigned       int       `json:"total_assigned"`
	TotalSpam           int       `json:"total_spam"`
	TotalDeleted        int       `json:"total_deleted"`
	AverageTimeToClose  int       `json:"average_time_to_close"`
	HourlyTickets       int       `json:"hourly_tickets"`
	DailyTickets        int       `json:"daily_tickets"`
	WeeklyTickets       int       `json:"weekly_tickets"`
	MonthlyTickets      int       `json:"monthly_tickets"`
	YearlyTickets       int       `json:"yearly_tickets"`
	CreatedOn           time.Time `json:"created_on" gorm:"autoCreateTime;not null" validate:"required"`
	ModifiedOn          time.Time `json:"modified_on" gorm:"autoUpdateTime;not null" validate:"required"`
}
