package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/nibroos/elearning-go/service/internal/config"
	"github.com/nibroos/elearning-go/service/internal/dtos"
	"github.com/nibroos/elearning-go/service/internal/models"
	"github.com/nibroos/elearning-go/service/internal/utils"
	"gorm.io/gorm"
)

type SubscribeRepository struct {
	db    *gorm.DB
	sqlDB *sqlx.DB
}

func NewSubscribeRepository(db *gorm.DB, sqlDB *sqlx.DB) *SubscribeRepository {
	return &SubscribeRepository{
		db:    db,
		sqlDB: sqlDB,
	}
}

func (r *SubscribeRepository) GetSubscribes(ctx context.Context, filters map[string]string) ([]dtos.SubscribeListDTO, int, error) {
	subscribes := []dtos.SubscribeListDTO{}
	var total int

	query := `SELECT *
    FROM ( 
        SELECT s.id, s.name, s.description, s.created_at, s.updated_at, s.deleted_at,
        cu.name as created_by_name,
        uu.name as updated_by_name

        FROM subscribes s
        JOIN users cu ON s.created_by_id = cu.id
        LEFT JOIN users uu ON s.updated_by_id = uu.id
    ) AS alias WHERE 1=1 AND deleted_at IS NULL`

	countQuery := `SELECT COUNT(*) FROM (
        SELECT s.id, s.name, s.description, s.created_at, s.updated_at, s.deleted_at,
        cu.name as created_by_name,
        uu.name as updated_by_name

        FROM subscribes s
        JOIN users cu ON s.created_by_id = cu.id
        LEFT JOIN users uu ON s.updated_by_id = uu.id
    ) AS alias WHERE 1=1 AND deleted_at IS NULL`

	var args []interface{}

	i := 1
	for key, value := range filters {
		switch key {
		case "name", "description":
			if value != "" {
				query += fmt.Sprintf(" AND %s ILIKE $%d", key, i)
				countQuery += fmt.Sprintf(" AND %s ILIKE $%d", key, i)
				args = append(args, "%"+value+"%")
				i++
			}
		}
	}

	if value, ok := filters["global"]; ok && value != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", i, i+1)
		countQuery += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", i, i+1)
		args = append(args, "%"+value+"%", "%"+value+"%")
		i += 2
	}

	countArgs := append([]interface{}{}, args...)

	// Channels for concurrent execution
	countChan := make(chan error)
	selectChan := make(chan error)

	// Goroutine for count query
	go func() {
		err := r.sqlDB.GetContext(ctx, &total, countQuery, countArgs...)
		countChan <- err
	}()

	orderColumn := utils.GetStringOrDefault(filters["order_column"], "id")
	orderDirection := utils.GetStringOrDefault(filters["order_direction"], "asc")
	query += fmt.Sprintf(" ORDER BY %s %s", orderColumn, orderDirection)

	perPage := utils.GetIntOrDefault(filters["per_page"], 10)
	currentPage := utils.GetIntOrDefault(filters["page"], 1)

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, perPage, (currentPage-1)*perPage)

	// Goroutine for select query
	go func() {
		err := r.sqlDB.SelectContext(ctx, &subscribes, query, args...)
		selectChan <- err
	}()

	// Wait for both goroutines to finish
	countErr := <-countChan
	selectErr := <-selectChan

	if countErr != nil {
		return nil, 0, countErr
	}

	if selectErr != nil {
		return nil, 0, selectErr
	}

	return subscribes, total, nil
}

func (r *SubscribeRepository) GetSubscribesFromRedis(ctx context.Context, filters map[string]string) ([]dtos.SubscribeListDTO, int, error) {
	var subscribes []dtos.SubscribeListDTO
	var total int

	// Fetch data from Redis
	data, err := config.RedisClient.Get(ctx, "subscribes").Result()
	if err != nil {
		return nil, 0, err
	}

	// Unmarshal the data
	if err := json.Unmarshal([]byte(data), &subscribes); err != nil {
		return nil, 0, err
	}

	// Apply filters
	filteredSubscribes := []dtos.SubscribeListDTO{}
	for _, subscribe := range subscribes {
		match := true
		for key, value := range filters {
			switch key {
			case "name":
				if value != "" && !utils.ContainsIgnoreCase(subscribe.Name, value) {
					match = false
				}
			case "description":
				if value != "" && !utils.ContainsIgnoreCase(subscribe.Description, value) {
					match = false
				}
			case "global":
				if value != "" && !utils.ContainsIgnoreCase(subscribe.Name, value) && !utils.ContainsIgnoreCase(subscribe.Description, value) {
					match = false
				}
			}
		}
		if match {
			filteredSubscribes = append(filteredSubscribes, subscribe)
		}
	}

	total = len(filteredSubscribes)

	// Apply pagination
	perPage := utils.GetIntOrDefault(filters["per_page"], 10)
	currentPage := utils.GetIntOrDefault(filters["page"], 1)
	start := (currentPage - 1) * perPage
	end := start + perPage
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	paginatedSubscribes := filteredSubscribes[start:end]

	return paginatedSubscribes, total, nil
}

// func (r *SubscribeRepository) FetchAndCacheSubscribes(ctx context.Context) error {
// 	var subscribes []dtos.SubscribeListDTO

// 	query := `SELECT s.id, s.name, s.description, s.created_at, s.updated_at, s.deleted_at,
//         cu.name as created_by_name,
//         uu.name as updated_by_name
//     FROM subscribes s
//     JOIN users cu ON s.created_by_id = cu.id
//     LEFT JOIN users uu ON s.updated_by_id = uu.id
//     WHERE s.deleted_at IS NULL`

// 	err := r.sqlDB.SelectContext(ctx, &subscribes, query)
// 	if err != nil {
// 		return err
// 	}

// 	return r.redisCache.FetchAndCacheSubscribes(ctx, subscribes)
// }

func (r *SubscribeRepository) GetSubscribeByID(ctx context.Context, id uint) (*dtos.SubscribeDetailDTO, error) {
	var subscribe dtos.SubscribeDetailDTO

	query := `SELECT s.*,
	cu.name as created_by_name,
	uu.name as updated_by_name

	FROM subscribes s
	JOIN users cu ON s.created_by_id = cu.id
	LEFT JOIN users uu ON s.updated_by_id = uu.id
	WHERE s.id = $1 AND s.deleted_at IS NULL`
	if err := r.sqlDB.Get(&subscribe, query, id); err != nil {
		return nil, err
	}

	return &subscribe, nil
}

// BeginTransaction starts a new transaction
func (r *SubscribeRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *SubscribeRepository) CreateSubscribe(tx *gorm.DB, subscribe *models.Subscribe) error {
	if err := tx.Create(subscribe).Error; err != nil {
		return err
	}
	return nil
}

func (r *SubscribeRepository) UpdateSubscribe(tx *gorm.DB, subscribe *models.Subscribe) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(subscribe).Error; err != nil {
			return err
		}
		return nil
	})

}

func (r *SubscribeRepository) DeleteSubscribe(tx *gorm.DB, id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// if err := tx.Unscoped().Delete(&models.Subscribe{}, id).Error; err != nil {
		if err := tx.Delete(&models.Subscribe{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *SubscribeRepository) RestoreSubscribe(tx *gorm.DB, id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var subscribe models.Subscribe
		if err := tx.Unscoped().First(&subscribe, id).Error; err != nil {
			return err
		}
		return tx.Model(&subscribe).Update("deleted_at", nil).Error
	})
}
