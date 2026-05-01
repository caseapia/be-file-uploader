package roadmap

import (
	"be-file-uploader/internal/models/requests"
	"be-file-uploader/internal/service/roadmap"
	"be-file-uploader/pkg/utils/account"
	"be-file-uploader/pkg/utils/validation"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	roadmapService *roadmap.Service
}

func NewHandler(roadmapService *roadmap.Service) *Handler {
	return &Handler{roadmapService: roadmapService}
}

func (h *Handler) GetRoadmap(ctx fiber.Ctx) error {
	list, err := h.roadmapService.RoadmapList(ctx)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, &fiber.Map{"roadmap": list})
}

func (h *Handler) AddRoadmapTask(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)
	var req requests.RoadmapAddRequest
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	task, err := h.roadmapService.AddTask(ctx, sender, req.Title)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusCreated, task)
}

func (h *Handler) EditRoadmapTask(ctx fiber.Ctx) error {
	sender := account.GetUserFromContext(ctx)
	var req requests.RoadmapUpdateRequest
	if err := validation.ParseAndValidate(ctx, &req); err != nil {
		return err
	}

	task, err := h.roadmapService.EditTask(ctx, sender, req.ID, req.Title, req.Status)
	if err != nil {
		return err
	}

	return validation.Response(ctx, fiber.StatusOK, task)
}
