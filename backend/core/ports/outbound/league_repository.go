package outbound

import "ludo-tournament/core/domain/models"

type LeagueRepository interface {
	Create(league *models.League) error
	GetByID(id string) (*models.League, error)
	Update(league *models.League) error
	SoftDelete(id string) error
	List() ([]models.League, error)
	ListByStatus(status models.LeagueStatus) ([]models.League, error)
}
