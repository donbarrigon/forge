package db

import (
	"context"
	"errors"
	"reflect"

	"github.com/donbarrigon/forge/errs"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoModel interface {
	GetID() bson.ObjectID
	Coll() string
}

type OdmModel interface {
	CollectionName() string
	GetID() bson.ObjectID
	SetID(id bson.ObjectID)
	BeforeCreate() *errs.Error
	BeforeUpdate() *errs.Error
	BeforeDelete() *errs.Error
	AfterCreate() *errs.Error
	AfterUpdate() *errs.Error
	AfterDelete() *errs.Error

	// funciones nesesarias para el move to trash
	Create() *errs.Error
	Delete() *errs.Error

	// funciones nesesarias para el history record
	GetOriginal() map[string]any
	GetDirty() map[string]any
	SetOriginal(original map[string]any)
	SetDirty(dirty map[string]any)
}

type Collection []OdmModel

type Odm struct {
	Model    OdmModel       `bson:"-" json:"-"`
	dirty    map[string]any `bson:"-" json:"-"`
	original map[string]any `bson:"-" json:"-"`
}

func (o *Odm) BeforeCreate() *errs.Error { return nil }
func (o *Odm) BeforeUpdate() *errs.Error { return nil }
func (o *Odm) BeforeDelete() *errs.Error { return nil }
func (o *Odm) AfterCreate() *errs.Error  { return nil }
func (o *Odm) AfterUpdate() *errs.Error  { return nil }
func (o *Odm) AfterDelete() *errs.Error  { return nil }

func (o *Odm) GetOriginal() map[string]any         { return o.original }
func (o *Odm) GetDirty() map[string]any            { return o.dirty }
func (o *Odm) SetOriginal(original map[string]any) { o.original = original }
func (o *Odm) SetDirty(dirty map[string]any)       { o.dirty = dirty }

func (o *Odm) FindByHexID(id string) *errs.Error {

	objectId, e := bson.ObjectIDFromHex(id)
	if e != nil {
		return errs.HexID(e)
	}
	filter := bson.D{bson.E{Key: "_id", Value: objectId}}
	if e := Mongo.Collection(o.Model.CollectionName()).FindOne(context.TODO(), filter).Decode(o.Model); e != nil {
		return errs.Mongo(e)
	}
	return nil
}

func (o *Odm) FindByID(id bson.ObjectID) *errs.Error {
	filter := bson.D{bson.E{Key: "_id", Value: id}}
	if e := Mongo.Collection(o.Model.CollectionName()).FindOne(context.TODO(), filter).Decode(o.Model); e != nil {
		return errs.Mongo(e)
	}
	return nil
}

func (o *Odm) First(field string, value any) *errs.Error {
	filter := bson.D{bson.E{Key: field, Value: value}}
	if e := Mongo.Collection(o.Model.CollectionName()).FindOne(context.TODO(), filter).Decode(o.Model); e != nil {
		return errs.Mongo(e)
	}
	return nil
}

func (o *Odm) FindOne(filter bson.D, opts ...options.Lister[options.FindOneOptions]) *errs.Error {
	if e := Mongo.Collection(o.Model.CollectionName()).FindOne(context.TODO(), filter, opts...).Decode(o.Model); e != nil {
		return errs.Mongo(e)
	}
	return nil
}

func (o *Odm) Find(result any, filter bson.D, opts ...options.Lister[options.FindOptions]) *errs.Error {
	ctx := context.TODO()
	cursor, e := Mongo.Collection(o.Model.CollectionName()).Find(ctx, filter, opts...)
	if e != nil {
		return errs.Mongo(e)
	}
	if e = cursor.All(ctx, result); e != nil {
		return errs.Mongo(e)
	}
	return nil
}

// busqueda eq
func (o *Odm) FindByField(result any, field string, value any, opts ...options.Lister[options.FindOptions]) *errs.Error {
	filter := bson.D{bson.E{Key: field, Value: value}}
	ctx := context.TODO()
	cursor, e := Mongo.Collection(o.Model.CollectionName()).Find(ctx, filter, opts...)
	if e != nil {
		return errs.Mongo(e)
	}
	if e = cursor.All(ctx, result); e != nil {
		return errs.Mongo(e)
	}
	return nil
}

func (o *Odm) Aggregate(result any, pipeline mongo.Pipeline) *errs.Error {
	ctx := context.TODO()
	cursor, e := Mongo.Collection(o.Model.CollectionName()).Aggregate(ctx, pipeline)
	if e != nil {
		return errs.Mongo(e)
	}
	if e = cursor.All(ctx, result); e != nil {
		return errs.Mongo(e)
	}
	return nil
}

func (o *Odm) AggregateOne(pipeline mongo.Pipeline) *errs.Error {
	ctx := context.TODO()
	cursor, e := Mongo.Collection(o.Model.CollectionName()).Aggregate(ctx, pipeline)
	if e != nil {
		return errs.Mongo(e)
	}
	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		if e := cursor.Decode(o.Model); e != nil {
			return errs.Mongo(e)
		}
	} else {
		return errs.NotFoundMsg("The document does not exist", errors.New("!cursor.Next(ctx)"))
	}
	return nil
}

