package wire

import (
	"project-POS-APP-golang-integer/internal/adaptor"
	"project-POS-APP-golang-integer/internal/data/repository"
	"project-POS-APP-golang-integer/internal/infra"
	"project-POS-APP-golang-integer/internal/usecase"
	mCustom "project-POS-APP-golang-integer/pkg/middleware"
	"project-POS-APP-golang-integer/pkg/utils"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	Route  *gin.Engine
	Stop   chan struct{}
	WG     *sync.WaitGroup
	Config utils.Configuration
}

func Wiring(db *gorm.DB, repo *repository.Repository, log *zap.Logger, config utils.Configuration) *App {
	r := gin.Default()
	r1 := r.Group("/api/v1")

	stop := make(chan struct{})
	wg := &sync.WaitGroup{}

	tx := infra.NewGormTxManager(db)
	usecase := usecase.NewUsecase(tx, db, repo, log, config)
	handler := adaptor.NewHandler(usecase, log, config)
	mw := mCustom.NewMiddlewareCustom(usecase, log)

	// Panggil ApiV1 dengan semua routes termasuk category
	ApiV1(r1, &handler, mw)

	return &App{
		Route:  r,
		Stop:   stop,
		WG:     wg,
		Config: config,
	}
}

func ApiV1(r *gin.RouterGroup, handler *adaptor.Handler, mw mCustom.MiddlewareCustom) {
	AuthRoute(r.Group("/auth"), handler, mw)
	UserRoute(r.Group("/users"), handler, mw)
	ProfileRoute(r.Group("/profile"), handler, mw)
	ReservationRoute(r.Group("/reservations"), handler, mw)
	InventoryRoute(r.Group("/inventories"), handler, mw)
	CategoryRoute(r.Group("/categories"), handler, mw)
	ProductRoute(r.Group("/products"), handler, mw)
}

func AuthRoute(r *gin.RouterGroup, handler *adaptor.Handler, mw mCustom.MiddlewareCustom) {
	r.POST("/login", handler.AuthHandler.Login)
	r.POST("/logout", mw.AuthMiddleware(), handler.AuthHandler.Logout)
	r.POST("/request-reset-password", mw.AuthMiddleware(), handler.AuthHandler.RequestResetPassword)
	r.POST("/reset-password", mw.AuthMiddleware(), handler.AuthHandler.ResetPassword)
}

func UserRoute(r *gin.RouterGroup, handler *adaptor.Handler, mw mCustom.MiddlewareCustom) {
	r.GET("/", mw.AuthMiddleware(), mw.RequirePermission("superadmin", "admin"), handler.UserHandler.GetUserList)
	r.POST("/", mw.AuthMiddleware(), mw.RequirePermission("superadmin", "admin"), handler.UserHandler.CreateUser)
	r.PUT("/:id", mw.AuthMiddleware(), mw.RequirePermission("superadmin", "admin"), handler.UserHandler.UpdateRole)
	r.DELETE("/:id", mw.AuthMiddleware(), mw.RequirePermission("superadmin", "admin"), handler.UserHandler.DeleteUser)
}

func ProfileRoute(r *gin.RouterGroup, handler *adaptor.Handler, mw mCustom.MiddlewareCustom) {
	r.GET("/", mw.AuthMiddleware(), handler.ProfileHandler.GetProfile)
	r.PUT("/", mw.AuthMiddleware(), handler.ProfileHandler.UpdateProfile)
}

func ReservationRoute(r *gin.RouterGroup, handler *adaptor.Handler, mw mCustom.MiddlewareCustom) {
	// Public routes
	r.GET("/available-tables", handler.ReservationHandler.GetAvailableTables)

	// Protected routes (need authentication)
	r.Use(mw.AuthMiddleware())

	r.POST("/", handler.ReservationHandler.CreateReservation)
	r.GET("/", handler.ReservationHandler.GetReservations)
	r.GET("/:id", handler.ReservationHandler.GetReservationByID)
	r.PUT("/:id/status", handler.ReservationHandler.UpdateReservationStatus)
	r.POST("/:id/cancel", handler.ReservationHandler.CancelReservation)
	r.POST("/:id/checkin", handler.ReservationHandler.CheckIn)
}

func InventoryRoute(r *gin.RouterGroup, handler *adaptor.Handler, mw mCustom.MiddlewareCustom) {
	r.GET("/", mw.AuthMiddleware(), handler.InventoryLogHandler.GetInventoryLogs)
	r.POST("/", mw.AuthMiddleware(), handler.InventoryLogHandler.CreateInventoryLog)
}

func CategoryRoute(r *gin.RouterGroup, handler *adaptor.Handler, mw mCustom.MiddlewareCustom) {
	// 🔥 PUBLIC ROUTES - tidak perlu auth
	r.GET("", handler.CategoryHandler.GetAllCategories)    // GET /api/v1/categories
	r.GET("/:id", handler.CategoryHandler.GetCategoryByID) // GET /api/v1/categories/:id
	// 🔥 PROTECTED ROUTES - perlu auth DAN admin permission
	protected := r.Group("")
	protected.Use(
		mw.AuthMiddleware(),
		mw.RequirePermission("admin", "superadmin"), // 🔥 GUNAKAN middleware yang ADA
	)
	protected.POST("", handler.CategoryHandler.CreateCategory)       // POST /api/v1/categories
	protected.PUT("/:id", handler.CategoryHandler.UpdateCategory)    // PUT /api/v1/categories/:id
	protected.DELETE("/:id", handler.CategoryHandler.DeleteCategory) // DELETE /api/v1/categories/:id
}

func ProductRoute(r *gin.RouterGroup, handler *adaptor.Handler, mw mCustom.MiddlewareCustom) {
	// 🔥 PUBLIC ROUTES - tidak perlu auth
	r.GET("", handler.ProductHandler.GetAllProducts) // TODO: akan dibuat nanti
	r.GET("/:id", handler.ProductHandler.GetProductByID)

	// 🔥 PROTECTED ROUTES - perlu auth DAN admin permission
	protected := r.Group("")
	protected.Use(
		mw.AuthMiddleware(),
		mw.RequirePermission("admin", "superadmin"),
	)

	protected.POST("", handler.ProductHandler.CreateProduct)
	protected.PUT("/:id", handler.ProductHandler.UpdateProduct)
	protected.DELETE("/:id", handler.ProductHandler.DeleteProduct)
}
