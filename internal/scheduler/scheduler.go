package scheduler

import (
	"log"
	"renew-guard/internal/config"
	"renew-guard/internal/services"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron                *cron.Cron
	notificationService services.NotificationService
	config              *config.SchedulerConfig
}

func NewScheduler(notificationService services.NotificationService, cfg *config.SchedulerConfig) *Scheduler {
	// Create cron with seconds support and logging
	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.Default())))

	return &Scheduler{
		cron:                c,
		notificationService: notificationService,
		config:              cfg,
	}
}

// Start begins the scheduled jobs
func (s *Scheduler) Start() error {
	if !s.config.Enabled {
		log.Println("Scheduler is disabled")
		return nil
	}

	log.Printf("Starting scheduler with cron expression: %s", s.config.CronExpression)

	// Add notification check job
	_, err := s.cron.AddFunc(s.config.CronExpression, func() {
		log.Println("Running scheduled notification check...")
		if err := s.notificationService.CheckAndSendNotifications(s.config.NotificationDaysBefore); err != nil {
			log.Printf("Error running notification check: %v", err)
		}
	})

	if err != nil {
		return err
	}

	s.cron.Start()
	log.Println("Scheduler started successfully")

	return nil
}

// Stop halts all scheduled jobs
func (s *Scheduler) Stop() {
	log.Println("Stopping scheduler...")
	s.cron.Stop()
	log.Println("Scheduler stopped")
}

// RunNow triggers the notification check immediately (useful for testing)
func (s *Scheduler) RunNow() error {
	log.Println("Running notification check manually...")
	return s.notificationService.CheckAndSendNotifications(s.config.NotificationDaysBefore)
}
