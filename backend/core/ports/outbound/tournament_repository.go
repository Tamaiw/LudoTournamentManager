package outbound

import "ludo-tournament/core/domain/models"

type TournamentRepository interface {
	Create(tournament *models.Tournament) error
	GetByID(id string) (*models.Tournament, error)
	Update(tournament *models.Tournament) error
	SoftDelete(id string) error
	List() ([]models.Tournament, error)
	ListByStatus(status models.TournamentStatus) ([]models.Tournament, error)
}
