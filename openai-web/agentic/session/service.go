package session

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"google.golang.org/adk/session"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var _ session.Service = &PostgresSessionService{}

type PostgresSessionService struct {
	db *gorm.DB
}

func NewPostgresSessionService(db *gorm.DB) (session.Service, error) {
	err := db.AutoMigrate(&SessionModel{}, &EventModel{})
	if err != nil {
		return nil, err
	}
	return &PostgresSessionService{db: db}, nil
}

func (s *PostgresSessionService) Create(ctx context.Context, req *session.CreateRequest) (*session.CreateResponse, error) {
	if req.AppName == "" || req.UserID == "" {
		return nil, fmt.Errorf("app_name and user_id are required")
	}

	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = uuid.NewString()
	}

	m := &SessionModel{
		SessionID:     sessionID,
		SessionName:   req.AppName,
		SessionUserID: req.UserID,
		SessionState:  datatypes.JSONMap(req.State),
		UpdatedAt:     time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, fmt.Errorf("session %s already exists or db error: %w", sessionID, err)
	}

	return &session.CreateResponse{Session: m}, nil
}

func (s *PostgresSessionService) Get(ctx context.Context, req *session.GetRequest) (*session.GetResponse, error) {
	var res SessionModel
	err := s.db.WithContext(ctx).
		Where("app_name = ? AND user_id = ? AND session_id = ?", req.AppName, req.UserID, req.SessionID).
		First(&res).Error

	if err != nil {
		return nil, err
	}

	var eventModels []EventModel
	query := s.db.WithContext(ctx).
		Where("app_name = ? AND user_id = ? AND session_id = ?", req.AppName, req.UserID, req.SessionID).
		Order("timestamp DESC")

	if !req.After.IsZero() {
		query = query.Where("timestamp >= ?", req.After)
	}
	if req.NumRecentEvents > 0 {
		query = query.Limit(req.NumRecentEvents)
	}

	if err := query.Find(&eventModels).Error; err != nil {
		return nil, err
	}

	for i := len(eventModels) - 1; i >= 0; i-- {
		var e session.Event
		json.Unmarshal(eventModels[i].EventData, &e)
	}

	return &session.GetResponse{Session: &res}, nil
}

func (s *PostgresSessionService) List(ctx context.Context, req *session.ListRequest) (*session.ListResponse, error) {
	var models []SessionModel
	query := s.db.WithContext(ctx).Where("app_name = ?", req.AppName)

	if req.UserID != "" {
		query = query.Where("user_id = ?", req.UserID)
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}

	var sessions []session.Session
	for _, m := range models {
		sessions = append(sessions, &m)
	}

	return &session.ListResponse{Sessions: sessions}, nil
}

func (s *PostgresSessionService) Delete(ctx context.Context, req *session.DeleteRequest) error {
	return s.db.WithContext(ctx).
		Where("app_name = ? AND user_id = ? AND session_id = ?", req.AppName, req.UserID, req.SessionID).
		Delete(&SessionModel{}).Error
}

func (s *PostgresSessionService) AppendEvent(ctx context.Context, curSession session.Session, event *session.Event) error {
	if event.Partial {
		return nil
	}

	delta := make(map[string]any)
	for k, v := range event.Actions.StateDelta {
		if !strings.HasPrefix(k, session.KeyPrefixTemp) {
			delta[k] = v
		}
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		eventJSON, _ := json.Marshal(event)
		eModel := EventModel{
			ID:        uuid.NewString(),
			SessionID: curSession.ID(),
			EventData: eventJSON,
			Timestamp: event.Timestamp,
		}
		if err := tx.Create(&eModel).Error; err != nil {
			return err
		}

		if len(delta) > 0 {
			deltaJSON, _ := json.Marshal(delta)
			err := tx.Exec(`
                UPDATE sessions 
                SET state = state || ?, updated_at = ? 
                WHERE app_name = ? AND user_id = ? AND session_id = ?`,
				deltaJSON, event.Timestamp, curSession.AppName(), curSession.UserID(), curSession.ID()).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}
