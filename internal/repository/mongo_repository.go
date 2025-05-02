package repository

import (
	"context"
	"time"

	"notificationservice/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoRepository struct {
    client     *mongo.Client
    database   string
    collection string
}

func NewMongoRepository(uri, database string) (*MongoRepository, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
    if err != nil {
        return nil, err
    }

    return &MongoRepository{
        client:     client,
        database:   database,
        collection: "notifications",
    }, nil
}

func (r *MongoRepository) SaveNotification(notification *models.Notification) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    collection := r.client.Database(r.database).Collection(r.collection)

    notification.ID = primitive.NewObjectID()
    notification.CreatedAt = time.Now()
    notification.Status = models.DeliveryStatus{
        Status:    models.Pending,
        UpdatedAt: time.Now(),
    }

    _, err := collection.InsertOne(ctx, notification)
    return err
}

func (r *MongoRepository) GetUnreadNotifications(userId string) ([]models.Notification, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    collection := r.client.Database(r.database).Collection(r.collection)
    
    filter := bson.M{
        "userId": userId,
        "isReceived": false,
    }

    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }

    var notifications []models.Notification
    if err = cursor.All(ctx, &notifications); err != nil {
        return nil, err
    }

    return notifications, nil
}

func (r *MongoRepository) Ping(ctx context.Context) error {
    return r.client.Ping(ctx, readpref.Primary())
}

func (r *MongoRepository) Disconnect(ctx context.Context) error {
    return r.client.Disconnect(ctx)
}