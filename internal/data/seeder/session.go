package data

import (
	"project-POS-APP-golang-integer/internal/data/entity"
	"project-POS-APP-golang-integer/pkg/utils"
	"time"
)

func SessionSeeds() []entity.Session{
	token1, _ := utils.GenerateRandomToken(16)
	token2, _ := utils.GenerateRandomToken(16)
	token3, _ := utils.GenerateRandomToken(16)
	return []entity.Session{
		{
			UserID: 1,
			Token: token1.String(),
			ExpiresAt: time.Now().AddDate(0,0,5),
			CreatedAt: time.Now(),
		},
		{
			UserID: 2,
			Token: token2.String(),
			ExpiresAt: time.Now().AddDate(0,0,5),
			CreatedAt: time.Now(),
		},
		{
			UserID: 3,
			Token: token3.String(),
			ExpiresAt: time.Now().AddDate(0,0,5),
			CreatedAt: time.Now(),
		},
	}
}