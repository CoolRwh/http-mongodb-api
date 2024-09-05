package controller

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"http-mongodb-api/db"
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
		c.JSON(http.StatusOK, Message{
			Code: 422,
			Msg:  "ID格式错误！" + err.Error(),
			Data: nil,
		})
		return
	}

	result, err := collection.UpdateOne(context.TODO(), bson.D{{"_id", objectID}}, bson.M{"$set": params.Update})
	if err != nil {
		c.JSON(http.StatusOK, Message{
			Code: 423,
			Msg:  "修改失败！" + err.Error(),
			Data: nil,
		})
	}
	c.JSON(http.StatusOK, Message{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
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
		c.JSON(http.StatusOK, Message{
			Code: 422,
			Msg:  "参数校验失败！" + err.Error(),
			Data: nil,
		})
		return
	}
	filter, err := CheckFilterData(params.Filter)
	if err != nil {
		c.JSON(http.StatusOK, Message{
			Code: 422,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}
	collection := db.MongoDbClient.Database(params.Database).Collection(params.Collection)
	var items []bson.M
	find, err := collection.Find(context.TODO(), filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusOK, Message{
				Code: http.StatusNotFound,
				Msg:  "未找到匹配的文档",
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, Message{
			Code: 500,
			Msg:  "查询失败" + err.Error(),
			Data: nil,
		})
		return
	}
	_ = find.All(context.TODO(), &items)
	c.JSON(http.StatusOK, Message{
		Code: 0,
		Msg:  "ok",
		Data: struct {
			Items []bson.M `json:"items"`
			Total int      `json:"total"`
		}{
			Items: items,
			Total: len(items),
		},
	})
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
		c.JSON(http.StatusOK, Message{
			Code: 422,
			Msg:  "参数校验失败！" + err.Error(),
			Data: nil,
		})
		return
	}
	filter, err := CheckFilterData(params.Filter)
	if err != nil {
		c.JSON(http.StatusOK, Message{
			Code: 422,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}
	collection := db.MongoDbClient.Database(params.Database).Collection(params.Collection)
	var result bson.M
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusOK, Message{
				Code: http.StatusNotFound,
				Msg:  "未找到匹配的文档",
				Data: nil,
			})
			return
		}
		c.JSON(http.StatusOK, Message{
			Code: 500,
			Msg:  "查询失败" + err.Error(),
			Data: nil,
		})
		return
	}
	c.JSON(http.StatusOK, Message{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
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
		c.JSON(http.StatusOK, Message{
			Code: 422,
			Msg:  "参数校验失败！" + err.Error(),
			Data: nil,
		})
		return
	}
	documents := make([]interface{}, len(params.Data))
	for i, doc := range params.Data {
		if bsonDoc, ok := doc.(map[string]interface{}); ok {
			documents[i] = bson.M(bsonDoc)
		} else {
			c.JSON(http.StatusOK, Message{
				Code: 422,
				Msg:  "数据格式错误！",
				Data: nil,
			})
			return
		}
	}
	// 选择数据库和集合
	collection := db.MongoDbClient.Database(params.Database).Collection(params.Collection)
	manyResult, err := collection.InsertMany(context.TODO(), documents)
	if err != nil {
		c.JSON(http.StatusOK, Message{
			Code: 423,
			Msg:  "添加失败！" + err.Error(),
			Data: nil,
		})
		return
	}
	c.JSON(http.StatusOK, Message{
		Code: 0,
		Msg:  "ok",
		Data: manyResult,
	})
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
		c.JSON(http.StatusOK, Message{
			Code: 422,
			Msg:  "参数校验失败！" + err.Error(),
			Data: nil,
		})
		return
	}
	var doc bson.M
	switch v := params.Data.(type) {
	case map[string]interface{}:
		doc = bson.M(v)
	default:
		c.JSON(http.StatusOK, Message{
			Code: 422,
			Msg:  "数据格式错误！",
			Data: nil,
		})
		return
	}
	collection := db.MongoDbClient.Database(params.Database).Collection(params.Collection)
	result, err := collection.InsertOne(context.TODO(), doc)
	if err != nil {
		c.JSON(http.StatusOK, Message{
			Code: 423,
			Msg:  "添加失败！" + err.Error(),
			Data: nil,
		})
	}
	c.JSON(http.StatusOK, Message{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
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
		c.JSON(http.StatusOK, Message{
			Code: 422,
			Msg:  "ID格式错误！" + err.Error(),
			Data: nil,
		})
		return
	}

	result, err := collection.UpdateOne(context.TODO(), bson.D{{"_id", objectID}}, bson.M{"$set": params.Update})
	if err != nil {
		c.JSON(http.StatusOK, Message{
			Code: 423,
			Msg:  "修改失败！" + err.Error(),
			Data: nil,
		})
	}
	c.JSON(http.StatusOK, Message{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
	return
}

func CheckFilterData(data any) (bson.M, error) {
	filterBSON, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("过滤条件格式错误")
	}
	filter := bson.M{}
	for key, value := range filterBSON {
		filter[key] = value
	}
	return filter, nil
}

// Helper function to convert interface{} to bson.M
func toBSON(data interface{}) (bson.M, error) {
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
