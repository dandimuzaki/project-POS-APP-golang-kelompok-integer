package mocks

import "project-POS-APP-golang-integer/internal/data/repository"

type Repository struct {
	UserRepo repository.UserRepository
	ProfileRepo repository.ProfileRepository
}
