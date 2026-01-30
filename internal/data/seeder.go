package data

import (
	data "project-POS-APP-golang-integer/internal/data/seeder"

	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		users, err := data.UserSeeds(tx)
		if err != nil {
			return err
		}
		profiles, err := data.ProfileSeeds(tx, users)
		if err != nil {
			return err
		}
		err = data.ShiftSeeds(tx, profiles)
		if err != nil {
			return err
		}
		err = data.SessionSeeds(tx, users)
		if err != nil {
			return err
		}
		err = data.OTPSeeds(tx, users)
		if err != nil {
			return err
		}

		categories, err := data.CategorySeeds(tx)
		if err != nil {
			return err
		}
		products, err := data.ProductSeeds(tx, categories)
		if err != nil {
			return err
		}
		err = data.InventoryLogSeeds(tx, products)
		if err != nil {
			return err
		}

		customers, err := data.CustomerSeeds(db)
		if err != nil {
			return err
		}
		tables, err := data.TableSeeds(db)
		if err != nil {
			return err
		}
		err = data.ReservationSeeds(db, customers, tables)
		if err != nil {
			return err
		}
		orders, err := data.OrderSeeds(db, tables)
		if err != nil {
			return err
		}
		err = data.OrderItemSeeds(tx, products, orders)
		if err != nil {
			return err
		}

		pMethods, err := data.PaymentMethodSeeds(db)
		if err != nil {
			return err
		}
		err = data.TransactionSeeds(db, orders, pMethods)
		if err != nil {
			return err
		}
		err = data.NotificationSeeds(db, users)
		if err != nil {
			return err
		}

		return nil
	})
}
