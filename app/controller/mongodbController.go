package controller

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"http-mongodb-api/app/models/response"
	"http-mongodb-api/pkg/db"
	"net/http"
)

var Mongodb = new(MongodbController)

type MongodbController struct{}

// Message 消息
type Message struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
	Data any    `json:"data"`
}

func (t *MongodbController) Cs(c *gin.Context) {
	var params struct {
		Database   string `json:"database"`
		Collection string `json:"collection"`
		ID         string `json:"id"`
		Update     bson.M `json:"update"`
		Options    any    `json:"options"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&params)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.RequestParamError, "参数校验失败！"+err.Error()))

		return
	}

	// 选择数据库和集合
	collection := db.MongoDbClient.Database(params.Database).Collection(params.Collection)
	// Convert ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.RequestParamError, "ID格式错误！"+err.Error()))
		return
	}
	result, err := collection.UpdateOne(context.TODO(), bson.D{{"_id", objectID}}, bson.M{"$set": params.Update})
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.DataUpdateError, "修改失败！"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Ok(result))
	return
}

func (t *MongodbController) Find(c *gin.Context) {

	var params struct {
		Database   string      `json:"database"`
		Collection string      `json:"collection"`
		Filter     interface{} `json:"filter"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&params)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.RequestParamError, err.Error()))
		return
	}
	filter, err := CheckFilterData(params.Filter)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.RequestParamError, err.Error()))
		return
	}
	collection := db.MongoDbClient.Database(params.Database).Collection(params.Collection)
	var items []bson.M
	find, err := collection.Find(context.TODO(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusOK, response.Result(response.DataNotExist, err.Error()))
			return
		}
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	_ = find.All(context.TODO(), &items)
	var responseData struct {
		Items []bson.M `json:"items"`
		Total int      `json:"total"`
	}
	responseData.Items = items
	responseData.Total = len(items)
	c.JSON(http.StatusOK, response.Ok(responseData))
	return
}

func (t *MongodbController) FindOne(c *gin.Context) {
	var params struct {
		Database   string      `json:"database"`
		Collection string      `json:"collection"`
		Filter     interface{} `json:"filter"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&params)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.RequestParamError, "参数校验失败！"+err.Error()))
		return
	}
	filter, err := CheckFilterData(params.Filter)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.RequestParamError, err.Error()))
		return
	}
	collection := db.MongoDbClient.Database(params.Database).Collection(params.Collection)
	var result bson.M
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusOK, response.Result(response.DataNotExist, nil))
			return
		}
		c.JSON(http.StatusOK, response.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Ok(result))
	return
}

func (t *MongodbController) InsertMany(c *gin.Context) {
	var params struct {
		Database   string        `json:"database"`
		Collection string        `json:"collection"`
		Data       []interface{} `json:"data"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&params)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.RequestParamError, err.Error()))
		return
	}
	documents := make([]interface{}, len(params.Data))
	for i, doc := range params.Data {
		if bsonDoc, ok := doc.(map[string]interface{}); ok {
			documents[i] = bson.M(bsonDoc)
		} else {
			c.JSON(http.StatusOK, response.Result(response.RequestParamError, err.Error()))
			return
		}
	}
	// 选择数据库和集合
	collection := db.MongoDbClient.Database(params.Database).Collection(params.Collection)
	manyResult, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.AddDataError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Ok(manyResult))
	return
}

