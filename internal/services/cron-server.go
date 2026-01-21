package services

import (
	"CardFlow/internal/models"
	"CardFlow/internal/repositories"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type CronService interface{
	NotifyCardsExpiringSoon(ctx context.Context)
	ExpireCards(ctx context.Context)

}

type cronService struct {
    userRepo repositories.UserRepository
	cardRepo repositories.CardRepository
}

func NewCronService(userrepo repositories.UserRepository, cardRepo repositories.CardRepository) CronService {
    return &cronService{userRepo:userrepo, cardRepo: cardRepo}
}

func CronJobs(ctx context.Context, cronSvc CronService) {
	c := cron.New()
	// Schedules a cron job to execute the findCardsExpiringin3days method of cronSvc every day at 7:00 AM.
	// If adding the cron job fails, the error is returned.
	// The cron expression "0 7 * * *" represents 7:00 AM every day.
	if _, err := c.AddFunc("0 7 * * *", func() { cronSvc.NotifyCardsExpiringSoon(ctx) }); err != nil {
		//log error for devs
	}
	if _, err := c.AddFunc("5 7 * * *", func() { cronSvc.ExpireCards(ctx) }); err != nil {
		//log error for devs
	}

	c.Start()
	<-ctx.Done()
	c.Stop()
}



func (s *cronService) NotifyCardsExpiringSoon(ctx context.Context) {
	start := time.Now().UTC().Truncate(24*time.Hour).AddDate(0, 0, 3)
	end := start.AddDate(0, 0, 1)

	cards, err := s.cardRepo.FindCardsExpiringBetween(ctx, start, end)
	if err != nil || len(cards) == 0 {
		return
	}

	userIDs := extractUserIDs(cards)

	users, err := s.userRepo.FindUsersByIDs(ctx, userIDs)
	if err != nil {
		return
	}

	userMap := mapUsersByID(users)

	for _, card := range cards {
		user := userMap[card.UserID]
		go SendEmail(map[string]string{
			"Email":     user.Email,
			"FirstName": user.FirstName,
			"LastFour":  card.LastFour,
			"Status":      "expiring",
		})
	}
}

func (s *cronService) ExpireCards(ctx context.Context) {
	start := time.Now().UTC().Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)

	cards, err := s.cardRepo.ExpireCardsBetween(ctx, start, end)
	if err != nil || len(cards) == 0 {
		return
	}

	userIDs := extractUserIDs(cards)

	users, err := s.userRepo.FindUsersByIDs(ctx, userIDs)
	if err != nil {
		return
	}

	userMap := mapUsersByID(users)

	for _, card := range cards {
		user := userMap[card.UserID]
		go SendEmail(map[string]string{
			"Email":     user.Email,
			"FirstName": user.FirstName,
			"LastFour":  card.LastFour,
			"Status":      "expired",
		})
	}
}


func extractUserIDs(cards []models.Card) []uuid.UUID {
	set := make(map[uuid.UUID]struct{})
	for _, card := range cards {
		set[card.UserID] = struct{}{}
	}
	ids := make([]uuid.UUID, 0, len(set))
	for id := range set {
		ids = append(ids, id)
	}
	return ids
}

func mapUsersByID(users []models.User) map[uuid.UUID]models.User {
	m := make(map[uuid.UUID]models.User)
	for _, u := range users {
		m[u.ID] = u
	}
	return m
}
