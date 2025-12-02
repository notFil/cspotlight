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
	ListProjects() ([]*models.ProjectFetchDTO, error)
	ListProjectsByTeamID(teamID string) ([]*models.ProjectFetchDTO, error)
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
	err = r.db.Preload("Team").First(&project, "id = ?", id).Error
	return project, err
}

func (r *projectRepository) UpdateProject(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *projectRepository) DeleteProject(id string) error {
	return r.db.Delete(&models.Project{}, "id = ?", id).Error
}

func (r *projectRepository) ListProjects() (projects []*models.ProjectFetchDTO, err error) {
	if err = r.db.Table("projects p").
		Joins("LEFT JOIN teams t ON t.id = p.team_id").
		Select(`p.*, t.name as team_name, (
			SELECT MAX(created_at)
			FROM csp_reports r
			WHERE r.project_id = p.id
		) AS last_active`).
		Scan(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) ListProjectsByTeamID(teamID string) (projects []*models.ProjectFetchDTO, err error) {
	if err = r.db.Table("projects p").
		Joins("LEFT JOIN teams t ON t.id = p.team_id").
		Select(`p.*, t.name as team_name, (
			SELECT MAX(created_at)
			FROM csp_reports r
			WHERE r.project_id = p.id
		) AS last_active`).
		Where("p.team_id = ?", teamID).
		Scan(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}
