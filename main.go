package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"labix.org/v2/mgo/bson"
)

const (
	URL        string = "127.0.0.1"
	USERNAME   string = "root"
	PASSWORD   string = "liurui97128224"
	PORT       uint   = 27017
	DATABASE   string = "test"
	COLLECTION string = "cron_log"
)

type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

type LogRecord struct {
	JobName  string    `bson:"jobName"`
	Command  string    `bson:"command"`
	Err      string    `bson:"err"`
	Content  string    `bson:"content"`
	ExecTime TimePoint `bson:"execTime"`
}

func main() {
	record := &LogRecord{
		JobName: "job11",
		Command: "echo hello",
		Err:     "",
		Content: "hello",
		ExecTime: TimePoint{
			StartTime: time.Now().Unix() - 10,
			EndTime:   time.Now().Unix(),
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second) //生成连接mongodb的上下文
	defer cancel()
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin", USERNAME, PASSWORD, URL, PORT, DATABASE)) //使用authSource才能在连接后操作mongodb里面的数据
	client, err := mongo.Connect(ctx, clientOptions)                                                                                              //连接到mongodb
	if err != nil {
		panic(err)
	}
	database := client.Database(DATABASE)            //选择mongodb的数据库
	collection := database.Collection(COLLECTION)    //选择mongodb数据库里面的集合
	result, err := collection.InsertOne(ctx, record) //向集合中插入数据
	if err != nil {
		panic(err)
	}
	docId := result.InsertedID.(primitive.ObjectID)
	fmt.Println("自增ID:", docId.Hex())
	cur1, err := collection.Find(ctx, bson.M{}) //Find方法是读取mongodb里面的所有数据，bson.M{}表示mongodb里面的数据是Map，而bson.D{}表示mongodb里面的数据是Data
	if err != nil {
		panic(err)
	}
	defer cur1.Close(ctx)
	// 下面的for循环就是逐个读取cur的内容，并解析成struct
	// for cur.Next(ctx) {
	// 	res := LogRecord{}
	// 	err = cur.Decode(&res)
	// 	fmt.Println(res)
	// }
	// 下面的All方法更加好，直接把cur解析成struct切片
	var records []LogRecord
	cur1.All(ctx, &records)
	fmt.Println(records)
	var rec LogRecord
	filter := bson.M{"jobName": "job11"}
	err = collection.FindOne(ctx, filter).Decode(&rec) //通过过滤器来找到1个结果
	if err != nil {
		panic(err)
	}
	fmt.Println(rec)
	filter = bson.M{"execTime": bson.M{"startTime": 1642824877, "endTime": 1642824887}} //这个示例展示了map里嵌套map的bson表示方法
	err = collection.FindOne(ctx, filter).Decode(&rec)
	if err != nil {
		panic(err)
	}
	fmt.Println(rec)
}
