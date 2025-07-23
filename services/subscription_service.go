package services

import (
	"errors"
	"sbs/models"
	"sbs/repositories"
	"time"
)

type SubscriptionService struct {
	repo *repositories.SubscriptionRepository
}

func NewSubscriptionService(repo *repositories.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) Create(sub *models.Subscription) error {
	return s.repo.Create(sub)
}

func (s *SubscriptionService) GetByID(id int) (*models.Subscription, error) {
	return s.repo.GetByID(id)
}

func (s *SubscriptionService) Update(sub *models.Subscription) error {
	return s.repo.Update(sub)
}

func (s *SubscriptionService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *SubscriptionService) List() ([]models.Subscription, error) {
	return s.repo.List()
}

func (s *SubscriptionService) SumSubscriptions(userID, serviceName string, startDate, endDate time.Time) (int, error) {
	if s.repo == nil {
		return 0, errors.New("repository is not initialized")
	}
	return s.repo.SumSubscriptions(userID, serviceName, startDate, endDate)
}