func (t *MongodbController) InsertOne(c *gin.Context) {
	var params struct {
		Database   string      `json:"database"`
		Collection string      `json:"collection"`
		Data       interface{} `json:"data"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&params)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.RequestParamError, err.Error()))
		return
	}
	var doc bson.M
	switch v := params.Data.(type) {
	case map[string]interface{}:
		doc = bson.M(v)
	default:
		c.JSON(http.StatusOK, response.Result(response.RequestParamError, err.Error()))
		return
	}
	collection := db.MongoDbClient.Database(params.Database).Collection(params.Collection)
	result, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.AddDataError, err.Error()))
	}
	c.JSON(http.StatusOK, response.Ok(result))
	return
}

// UpdateById 根据ID修改资料
func (t *MongodbController) UpdateById(c *gin.Context) {
	var params struct {
		Database   string      `json:"database"`
		Collection string      `json:"collection"`
		ID         string      `json:"id"`
		Update     bson.M      `json:"update"`
		Options    interface{} `json:"options"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&params)
	if err != nil {
		c.JSON(http.StatusOK, Message{
			Code: 422,
			Msg:  "参数校验失败！" + err.Error(),
			Data: nil,
		})
		return
	}
	// 选择数据库和集合
	collection := db.MongoDbClient.Database(params.Database).Collection(params.Collection)
	// Convert ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.RequestParamError, "ID格式错误 "+err.Error()))
		return
	}

	result, err := collection.UpdateOne(context.TODO(), bson.D{{"_id", objectID}}, bson.M{"$set": params.Update})
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.DataUpdateError, err.Error()))
	}
	c.JSON(http.StatusOK, response.Ok(result))
	return
}

// UpdateOne 更新一条数据
func (t *MongodbController) UpdateOne(c *gin.Context) {
	var params struct {
		Database   string      `json:"database"`
		Collection string      `json:"collection"`
		Filter     interface{} `json:"filter"`
		Update     bson.M      `json:"update"`
		Options    interface{} `json:"options"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&params)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.RequestParamError, err.Error()))
		return
	}
	filter, err := CheckFilterData(params.Filter)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.RequestParamError, err.Error()))
		return
	}
	// 选择数据库和集合
	collection := db.MongoDbClient.Database(params.Database).Collection(params.Collection)

	result, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$set": params.Update})
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.DataUpdateError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Ok(result))
	return
}

func (t *MongodbController) UpdateMany(c *gin.Context) {
	var params struct {
		Database   string      `json:"database"`
		Collection string      `json:"collection"`
		Filter     interface{} `json:"filter"`
		Update     bson.M      `json:"update"`
		Options    interface{} `json:"options"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&params)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.RequestParamError, err.Error()))
		return
	}
	filter, err := CheckFilterData(params.Filter)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.RequestParamError, err.Error()))
		return
	}
	// 选择数据库和集合
	collection := db.MongoDbClient.Database(params.Database).Collection(params.Collection)
	result, err := collection.UpdateMany(context.TODO(), filter, bson.M{"$set": params.Update})
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.DataUpdateError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Ok(result))
	return
}

func (t *MongodbController) DeleteById(c *gin.Context) {
	var params struct {
		Database   string      `json:"database"`
		Collection string      `json:"collection"`
		ID         string      `json:"id"`
		Options    interface{} `json:"options"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&params)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.RequestParamError, err.Error()))
		return
	}
	// Convert ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.RequestParamError, "ID 数据格式异常"+err.Error()))
		return
	}
	// 选择数据库和集合
	collection := db.MongoDbClient.Database(params.Database).Collection(params.Collection)

	result, err := collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.DataUpdateError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Ok(result))
	return
}

func (t *MongodbController) DeleteMany(c *gin.Context) {
	var params struct {
		Database   string      `json:"database"`
		Collection string      `json:"collection"`
		Filter     interface{} `json:"filter"`
		Options    interface{} `json:"options"`
	}
	err := json.NewDecoder(c.Request.Body).Decode(&params)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.RequestParamError, err.Error()))
		return
	}
	filter, err := CheckFilterData(params.Filter)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.RequestParamError, err.Error()))
		return
	}
	// 选择数据库和集合
	collection := db.MongoDbClient.Database(params.Database).Collection(params.Collection)
	result, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusOK, response.Result(response.DataUpdateError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Ok(result))
	return
}

func CheckFilterData(data any) (bson.M, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		filter := bson.M{}
		for key, value := range v {
			filter[key] = value
		}
		return filter, nil
	default:
		return nil, errors.New("不支持的过滤条件格式")
	}
}

// Helper function to convert interface{} to bson.M
func dataToBSON(data interface{}) (bson.M, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		filter := bson.M{}
		for key, value := range v {
			filter[key] = value
		}
		return filter, nil
	default:
		return nil, errors.New("不支持的过滤条件格式")
	}
}
