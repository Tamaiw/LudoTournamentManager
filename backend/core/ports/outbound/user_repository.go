package outbound

import "ludo-tournament/core/domain/models"

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	SoftDelete(id string) error
	List() ([]models.User, error)
	UpdateLastActive(id string) error
}
