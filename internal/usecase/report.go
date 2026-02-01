package usecase

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/dto/response"
	"time"

	"go.uber.org/zap"
)

type ReportService interface{
	GetSales(ctx context.Context, f request.PeriodRequest) (*response.SalesResponse, error)
	GetRevenue(ctx context.Context, f request.PeriodRequest) (*response.RevenueResponse, error)
	GetProductPerformance(ctx context.Context, f request.PeriodRequest) ([]response.ProductPerformance, error)
}

type reportService struct {
	tx   TxManager
	repo *repository.Repository
	log  *zap.Logger
}

func NewReportService(tx TxManager, repo *repository.Repository, log *zap.Logger) ReportService {
	return &reportService{
		tx:   tx,
		repo: repo,
		log:  log,
	}
}

func (s *reportService) GetSales(ctx context.Context, f request.PeriodRequest) (*response.SalesResponse, error) {
	var query request.PeriodQuery
	if f.From != "" {
		from, _ := time.Parse("02-01-2006", f.From)
		query.From = from
	} else {
		from, _ := time.Parse("02-01-2006", "01-01-2000")
		query.From = from
	}

	if f.To != "" {
		to, _ := time.Parse("02-01-2006", f.To)
		query.To = to
	} else {
		query.To = time.Now()
	}

	total, err := s.repo.ReportRepo.GetTotalSales(ctx, query)
	if err != nil {
		s.log.Error("Error get total sales", zap.Error(err))
		return nil, err
	}

	avgDaily, err := s.repo.ReportRepo.GetAverageDailySales(ctx, query)
	if err != nil {
		s.log.Error("Error get average daily sales", zap.Error(err))
		return nil, err
	}

	perHour, err := s.repo.ReportRepo.GetHourlySales(ctx, query)
	if err != nil {
		s.log.Error("Error get hourly sales", zap.Error(err))
		return nil, err
	}

	daily, err := s.repo.ReportRepo.GetDailySales(ctx, query)
	if err != nil {
		s.log.Error("Error get daily sales", zap.Error(err))
		return nil, err
	}

	monthly, err := s.repo.ReportRepo.GetMonthlySales(ctx, query)
	if err != nil {
		s.log.Error("Error get monthly sales", zap.Error(err))
		return nil, err
	}

	res := response.SalesResponse{
		Summary: response.SalesSummary{
			Total: total,
			AverageDaily: avgDaily,
		},
		PerHour: perHour,
		Daily: daily,
		Monthly: monthly,
	}

	return &res, nil
}

func (s *reportService) GetRevenue(ctx context.Context, f request.PeriodRequest) (*response.RevenueResponse, error) {
	var query request.PeriodQuery
	if f.From != "" {
		from, _ := time.Parse("02-01-2006", f.From)
		query.From = from
	} else {
		from, _ := time.Parse("02-01-2006", "01-01-2000")
		query.From = from
	}

	if f.To != "" {
		to, _ := time.Parse("02-01-2006", f.To)
		query.To = to
	} else {
		query.To = time.Now()
	}

	total, err := s.repo.ReportRepo.GetTotalRevenue(ctx, query)
	if err != nil {
		s.log.Error("Error get total revenue", zap.Error(err))
		return nil, err
	}

	avgDaily, err := s.repo.ReportRepo.GetAverageDailyRevenue(ctx, query)
	if err != nil {
		s.log.Error("Error get average daily revenue", zap.Error(err))
		return nil, err
	}

	perHour, err := s.repo.ReportRepo.GetHourlyRevenue(ctx, query)
	if err != nil {
		s.log.Error("Error get hourly revenue", zap.Error(err))
		return nil, err
	}

	daily, err := s.repo.ReportRepo.GetDailyRevenue(ctx, query)
	if err != nil {
		s.log.Error("Error get daily revenue", zap.Error(err))
		return nil, err
	}

	monthly, err := s.repo.ReportRepo.GetMonthlyRevenue(ctx, query)
	if err != nil {
		s.log.Error("Error get monthly revenue", zap.Error(err))
		return nil, err
	}

	res := response.RevenueResponse{
		Summary: response.RevenueSummary{
			Total: total,
			AverageDaily: avgDaily,
		},
		PerHour: perHour,
		Daily: daily,
		Monthly: monthly,
	}

	return &res, nil
}

func (s *reportService) GetProductPerformance(ctx context.Context, f request.PeriodRequest) ([]response.ProductPerformance, error) {
	var query request.PeriodQuery
	if f.From != "" {
		from, _ := time.Parse("02-01-2006", f.From)
		query.From = from
	} else {
		from, _ := time.Parse("02-01-2006", "01-01-2000")
		query.From = from
	}

	if f.To != "" {
		to, _ := time.Parse("02-01-2006", f.To)
		query.To = to
	} else {
		query.To = time.Now()
	}
	
	result, err := s.repo.ReportRepo.GetProductPerformance(ctx, query)
	if err != nil {
		s.log.Error("Error get product performance", zap.Error(err))
		return nil, err
	}
	return result, nil
}