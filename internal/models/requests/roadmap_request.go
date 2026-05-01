package requests

import roadmapEnum "be-file-uploader/pkg/enums/roadmap"

type RoadmapAddRequest struct {
	Title string `json:"title" validate:"required,min=3,max=200"`
}

type RoadmapUpdateRequest struct {
	ID     int                `json:"id" validate:"required"`
	Title  string             `json:"title" validate:"required,min=3,max=200"`
	Status roadmapEnum.Status `json:"status"`
}
