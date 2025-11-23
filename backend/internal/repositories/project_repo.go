package repositories

import (
	"github.com/notFil/cspotlight/internal/models"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	CreateProject(project *models.Project) error
	GetProjectByID(id string) (*models.Project, error)
	UpdateProject(project *models.Project) error
	DeleteProject(id string) error
	ListProjects() ([]*models.Project, error)
	ListProjectsByTeamID(teamID string) ([]*models.Project, error)
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) CreateProject(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *projectRepository) GetProjectByID(id string) (project *models.Project, err error) {
	err = r.db.First(&project, "id = ?", id).Error
	return project, err
}

func (r *projectRepository) UpdateProject(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepository) DeleteProject(id string) error {
	return r.db.Delete(&models.Project{}, "id = ?", id).Error
}

func (r *projectRepository) ListProjects() (projects []*models.Project, err error) {
	if err = r.db.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) ListProjectsByTeamID(teamID string) (projects []*models.Project, err error) {
	if err = r.db.Where("team_id = ?", teamID).Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}
