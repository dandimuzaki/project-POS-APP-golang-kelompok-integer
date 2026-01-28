package repository

import (
	"project-POS-APP-golang-integer/internal/data/entity"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(customer *entity.Customer, tx ...*gorm.DB) error
	FindByID(id uint, tx ...*gorm.DB) (*entity.Customer, error)
	FindByPhone(phone string, tx ...*gorm.DB) (*entity.Customer, error)
	FindByEmail(email string, tx ...*gorm.DB) (*entity.Customer, error)
	FindAll(params CustomerQueryParams, tx ...*gorm.DB) ([]entity.Customer, int64, error)
	Update(customer *entity.Customer, tx ...*gorm.DB) error
	Delete(id uint, tx ...*gorm.DB) error
}

type customerRepository struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewCustomerRepository(db *gorm.DB, log *zap.Logger) CustomerRepository {
	return &customerRepository{
		db:  db,
		log: log.With(zap.String("repository", "customer")),
	}
}

type CustomerQueryParams struct {
	Search string // Search by name, phone, or email
	Offset int    // Pagination offset
	Limit  int    // Pagination limit
}

// getDB returns either transaction db or regular db
func (cr *customerRepository) getDB(tx ...*gorm.DB) *gorm.DB {
	if len(tx) > 0 && tx[0] != nil {
		return tx[0]
	}
	return cr.db
}

func (cr *customerRepository) Create(customer *entity.Customer, tx ...*gorm.DB) error {
	db := cr.getDB(tx...)

	cr.log.Debug("Creating new customer",
		zap.String("phone", customer.Phone))

	err := db.Create(customer).Error
	if err != nil {
		cr.log.Error("Failed to create customer",
			zap.String("phone", customer.Phone),
			zap.Error(err))
		return err
	}

	cr.log.Info("Customer created successfully",
		zap.Uint("id", customer.ID),
		zap.String("phone", customer.Phone))
	return nil
}

func (cr *customerRepository) FindByID(id uint, tx ...*gorm.DB) (*entity.Customer, error) {
	db := cr.getDB(tx...)

	cr.log.Debug("Finding customer by ID", zap.Uint("id", id))

	var customer entity.Customer
	err := db.First(&customer, id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			cr.log.Warn("Customer not found", zap.Uint("id", id))
		} else {
			cr.log.Error("Failed to find customer by ID",
				zap.Uint("id", id),
				zap.Error(err))
		}
		return nil, err
	}

	cr.log.Debug("Customer found", zap.Uint("id", id))
	return &customer, nil
}

func (cr *customerRepository) FindByPhone(phone string, tx ...*gorm.DB) (*entity.Customer, error) {
	db := cr.getDB(tx...)

	cr.log.Debug("Finding customer by phone", zap.String("phone", phone))

	var customer entity.Customer
	err := db.Where("phone = ?", phone).First(&customer).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			cr.log.Debug("Customer not found by phone", zap.String("phone", phone))
		} else {
			cr.log.Error("Failed to find customer by phone",
				zap.String("phone", phone),
				zap.Error(err))
		}
		return nil, err
	}

	cr.log.Debug("Customer found by phone",
		zap.String("phone", phone),
		zap.Uint("id", customer.ID))
	return &customer, nil
}

func (cr *customerRepository) FindByEmail(email string, tx ...*gorm.DB) (*entity.Customer, error) {
	db := cr.getDB(tx...)

	cr.log.Debug("Finding customer by email", zap.String("email", email))

	var customer entity.Customer
	err := db.Where("email = ?", email).First(&customer).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			cr.log.Debug("Customer not found by email", zap.String("email", email))
		} else {
			cr.log.Error("Failed to find customer by email",
				zap.String("email", email),
				zap.Error(err))
		}
		return nil, err
	}

	return &customer, nil
}

func (cr *customerRepository) FindAll(params CustomerQueryParams, tx ...*gorm.DB) ([]entity.Customer, int64, error) {
	db := cr.getDB(tx...)

	cr.log.Debug("Finding customers",
		zap.String("search", params.Search),
		zap.Int("offset", params.Offset),
		zap.Int("limit", params.Limit))

	var customers []entity.Customer
	var total int64

	query := db.Model(&entity.Customer{})

	// Apply search filter if provided
	if params.Search != "" {
		search := "%" + params.Search + "%"
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR phone ILIKE ? OR email ILIKE ?",
			search, search, search, search)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		cr.log.Error("Failed to count customers", zap.Error(err))
		return nil, 0, err
	}

	// Apply pagination
	if params.Limit > 0 {
		query = query.Offset(params.Offset).Limit(params.Limit)
	}

	// Execute query with ordering
	if err := query.Order("created_at DESC").Find(&customers).Error; err != nil {
		cr.log.Error("Failed to find customers",
			zap.String("search", params.Search),
			zap.Error(err))
		return nil, 0, err
	}

	cr.log.Debug("Customers found",
		zap.Int("count", len(customers)),
		zap.Int64("total", total))
	return customers, total, nil
}

func (cr *customerRepository) Update(customer *entity.Customer, tx ...*gorm.DB) error {
	db := cr.getDB(tx...)

	cr.log.Debug("Updating customer",
		zap.Uint("id", customer.ID),
		zap.String("phone", customer.Phone))

	err := db.Save(customer).Error
	if err != nil {
		cr.log.Error("Failed to update customer",
			zap.Uint("id", customer.ID),
			zap.Error(err))
		return err
	}

	cr.log.Info("Customer updated", zap.Uint("id", customer.ID))
	return nil
}

func (cr *customerRepository) Delete(id uint, tx ...*gorm.DB) error {
	db := cr.getDB(tx...)

	cr.log.Debug("Deleting customer", zap.Uint("id", id))

	err := db.Delete(&entity.Customer{}, id).Error
	if err != nil {
		cr.log.Error("Failed to delete customer",
			zap.Uint("id", id),
			zap.Error(err))
		return err
	}

	cr.log.Info("Customer deleted", zap.Uint("id", id))
	return nil
}
