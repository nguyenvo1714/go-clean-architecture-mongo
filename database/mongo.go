package database

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"reflect"
	"time"
)

type Database interface {
	Collection(string) Collection
	Client() Client
}

type Collection interface {
	FindOne(context.Context, interface{}) SingleResult
	InsertOne(context.Context, interface{}) (interface{}, error)
	InsertMany(context.Context, []interface{}) ([]interface{}, error)
	DeleteOne(context.Context, interface{}) (int64, error)
	Find(context.Context, interface{}, ...*options.FindOptions) (Cursor, error)
	CountDocuments(context.Context, interface{}, ...*options.CountOptions) (int64, error)
	Aggregate(context.Context, interface{}) (Cursor, error)
	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateMany(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

type SingleResult interface {
	Decode(interface{}) error
}

type Cursor interface {
	Close(context.Context) error
	Decode(interface{}) error
	Next(context.Context) bool
	All(context.Context, interface{}) error
}

type Client interface {
	Database(string) Database
	Disconnect(context.Context) error
	StartSession() (mongo.Session, error)
	UseSession(ctx context.Context, fn func(sessionContext mongo.SessionContext) error) error
	Ping(context.Context) error
}

type mongoClient struct {
	cl *mongo.Client
}

type mongoDatabase struct {
	db *mongo.Database
}

type mongoCollection struct {
	coll *mongo.Collection
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

type mongoCursor struct {
	mc *mongo.Cursor
}

type mongoSession struct {
	mongo.Session
}

type nullAwareDecoder struct {
	defDecoder bsoncodec.ValueDecoder
	zeroValue  reflect.Value
}

func (m mongoDatabase) Collection(s string) Collection {
	collection := m.db.Collection(s)

	return &mongoCollection{coll: collection}
}

func (m mongoDatabase) Client() Client {
	client := m.db.Client()

	return &mongoClient{cl: client}
}

func (d *nullAwareDecoder) DecodeValue(dctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if vr.Type() != bson.TypeNull {
		return d.defDecoder.DecodeValue(dctx, vr, val)
	}

	if !val.CanSet() {
		return errors.New("value not settable")
	}
	if err := vr.ReadNull(); err != nil {
		return err
	}
	// Set the zero value of val's type:
	val.Set(d.zeroValue)
	return nil
}

func NewClient(ctx context.Context, connection string) (Client, error) {
	time.Local = time.UTC
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection))

	return &mongoClient{
		cl: client,
	}, err
}

func (m mongoClient) Database(s string) Database {
	db := m.cl.Database(s)

	return &mongoDatabase{db: db}
}

func (m mongoClient) Disconnect(ctx context.Context) error {
	return m.cl.Disconnect(ctx)
}

func (m mongoClient) StartSession() (mongo.Session, error) {
	session, err := m.cl.StartSession()

	return mongoSession{session}, err
}

func (m mongoClient) UseSession(ctx context.Context, fn func(sessionContext mongo.SessionContext) error) error {
	return m.cl.UseSession(ctx, fn)
}

func (m mongoClient) Ping(ctx context.Context) error {
	return m.cl.Ping(ctx, readpref.Primary())
}

func (m mongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResult {
	singleResult := m.coll.FindOne(ctx, filter)

	return &mongoSingleResult{sr: singleResult}
}

func (m mongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	result, err := m.coll.InsertOne(ctx, document)

	return result.InsertedID, err
}

func (m mongoCollection) InsertMany(ctx context.Context, documents []interface{}) ([]interface{}, error) {
	result, err := m.coll.InsertMany(ctx, documents)

	return result.InsertedIDs, err
}

func (m mongoCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	res, err := m.coll.DeleteOne(ctx, filter)

	return res.DeletedCount, err
}

func (m mongoCollection) Find(ctx context.Context, filter interface{}, findOptions ...*options.FindOptions) (Cursor, error) {
	cursor, err := m.coll.Find(ctx, filter, findOptions...)

	return &mongoCursor{mc: cursor}, err
}

func (m mongoCollection) CountDocuments(ctx context.Context, filter interface{}, countOptions ...*options.CountOptions) (int64, error) {
	return m.coll.CountDocuments(ctx, filter, countOptions...)
}

func (m mongoCollection) Aggregate(ctx context.Context, pipeline interface{}) (Cursor, error) {
	cursor, err := m.coll.Aggregate(ctx, pipeline)

	return &mongoCursor{mc: cursor}, err
}

func (m mongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, updateOptions ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return m.coll.UpdateOne(ctx, filter, update, updateOptions...)
}

func (m mongoCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, updateOptions ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return m.coll.UpdateMany(ctx, filter, update, updateOptions...)
}

func (m mongoSingleResult) Decode(val interface{}) error {
	return m.sr.Decode(val)
}

func (m *mongoCursor) Close(ctx context.Context) error {
	return m.mc.Close(ctx)
}

func (m *mongoCursor) Decode(val interface{}) error {
	return m.mc.Decode(val)
}

func (m *mongoCursor) Next(ctx context.Context) bool {
	return m.mc.Next(ctx)
}

func (m *mongoCursor) All(ctx context.Context, result interface{}) error {
	return m.mc.All(ctx, result)
}