func (o *Odm) Create() *errs.Error {
	if e := o.Model.BeforeCreate(); e != nil {
		return e
	}
	result, e := Mongo.Collection(o.Model.CollectionName()).InsertOne(context.TODO(), o.Model)
	if e != nil {
		return errs.Mongo(e)
	}
	o.Model.SetID(result.InsertedID.(bson.ObjectID))
	return o.Model.AfterCreate()
}

func (o *Odm) CreateBy(validator any) *errs.Error {
	if e := Fill(o.Model, validator); e != nil {
		return e
	}
	return o.Create()
}

// usela solo si tienes pereza.
func (o *Odm) CreateMany(data any) *errs.Error {

	v := reflect.ValueOf(data)

	if v.Kind() != reflect.Slice {
		return errs.InternalMsg("The data collection cannot be saved", errors.New("CreateMany only accepts slices of Model"))
	}
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i).Interface()
		if e := elem.(OdmModel).BeforeCreate(); e != nil {
			return e
		}
	}
	collection := Mongo.Collection(o.Model.CollectionName())
	result, e := collection.InsertMany(context.TODO(), data)
	if e != nil {
		return errs.Mongo(e)
	}
	he := []error{}
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i).Interface()
		elem.(OdmModel).SetID(result.InsertedIDs[i].(bson.ObjectID))
		if e := elem.(OdmModel).AfterCreate(); e != nil {
			he = append(he, e)
		}
	}
	if len(he) > 0 {
		err := errs.InternalMsg("Something went wrong while saving the data collection", nil)
		err.Data = he
		return err
	}
	return nil
}

func (o *Odm) Update() *errs.Error {
	if e := o.Model.BeforeUpdate(); e != nil {
		return e
	}
	filter := bson.D{bson.E{Key: "_id", Value: o.Model.GetID()}}
	update := bson.D{bson.E{Key: "$set", Value: o.Model}}

	result, e := Mongo.Collection(o.Model.CollectionName()).UpdateOne(context.TODO(), filter, update)
	if e != nil {
		return errs.Mongo(e)
	}
	if result.MatchedCount == 0 {
		return errs.NotFoundMsg("The document to update does not exist", errors.New("!result.MatchedCount == 0"))
	}

	if result.ModifiedCount == 0 {
		return errs.ConflictMsg("No changes were applied when saving the document", errors.New("!result.ModifiedCount == 0"))
	}
	return o.Model.AfterUpdate()
}

func (o *Odm) UpdateBy(validator any) *errs.Error {

	if e := Filld(o.Model, validator); e != nil {
		return e
	}
	return o.Update()
}

// OjO no usa el hook BeforeUpdate

func (o *Odm) Delete() *errs.Error {
	if e := o.Model.BeforeDelete(); e != nil {
		return e
	}

	filter := bson.D{bson.E{Key: "_id", Value: o.Model.GetID()}}

	result, e := Mongo.Collection(o.Model.CollectionName()).DeleteOne(context.TODO(), filter)
	if e != nil {
		return errs.Mongo(e)
	}
	if result.DeletedCount == 0 {
		return errs.ConflictMsg("The document was not deleted", errors.New("!result.DeletedCount == 0"))
	}
	return o.Model.AfterDelete()
}
