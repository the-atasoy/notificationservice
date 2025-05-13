package repository

import (
	"context"
	"time"

	"notificationservice/internal/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (repository *MongoRepository) SaveNotification(notification *models.Notification) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    collection := repository.client.Database(repository.database).Collection(repository.collection)

    notification.ID = primitive.NewObjectID()
    notification.CreatedAt = time.Now()
    notification.DeliveryStatus = models.DeliveryStatus{
        NotificationStatus:    models.Pending,
        UpdatedAt: time.Now(),
    }

    _, err := collection.InsertOne(ctx, notification)
    return err
}

func (repository *MongoRepository) GetUnreadNotifications(userId string) ([]models.Notification, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    collection := repository.client.Database(repository.database).Collection(repository.collection)
    
    filter := bson.M{
        "userId": userId,
        "receivedAt": bson.M{"$exists": false},
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

func (repository *MongoRepository) GetUnsentNotifications(externalId uuid.UUID) (*models.Notification, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    collection := repository.client.Database(repository.database).Collection(repository.collection)

    filter := bson.M{
        "externalId": externalId,
        "deliveryStatus.notificationStatus": bson.M{
            "$in": bson.A{models.Failed, models.Pending},
            },
    }
    var notification models.Notification
    err := collection.FindOne(ctx, filter).Decode(&notification)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, nil
        }
        return nil, err
    }

    return &notification, nil
}

func (repository *MongoRepository) UpdateNotificationStatus(notificationID primitive.ObjectID, status models.DeliveryStatus) error {
    status.UpdatedAt = time.Now()

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    collection := repository.client.Database(repository.database).Collection(repository.collection)

    filter := bson.M{"_id": notificationID}
    update := bson.M{
        "$set": bson.M{
            "deliveryStatus": status,
        },
    }

    _, err := collection.UpdateOne(ctx, filter, update)
    return err
}