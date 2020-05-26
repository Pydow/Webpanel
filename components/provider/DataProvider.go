package provider

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

var Cache redis.Conn
var Coll *mongo.Collection
var mongoURI = ""

type User struct {
	_id string
	Username string
}

func GetUserByCookie(w http.ResponseWriter, r *http.Request) string {

	c, err := r.Cookie("SESSION_ID")
	if err != nil {
		return "none"
	}
	sessionToken := c.Value

	response, err := Cache.Do("GET", "sessions/" + sessionToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return "none"
	}

	if response == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return "none"
	}

	s, _ := redis.String(Cache.Do("GET", "sessions/" + sessionToken))

	return s
}

func GetUser(name string) User {
	var res User
	_ = Coll.FindOne(context.TODO(), bson.D{{"username", name}}).Decode(&res)
	return res
}


func StartMongoDB() {
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, e := mongo.Connect(context.TODO(), clientOptions)

	if e != nil {
		panic(e)
	}

	Coll = client.Database("webpanel").Collection("users")
}

func StartRedis() {
	conn, err := redis.DialURL("redis://localhost")

	if err != nil {
		panic(err)
	}

	Cache = conn
}
