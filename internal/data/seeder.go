package data

import (
	"fmt"
	data "project-POS-APP-golang-integer/internal/data/seeder"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedAll(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		seeds := dataSeeds()
		for i := range seeds {
			err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(seeds[i]).Error
			if err != nil {
				name := reflect.TypeOf(seeds[i]).String()
				errMessage := err.Error()
				return fmt.Errorf("%s seeder failed with %s", name, errMessage)
			}
		}
		return nil
	})
}

func dataSeeds() []interface{} {
	return []interface{}{
		data.UserSeeds(),
		data.ProfileSeeds(),
		data.ShiftSeeds(),
		data.OTPSeeds(),
		data.SessionSeeds(),
		data.CategorySeeds(),
		data.ProductSeeds(),
		data.InventoryLogSeeds(),
		data.TableSeeds(),
		data.CustomerSeeds(),
		data.OrderSeeds(),
		data.OrderItemSeeds(),
		data.ReservationSeeds(),
		data.PaymentMethodSeeds(),
		data.TransactionSeeds(),
		data.NotificationSeeds(),
	}
}