package models

type FindSubGroup struct {
	SubGroupID string `json:"SubGroupID" binding:"required"`
}
