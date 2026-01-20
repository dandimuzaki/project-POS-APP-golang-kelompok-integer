package adaptor

import (
	"net/http"
	"strconv"
	"travel-api/internal/dto"
	"travel-api/internal/usecase"
	"travel-api/pkg/utils"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type TourHandler struct {
	service *usecase.Usecase
	Logger *zap.Logger
	Config utils.Configuration
}

func NewTourAdaptor(service *usecase.Usecase, log *zap.Logger, config utils.Configuration) TourHandler {
	return TourHandler{
		service: service,
		Logger: log,
		Config: config,
	}
}

// GetAllTours menangani GET /tours
func (h *TourHandler) GetAllTours(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	date := r.URL.Query().Get("date")
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sort_by")

	// Construct DTO
	req := dto.TourFilterRequest{
		Date:   date,
		Search: search,
		SortBy: sortBy,
	}

	result, pagination, err := h.service.TourService.GetListTours(ctx, req)
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadGateway, "", nil)
		return
	}

	utils.ResponsePagination(w, http.StatusOK, "success get data", result, pagination)
}

func (h *TourHandler) GetTourDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	result, err := h.service.TourService.ScheduleByID(ctx, uint(id))
	if err != nil {
		utils.ResponseBadRequest(w, http.StatusBadGateway, "", nil)
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "success get data", result)
}