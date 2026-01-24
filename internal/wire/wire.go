package wire

import (
	"project-POS-APP-golang-integer/internal/adaptor"
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/usecase"
	mCustom "project-POS-APP-golang-integer/pkg/middleware"
	"project-POS-APP-golang-integer/pkg/utils"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	Route *gin.Engine
	Stop chan struct{}
	WG *sync.WaitGroup
	Config utils.Configuration
}

func Wiring(repo *repository.Repository, log *zap.Logger, config utils.Configuration) *App {
	r := gin.Default()
	r1 := r.Group("/api/v1")

	stop := make(chan struct{})
	wg := &sync.WaitGroup{}

	usecase := usecase.NewUsecase(repo, log)
	handler := adaptor.NewHandler(usecase, log, config)
	mw := mCustom.NewMiddlewareCustom(usecase, log)
	ApiV1(r1, &handler, mw)

	return &App{
		Route: r,
		Stop: stop,
		WG: wg,
		Config: config,
	}
}

func ApiV1(r *gin.RouterGroup, handler *adaptor.Handler, mw mCustom.MiddlewareCustom) {
	users := r.Group("/users")
	users.GET("/", handler.UserHandler.GetListUsers)
}