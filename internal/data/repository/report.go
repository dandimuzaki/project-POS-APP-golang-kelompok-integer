package repository

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/dto/response"
	"project-POS-APP-golang-integer/internal/infra"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ReportRepository interface{
	GetTotalSales(ctx context.Context, f request.PeriodQuery) (int64, error)
	GetAverageDailySales(ctx context.Context, f request.PeriodQuery) (float64, error)
	GetHourlySales(ctx context.Context, f request.PeriodQuery) ([]response.PerHourSales, error)
	GetDailySales(ctx context.Context, f request.PeriodQuery) ([]response.DailySales, error)
	GetMonthlySales(ctx context.Context, f request.PeriodQuery) ([]response.MonthlySales, error)
	GetTotalRevenue(ctx context.Context, f request.PeriodQuery) (float64, error)
	GetAverageDailyRevenue(ctx context.Context, f request.PeriodQuery) (float64, error)
	GetHourlyRevenue(ctx context.Context, f request.PeriodQuery) ([]response.PerHourRevenue, error)
	GetDailyRevenue(ctx context.Context, f request.PeriodQuery) ([]response.DailyRevenue, error)
	GetMonthlyRevenue(ctx context.Context, f request.PeriodQuery) ([]response.MonthlyRevenue, error)
}

type reportRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewReportRepo(db *gorm.DB, log *zap.Logger) ReportRepository {
	return &reportRepository{
		db:  db,
		log: log,
	}
}

func (r *reportRepository) GetTotalSales(ctx context.Context, f request.PeriodQuery) (int64, error) {
	db := infra.GetDB(ctx, r.db)
	var total int64
	err := db.Model(&entity.Order{}).
		Select("sum(oi.quantity) as total").
		Joins("JOIN order_items oi ON oi.order_id = orders.id").
		Where("orders.created_at > ? AND orders.created_at < ? AND orders.status = ?", f.From, f.To, entity.OrderStatusCompleted).
		Scan(&total).Error

	if err != nil {
		r.log.Error("Error get total sales", zap.Error(err))
		return 0, err
	}

	return total, nil
}

func (r *reportRepository) GetAverageDailySales(ctx context.Context, f request.PeriodQuery) (float64, error) {
	db := infra.GetDB(ctx, r.db)
	var avg float64
	err := db.Model(&entity.Order{}).
		Select("avg(oi.quantity) as average").
		Joins("JOIN order_items oi ON oi.order_id = orders.id").
		Where("orders.created_at > ? AND orders.created_at < ? AND orders.status = ?", f.From, f.To, entity.OrderStatusCompleted).
		Group("TO_CHAR(orders.created_at, 'DD-MM-YYYY')").
		Scan(&avg).Error

	if err != nil {
		r.log.Error("Error get average daily sales", zap.Error(err))
		return 0, err
	}

	return avg, nil
}

func (r *reportRepository) GetHourlySales(ctx context.Context, f request.PeriodQuery) ([]response.PerHourSales, error) {
	db := infra.GetDB(ctx, r.db)
	var perHour []response.PerHourSales
	err := db.Model(&entity.Order{}).
		Select("TO_CHAR(DATE_TRUNC('hour', orders.created_at), 'HH:00') AS hour, avg(oi.quantity) as average_sales").
		Joins("JOIN order_items oi ON oi.order_id = orders.id").
		Where("orders.created_at > ? AND orders.created_at < ? AND orders.status = ?", f.From, f.To, entity.OrderStatusCompleted).
		Group("hour").
		Order("hour desc").
		Scan(&perHour).Error

	if err != nil {
		r.log.Error("Error get hourly sales", zap.Error(err))
		return nil, err
	}

	return perHour, nil
}

func (r *reportRepository) GetDailySales(ctx context.Context, f request.PeriodQuery) ([]response.DailySales, error) {
	db := infra.GetDB(ctx, r.db)
	var daily []response.DailySales
	err := db.Model(&entity.Order{}).
		Select("TO_CHAR(orders.created_at, 'DD-MM-YYYY') as date, sum(oi.quantity) as sales").
		Joins("JOIN order_items oi ON oi.order_id = orders.id").
		Where("orders.created_at > ? AND orders.created_at < ? AND orders.status = ?", f.From, f.To, entity.OrderStatusCompleted).
		Group("date").
		Order("date desc").
		Scan(&daily).Error

	if err != nil {
		r.log.Error("Error get daily sales", zap.Error(err))
		return nil, err
	}

	return daily, nil
}

