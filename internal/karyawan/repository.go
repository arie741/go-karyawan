package karyawan

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(collection *mongo.Collection) *Repository {
	return &Repository{collection: collection}
}

func (r *Repository) GetAll(offset int, limit int, filter bson.M) ([]Karyawan, error) {
	sort := bson.D{{"Name", 1}}
	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(limit)).SetSort(sort)
	cursor, err := r.collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}

	var karyawans []Karyawan
	if err := cursor.All(context.TODO(), &karyawans); err != nil {
		return nil, err
	}

	return karyawans, nil
}

func (r *Repository) CountAll(filter bson.M) (int64, error) {
	totalDocuments, err := r.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return totalDocuments, nil
}

func (r *Repository) FindById(id string) (Karyawan, error) {
	filter := bson.D{{"Id", id}}

	var karyawan Karyawan
	if err := r.collection.FindOne(context.TODO(), filter).Decode(&karyawan); err != nil {
		return karyawan, err
	}

	return karyawan, nil
}

func (r *Repository) Add(karyawan *Karyawan) (Karyawan, error) {
	_, err := r.collection.InsertOne(context.TODO(), karyawan)
	if err != nil {
		return *karyawan, err
	}

	return *karyawan, nil
}

func (r *Repository) Update(karyawan *Karyawan) (Karyawan, error) {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	filter := bson.M{"Id": karyawan.Id}

	update := bson.M{"$set": bson.M{
		"Name":      karyawan.Name,
		"BirthDate": karyawan.BirthDate,
		"Salary":    karyawan.Salary,
		"Position":  karyawan.Position,
		"Joined":    karyawan.Joined,
	}}

	var updatedKaryawan Karyawan
	err := r.collection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&updatedKaryawan)
	if err != nil {
		return *karyawan, err
	}

	return updatedKaryawan, nil
}

func (r *Repository) Delete(id string) error {
	filter := bson.M{"Id": id}
	_, err := r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}
