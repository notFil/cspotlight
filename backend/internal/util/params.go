package util

type ProjectParams struct {
	ProjectID string `uri:"projectID" binding:"required,uuid"`
}

type UserParams struct {
	UserID string `uri:"userID" binding:"required,uuid"`
}

type TeamParams struct {
	TeamID string `uri:"teamID" binding:"required,uuid"`
}
