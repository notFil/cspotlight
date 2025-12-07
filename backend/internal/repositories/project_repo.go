package repositories

import (
	"context"

	"github.com/notFil/cspotlight/internal/models"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, project *models.Project) error
	GetProjectByID(ctx context.Context, id string) (*models.Project, error)
	UpdateProject(ctx context.Context, project *models.Project) error
	DeleteProject(ctx context.Context, id string) error
	ListProjects(ctx context.Context) ([]*models.ProjectFetchDTO, error)
	ListProjectsByTeamID(ctx context.Context, teamID string) ([]*models.ProjectFetchDTO, error)
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) CreateProject(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *projectRepository) GetProjectByID(ctx context.Context, id string) (project *models.Project, err error) {
	err = r.db.WithContext(ctx).Preload("Team").First(&project, "id = ?", id).Error
	return project, err
}

func (r *projectRepository) UpdateProject(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *projectRepository) DeleteProject(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Project{}, "id = ?", id).Error
}

func (r *projectRepository) ListProjects(ctx context.Context) (projects []*models.ProjectFetchDTO, err error) {
	if err = r.db.WithContext(ctx).Table("projects p").
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

func (r *projectRepository) ListProjectsByTeamID(ctx context.Context, teamID string) (projects []*models.ProjectFetchDTO, err error) {
	if err = r.db.WithContext(ctx).Table("projects p").
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
