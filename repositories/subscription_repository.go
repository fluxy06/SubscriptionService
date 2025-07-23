package repositories

import (
	"database/sql"
	"sbs/models"
	"time"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create добавляет новую подписку и заполняет ID, CreatedAt, UpdatedAt в модели
func (r *SubscriptionRepository) Create(s *models.Subscription) error {
	query := `
        INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at;
    `
	return r.db.QueryRow(
		query,
		s.ServiceName,
		s.Price,
		s.UserID,
		s.StartDate,
		s.EndDate,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *SubscriptionRepository) SumByFilter(userID, serviceName string, startDate, endDate time.Time) (int, error) {
	var total int
	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE user_id = $1
		  AND service_name = $2
		  AND start_date <= $4
		  AND (end_date IS NULL OR end_date >= $3)
	`
	err := r.db.QueryRow(query, userID, serviceName, startDate, endDate).Scan(&total)
	return total, err
}

func (r *SubscriptionRepository) SumSubscriptions(userID, serviceName string, startDate, endDate time.Time) (int, error) {
	var total int
	query := `
        SELECT COALESCE(SUM(price), 0) FROM subscriptions
        WHERE user_id = $1
          AND service_name = $2
          AND start_date <= $4
          AND (end_date IS NULL OR end_date >= $3)
    `
	// start_date <= endDate AND (end_date IS NULL OR end_date >= startDate)
	// так считаются подписки, которые пересекаются с периодом startDate - endDate

	err := r.db.QueryRow(query, userID, serviceName, startDate, endDate).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetByID возвращает подписку по ID
func (r *SubscriptionRepository) GetByID(id int) (*models.Subscription, error) {
	query := `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
        FROM subscriptions WHERE id = $1;
    `
	s := &models.Subscription{}
	err := r.db.QueryRow(query, id).Scan(
		&s.ID,
		&s.ServiceName,
		&s.Price,
		&s.UserID,
		&s.StartDate,
		&s.EndDate,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Update меняет данные подписки
func (r *SubscriptionRepository) Update(s *models.Subscription) error {
	query := `
        UPDATE subscriptions
        SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5, updated_at = now()
        WHERE id = $6;
    `
	_, err := r.db.Exec(query, s.ServiceName, s.Price, s.UserID, s.StartDate, s.EndDate, s.ID)
	return err
}

// Delete удаляет подписку по ID
func (r *SubscriptionRepository) Delete(id int) error {
	query := `DELETE FROM subscriptions WHERE id = $1;`
	_, err := r.db.Exec(query, id)
	return err
}

// List возвращает все подписки
func (r *SubscriptionRepository) List() ([]models.Subscription, error) {
	query := `
        SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
        FROM subscriptions;
    `
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		var s models.Subscription
		err := rows.Scan(
			&s.ID,
			&s.ServiceName,
			&s.Price,
			&s.UserID,
			&s.StartDate,
			&s.EndDate,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, nil
}
