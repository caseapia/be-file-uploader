package roadmap

import (
	"be-file-uploader/internal/repository/mysql"
	"be-file-uploader/internal/service/roadmap"
	"be-file-uploader/pkg/utils/validation"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	roadmapService *roadmap.Service
	repository     *mysql.Repository
}

func NewHandler(roadmapService *roadmap.Service, repository *mysql.Repository) *Handler {
	return &Handler{roadmapService: roadmapService, repository: repository}
}

func (h *Handler) GetRoadmapList(ctx fiber.Ctx) error {
	rmaplist, err := h.repository.GetRoadmapList(ctx.Context())
	if err != nil {
		return err
	}

	return validation.Response(ctx, 200, rmaplist)
}