func (r *reportRepository) GetMonthlySales(ctx context.Context, f request.PeriodQuery) ([]response.MonthlySales, error) {
	db := infra.GetDB(ctx, r.db)
	var results []response.MonthlySales
	err := db.Model(&entity.Order{}).
		Select("TO_CHAR(orders.created_at, 'YYYY-MM') as month, sum(oi.quantity) as sales").
		Joins("JOIN order_items oi ON oi.order_id = orders.id").
		Where("orders.created_at > ? AND orders.created_at < ? AND orders.status = ?", f.From, f.To, entity.OrderStatusCompleted).
		Group("month").
		Order("month desc").
		Scan(&results).Error
	
	if err != nil {
		r.log.Error("Error get monthly sales", zap.Error(err))
		return nil, err
	}

	return results, err
}

func (r *reportRepository) GetTotalRevenue(ctx context.Context, f request.PeriodQuery) (float64, error) {
	db := infra.GetDB(ctx, r.db)
	var total float64
	err := db.Model(&entity.Order{}).
		Select("sum(total) as total").
		Where("created_at > ? AND created_at < ? AND status = ?", f.From, f.To, entity.OrderStatusCompleted).
		Scan(&total).Error

	if err != nil {
		r.log.Error("Error get total revenue", zap.Error(err))
		return 0, err
	}

	return total, nil
}

func (r *reportRepository) GetAverageDailyRevenue(ctx context.Context, f request.PeriodQuery) (float64, error) {
	db := infra.GetDB(ctx, r.db)
	var avg float64
	err := db.Model(&entity.Order{}).
		Select("avg(total) as average").
		Where("created_at > ? AND created_at < ? AND status = ?", f.From, f.To, entity.OrderStatusCompleted).
		Group("TO_CHAR(created_at, 'DD-MM-YYYY')").
		Scan(&avg).Error

	if err != nil {
		r.log.Error("Error get average daily revenue", zap.Error(err))
		return 0, err
	}

	return avg, nil
}

func (r *reportRepository) GetHourlyRevenue(ctx context.Context, f request.PeriodQuery) ([]response.PerHourRevenue, error) {
	db := infra.GetDB(ctx, r.db)
	var perHour []response.PerHourRevenue
	err := db.Model(&entity.Order{}).
		Select("TO_CHAR(DATE_TRUNC('hour', created_at), 'HH:00') AS hour, avg(total) as average_revenue").
		Where("created_at > ? AND created_at < ? AND status = ?", f.From, f.To, entity.OrderStatusCompleted).
		Group("hour").
		Order("hour desc").
		Scan(&perHour).Error

	if err != nil {
		r.log.Error("Error get hourly revenue", zap.Error(err))
		return nil, err
	}

	return perHour, nil
}

func (r *reportRepository) GetDailyRevenue(ctx context.Context, f request.PeriodQuery) ([]response.DailyRevenue, error) {
	db := infra.GetDB(ctx, r.db)
	var daily []response.DailyRevenue
	err := db.Model(&entity.Order{}).
		Select("TO_CHAR(created_at, 'DD-MM-YYYY') as date, sum(total) as revenue").
		Where("created_at > ? AND created_at < ? AND status = ?", f.From, f.To, entity.OrderStatusCompleted).
		Group("date").
		Order("date desc").
		Scan(&daily).Error

	if err != nil {
		r.log.Error("Error get daily revenue", zap.Error(err))
		return nil, err
	}

	return daily, nil
}

func (r *reportRepository) GetMonthlyRevenue(ctx context.Context, f request.PeriodQuery) ([]response.MonthlyRevenue, error) {
	db := infra.GetDB(ctx, r.db)
	var results []response.MonthlyRevenue
	err := db.Model(&entity.Order{}).
		Select("TO_CHAR(created_at, 'YYYY-MM') as month, sum(total) as revenue").
		Where("created_at > ? AND created_at < ? AND status = ?", f.From, f.To, entity.OrderStatusCompleted).
		Group("month").
		Order("month desc").
		Scan(&results).Error
	
	if err != nil {
		r.log.Error("Error get monthly revenue", zap.Error(err))
		return nil, err
	}

	return results, err
}