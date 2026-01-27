package repository

// type SessionRepository interface {
// 	Create(ctx context.Context, userID int) (uuid.UUID, error)
// 	Revoke(ctx context.Context, token string) error
// 	ValidateToken(ctx context.Context, token string) (*int, error)
// }

// type sessionRepository struct {
// 	db *gorm.DB
// 	Logger *zap.Logger
// }

// func NewSessionRepo(db *gorm.DB, log *zap.Logger) SessionRepository {
// 	return &sessionRepository{
// 		db: db,
// 		Logger: log,
// 	}
// }

// func (r *sessionRepository) Create(ctx context.Context, userID int) (uuid.UUID, error) {
// 	// Create session after login and register
// 	token, err := utils.GenerateRandomToken(16)
// 	if err != nil {
// 		r.Logger.Error("Error create token: ", zap.Error(err))
// 		return uuid.Nil, err
// 	}

// 	session := entity.Session{
// 		UserID: uint(userID),
// 		Token: token,
// 		ExpiresAt: time.Now().AddDate(0,0,5),
// 		CreatedAt: time.Now(),
// 	}

// 	query := `INSERT INTO sessions (token, user_id, expired_at) VALUES ($1, $2, $3) RETURNING token`
// 	err = r.db.QueryRow(context.Background(), query, token, userID, time.Now().AddDate(0,1,0)).Scan(&token)
// 	if err != nil {
// 		r.Logger.Error("Error query create session: ", zap.Error(err))
// 		return uuid.Nil, err
// 	}

// 	return token, nil
// }

// func (r *sessionRepository) ValidateToken(ctx context.Context, token string) (*int, error) {
// 	// Validate token to authorize user
// 	var userID *int
// 	query := `SELECT user_id FROM sessions WHERE token = $1 AND expired_at > NOW() AND revoked_at IS NULL`
// 	err := r.db.QueryRow(context.Background(), query, token).Scan(&userID)
// 	if err != nil {
// 		r.Logger.Error("Error query validate token: ", zap.Error(err))
// 		return nil, err
// 	}

// 	return userID, nil
// }

// func (r *sessionRepository) Revoke(ctx context.Context, token string) (error) {
// 	// Revoke session after logout
// 	query := `UPDATE sessions SET revoked_at = NOW() WHERE token = $1`
// 	result, err := r.db.Exec(context.Background(), query, token)
// 	rowsAffected := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		r.Logger.Error("Error not found session: ", zap.Error(err))
// 		return nil
// 	}

// 	if err != nil {
// 		r.Logger.Error("Error query revoke session: ", zap.Error(err))
// 	}

// 	return err
// }