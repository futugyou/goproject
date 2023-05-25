package session

import (
	"context"
	"sync"
	"time"

	session "github.com/go-session/session/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	jsoniter "github.com/json-iterator/go"
)

var (
	_ session.ManagerStore = &managerStore{}
	_ session.Store        = &store{}
)

// NewStore Create an instance of a mongo store
func NewStore(url, dbName, cName string) session.ManagerStore {
	return newManagerStore(url, dbName, cName)
}

func newManagerStore(url, dbName, cName string) *managerStore {
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))

	return &managerStore{
		client: client,
		dbName: dbName,
		cName:  cName,
	}
}

type managerStore struct {
	client *mongo.Client
	dbName string
	cName  string
}

func (s *managerStore) getValue(ctx context.Context, sid string) (string, error) {
	item := new(sessionItem)
	coll := s.client.Database(s.dbName).Collection(s.cName)
	filter := bson.D{{Key: "_id", Value: sid}}
	err := coll.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil
		}

		return "", err
	} else if item.ExpiredAt.Before(time.Now()) {
		return "", nil
	}
	return item.Value, nil
}

func (s *managerStore) parseValue(value string) (map[string]interface{}, error) {
	var values map[string]interface{}
	if len(value) > 0 {
		err := jsoniter.Unmarshal([]byte(value), &values)
		if err != nil {
			return nil, err
		}
	}

	return values, nil
}

func (s *managerStore) Check(ctx context.Context, sid string) (bool, error) {
	val, err := s.getValue(ctx, sid)
	if err != nil {
		return false, err
	}
	return val != "", nil
}

func (s *managerStore) Create(ctx context.Context, sid string, expired int64) (session.Store, error) {
	return newStore(ctx, s, sid, expired, nil), nil
}

func (s *managerStore) Update(ctx context.Context, sid string, expired int64) (session.Store, error) {
	value, err := s.getValue(ctx, sid)
	if err != nil {
		return nil, err
	} else if value == "" {
		return newStore(ctx, s, sid, expired, nil), nil
	}

	coll := s.client.Database(s.dbName).Collection(s.cName)
	_, err = coll.UpdateByID(ctx, sid, bson.M{
		"$set": bson.M{
			"expired_at": time.Now().Add(time.Duration(expired) * time.Second),
		},
	})

	if err != nil {
		return nil, err
	}

	values, err := s.parseValue(value)
	if err != nil {
		return nil, err
	}

	return newStore(ctx, s, sid, expired, values), nil
}

func (s *managerStore) Delete(ctx context.Context, sid string) error {
	coll := s.client.Database(s.dbName).Collection(s.cName)
	filter := bson.D{{Key: "_id", Value: sid}}
	_, err := coll.DeleteOne(ctx, filter)
	return err
}

func (s *managerStore) Refresh(ctx context.Context, oldsid, sid string, expired int64) (session.Store, error) {
	value, err := s.getValue(ctx, oldsid)
	if err != nil {
		return nil, err
	} else if value == "" {
		return newStore(ctx, s, sid, expired, nil), nil
	}

	coll := s.client.Database(s.dbName).Collection(s.cName)
	filter := bson.D{{Key: "_id", Value: sid}}
	upsert := true
	op := &options.ReplaceOptions{
		Upsert: &upsert,
	}
	_, err = coll.ReplaceOne(ctx, filter, &sessionItem{
		ID:        sid,
		Value:     value,
		ExpiredAt: time.Now().Add(time.Duration(expired) * time.Second),
	}, op)
	if err != nil {
		return nil, err
	}
	_, err = coll.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	values, err := s.parseValue(value)
	if err != nil {
		return nil, err
	}

	return newStore(ctx, s, sid, expired, values), nil
}

func (s *managerStore) Close() error {
	s.client.Disconnect(context.Background())
	return nil
}

func newStore(ctx context.Context, s *managerStore, sid string, expired int64, values map[string]interface{}) *store {
	if values == nil {
		values = make(map[string]interface{})
	}

	return &store{
		client:  s.client,
		dbName:  s.dbName,
		cName:   s.cName,
		ctx:     ctx,
		sid:     sid,
		expired: expired,
		values:  values,
	}
}

type store struct {
	sync.RWMutex
	ctx     context.Context
	client  *mongo.Client
	dbName  string
	cName   string
	sid     string
	expired int64
	values  map[string]interface{}
}

func (s *store) Context() context.Context {
	return s.ctx
}

func (s *store) SessionID() string {
	return s.sid
}

func (s *store) Set(key string, value interface{}) {
	s.Lock()
	s.values[key] = value
	s.Unlock()
}

func (s *store) Get(key string) (interface{}, bool) {
	s.RLock()
	val, ok := s.values[key]
	s.RUnlock()
	return val, ok
}

func (s *store) Delete(key string) interface{} {
	s.RLock()
	v, ok := s.values[key]
	s.RUnlock()
	if ok {
		s.Lock()
		delete(s.values, key)
		s.Unlock()
	}
	return v
}

func (s *store) Flush() error {
	s.Lock()
	s.values = make(map[string]interface{})
	s.Unlock()
	return s.Save()
}

func (s *store) Save() error {
	var value string

	s.RLock()
	if len(s.values) > 0 {
		buf, err := jsoniter.Marshal(s.values)
		if err != nil {
			s.RUnlock()
			return err
		}
		value = string(buf)
	}
	s.RUnlock()

	coll := s.client.Database(s.dbName).Collection(s.cName)
	filter := bson.D{{Key: "_id", Value: s.sid}}
	upsert := true
	op := &options.ReplaceOptions{
		Upsert: &upsert,
	}
	_, err := coll.ReplaceOne(s.Context(), filter, &sessionItem{
		ID:        s.sid,
		Value:     value,
		ExpiredAt: time.Now().Add(time.Duration(s.expired) * time.Second),
	}, op)
	return err
}

// Data items stored in mongo
type sessionItem struct {
	ID        string    `bson:"_id"`
	Value     string    `bson:"value"`
	ExpiredAt time.Time `bson:"expired_at"`
}
