package outbound

import "ludo-tournament/core/domain/models"

type PlayerRepository interface {
	Create(player *models.Player) error
	GetByID(id string) (*models.Player, error)
	GetByUserID(userID string) (*models.Player, error)
	Update(player *models.Player) error
	SoftDelete(id string) error
}
