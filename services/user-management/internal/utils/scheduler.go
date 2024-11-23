package utils

import (
	"fmt"
	"log"
	"time"
	"user-management/internal/model"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// Scheduler handles scheduled tasks
type Scheduler struct {
	DB          *gorm.DB
	EmailSender EmailSender
}

// NewScheduler creates a new Scheduler instance
func NewScheduler(db *gorm.DB, emailSender EmailSender) *Scheduler {
	return &Scheduler{
		DB:          db,
		EmailSender: emailSender,
	}
}

// StartMonthlySummaryScheduler starts the monthly summary scheduler
func (s *Scheduler) StartMonthlySummaryScheduler() {
	c := cron.New()
	c.AddFunc("0 0 1 * * *", func() {
		s.sendMonthlySummaries()
	})
	c.Start()
}

// sendMonthlySummaries sends monthly summaries to all users
func (s *Scheduler) sendMonthlySummaries() {
	users := []model.User{}
	if result := s.DB.Find(&users); result.Error != nil {
		log.Printf("Error fetching users: %v", result.Error)
		return
	}
	for _, user := range users {
		go s.sendMonthlySummary(user)
	}
}

// sendMonthlySummary sends a monthly summary to a specific user
func (s *Scheduler) sendMonthlySummary(user model.User) {
	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	start := firstOfMonth.AddDate(0, -1, 0)
	end := firstOfMonth.Add(-time.Second)

	var activities []model.Activity
	if result := s.DB.Where("user_id = ? AND created_at >= ? AND created_at <= ?", user.ID, start, end).Find(&activities); result.Error != nil {
		log.Printf("Error fetching activities: %v", result.Error)
		return
	}

	activitySummary := fmt.Sprintf("You had %d activities last month.\n", len(activities))
	for _, activity := range activities {
		activitySummary += fmt.Sprintf("- %s at %s\n", activity.Description, activity.CreatedAt.Format(time.RFC822))
	}

	emailSubject := "Monthly Account Activity Summary"
	emailBody := "Dear " + user.Name + ",\n\nHere is your account activity summary for the last month:\n\n" + activitySummary

	if err := s.EmailSender.SendEmail(user.Email, emailSubject, emailBody); err != nil {
		log.Printf("Error sending email: %v", err)
	}
}
