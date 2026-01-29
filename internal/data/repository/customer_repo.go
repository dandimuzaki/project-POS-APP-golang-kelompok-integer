package repository

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/internal/dto/request"
	"project-POS-APP-golang-integer/internal/infra"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *entity.Customer) (*entity.Customer, error)
	FindByID(ctx context.Context, id uint) (*entity.Customer, error)
	FindByPhone(ctx context.Context, phone string) (*entity.Customer, error)
	FindByEmail(ctx context.Context, email string) (*entity.Customer, error)
	FindAll(ctx context.Context, params request.GetCustomersRequest) ([]entity.Customer, int64, error)
	Update(ctx context.Context, customer *entity.Customer) error
	Delete(ctx context.Context, id uint) error
}

type customerRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewCustomerRepo(db *gorm.DB, log *zap.Logger) CustomerRepository {
	return &customerRepository{
		db:     db,
		logger: log.With(zap.String("repository", "customer")),
	}
}

func (r *customerRepository) Create(ctx context.Context, customer *entity.Customer) (*entity.Customer, error) {
	db := infra.GetDB(ctx, r.db)

	r.logger.Info("Creating customer",
		zap.String("phone", customer.Phone),
		zap.String("name", customer.FirstName))

	err := db.Create(customer).Error
	if err != nil {
		r.logger.Error("Failed to create customer",
			zap.String("phone", customer.Phone),
			zap.Error(err))
		return nil, err
	}

	r.logger.Info("Customer created successfully",
		zap.Uint("id", customer.ID),
		zap.String("phone", customer.Phone))

	return customer, nil
}

func (r *customerRepository) FindByID(ctx context.Context, id uint) (*entity.Customer, error) {
	db := infra.GetDB(ctx, r.db)

	r.logger.Debug("Finding customer by ID", zap.Uint("id", id))

	var customer entity.Customer
	err := db.First(&customer, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warn("Customer not found", zap.Uint("id", id))
		} else {
			r.logger.Error("Failed to find customer",
				zap.Uint("id", id),
				zap.Error(err))
		}
		return nil, err
	}

	r.logger.Debug("Customer found", zap.Uint("id", id))
	return &customer, nil
}

func (r *customerRepository) FindByPhone(ctx context.Context, phone string) (*entity.Customer, error) {
	db := infra.GetDB(ctx, r.db)

	r.logger.Debug("Finding customer by phone", zap.String("phone", phone))

	var customer entity.Customer
	err := db.Where("phone = ?", phone).First(&customer).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("Customer not found by phone", zap.String("phone", phone))
		} else {
			r.logger.Error("Failed to find customer by phone",
				zap.String("phone", phone),
				zap.Error(err))
		}
		return nil, err
	}

	r.logger.Debug("Customer found by phone",
		zap.String("phone", phone),
		zap.Uint("id", customer.ID))

	return &customer, nil
}

func (r *customerRepository) FindByEmail(ctx context.Context, email string) (*entity.Customer, error) {
	db := infra.GetDB(ctx, r.db)

	r.logger.Debug("Finding customer by email", zap.String("email", email))

	var customer entity.Customer
	err := db.Where("email = ?", email).First(&customer).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Debug("Customer not found by email", zap.String("email", email))
		} else {
			r.logger.Error("Failed to find customer by email",
				zap.String("email", email),
				zap.Error(err))
		}
		return nil, err
	}

	return &customer, nil
}

func (r *customerRepository) FindAll(ctx context.Context, params request.GetCustomersRequest) ([]entity.Customer, int64, error) {
	db := infra.GetDB(ctx, r.db)

	r.logger.Debug("Finding customers",
		zap.String("search", params.Search),
		zap.Int("page", params.GetPage()),
		zap.Int("per_page", params.GetPerPage()))

	var customers []entity.Customer
	var total int64

	query := db.Model(&entity.Customer{})

	// Apply search filter
	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR phone ILIKE ? OR email ILIKE ?",
			search, search, search, search)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		r.logger.Error("Failed to count customers", zap.Error(err))
		return nil, 0, err
	}

	// Apply pagination
	offset := params.GetOffset()
	limit := params.GetPerPage()

	err := query.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&customers).Error

	if err != nil {
		r.logger.Error("Failed to find customers",
			zap.String("search", params.Search),
			zap.Error(err))
		return nil, 0, err
	}

	r.logger.Debug("Customers found",
		zap.Int("count", len(customers)),
		zap.Int64("total", total))

	return customers, total, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *entity.Customer) error {
	db := infra.GetDB(ctx, r.db)

	r.logger.Info("Updating customer",
		zap.Uint("id", customer.ID),
		zap.String("phone", customer.Phone))

	err := db.Save(customer).Error
	if err != nil {
		r.logger.Error("Failed to update customer",
			zap.Uint("id", customer.ID),
			zap.Error(err))
		return err
	}

	r.logger.Info("Customer updated", zap.Uint("id", customer.ID))
	return nil
}

func (r *customerRepository) Delete(ctx context.Context, id uint) error {
	db := infra.GetDB(ctx, r.db)

	r.logger.Info("Deleting customer", zap.Uint("id", id))

	err := db.Delete(&entity.Customer{}, id).Error
	if err != nil {
		r.logger.Error("Failed to delete customer",
			zap.Uint("id", id),
			zap.Error(err))
		return err
	}

	r.logger.Info("Customer deleted", zap.Uint("id", id))
	return nil
}
