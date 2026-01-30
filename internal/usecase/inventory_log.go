package usecase

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/dto/response"
	"time"

	"go.uber.org/zap"
)

type InventoryLogService interface {
	GetInventoryLogs(ctx context.Context, req request.InventoryLogsFilter) (*response.PaginatedResponse[response.InventoryLogResponse], error)
	CreateInventoryLog(ctx context.Context, req request.CreateInventoryLogRequest) (*response.InventoryLogResponse, error)
}

type inventoryLogService struct {
	tx TxManager
	repo *repository.Repository
	log *zap.Logger
}

func NewInventoryLogService(tx TxManager, repo *repository.Repository, log *zap.Logger) InventoryLogService {
	return &inventoryLogService{
		tx: tx,
		repo: repo,
		log: log,
	}
}

func (s *inventoryLogService) GetInventoryLogs(ctx context.Context, req request.InventoryLogsFilter) (*response.PaginatedResponse[response.InventoryLogResponse], error) {
	// Construct params
	params := repository.InventoryLogParams{
		Offset: req.GetOffset(),
		Limit: req.GetPerPage(),
	}

	logs, total, err := s.repo.InventoryLogRepo.GetInventoryLogs(ctx, params)
	if err != nil {
		s.log.Error("Error get inventory logs service", zap.Error(err))
		return nil, err
	}

	// Convert to DTO
	var res []response.InventoryLogResponse
	for _, l := range logs {
		log := response.InventoryLogResponse{
			ID: l.ID,
			ProductID: l.ProductID,
			Type: l.Type,
			QuantityChange: l.QuantityChange,
			CurrentStockAfter: l.CurrentStockAfter,
			ReferenceID: l.ReferenceID,
			ReferenceType: l.ReferenceType,
			Notes: l.Notes,
			CreatedBy: l.CreatedBy,
			CreatedAt: l.CreatedAt,
		}
		res = append(res, log)
	}	

	return response.NewPaginatedResponse(
		res,
		req.GetPage(),
		req.GetPerPage(),
		total,
	), nil
}

func (s *inventoryLogService) CreateInventoryLog(ctx context.Context, req request.CreateInventoryLogRequest) (*response.InventoryLogResponse, error) {
	inventory := entity.InventoryLog{
		ProductID: req.ProductID,
		QuantityChange: req.QuantityChange,
		Notes: req.Notes,
		CreatedBy: ctx.Value("user_id").(uint),
		CreatedAt: time.Now(),
	}

	err := s.tx.WithinTx(ctx, func(ctx context.Context) error {
		// Get current stock
		// product, err := s.repo.ProductRepo.GetProductByID(ctx, req.ProductID)
		// if err != nil {
		// 	return err
		// }

		// // Construct inventory log
		// inventory.CurrentStockAfter = product.Stock + inventory.QuantityChange

		if req.Action == "restock" {
			inventory.Type = entity.InventoryLogTypeIn
			inventory.ReferenceType = "purchase"
		}

		if req.Action == "adjustment" {
			inventory.Type = entity.InventoryLogTypeAdjustment
			inventory.ReferenceType = "adjustment"
		}

		log, err := s.repo.InventoryLogRepo.CreateInventoryLog(ctx, &inventory)
		if err != nil {
			return err
		}
		inventory.ID = log.ID
		return nil
	})

	if err != nil {
		s.log.Error("Error create inventory log transaction", zap.Error(err))
		return nil, err
	}

	// Construct response
	res := response.InventoryLogResponse{
		ID: inventory.ID,
		ProductID: inventory.ProductID,
		Type: inventory.Type,
		QuantityChange: inventory.QuantityChange,
		CurrentStockAfter: inventory.CurrentStockAfter,
		ReferenceType: inventory.ReferenceType,
		Notes: inventory.Notes,
		CreatedBy: inventory.CreatedBy,
		CreatedAt: inventory.CreatedAt,
	}
		
	return &res, nil
}