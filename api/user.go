package api

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
)

type UserInfo struct {
	Id       primitive.ObjectID `bson:"_id" json:"_id"`
	Name     string `bson:"name"  json:"name"`
	Email    string `bson:"email"  json:"email"`
	Password string `bson:"password"  json:"password"`
	Phone    string `bson:"phone"  json:"phone"`
	Desc     string `bson:"desc"  json:"desc"`
	Avatar     string `bson:"avatar"  json:"avatar"`
}

func UserLogin(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
	res := Response{}
	var param map[string]interface{}
	body, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(body, &param)
	log.Info(param)

	email, ok := param["email"].(string)
	if !ok {
		res.SetErrorCode(ErrParam)
		json.NewEncoder(w).Encode(res)
	}
	pwd, ok := param["password"].(string)
	if !ok {
		res.SetErrorCode(ErrParam)
		json.NewEncoder(w).Encode(res)
	}
	log.Info(email, pwd)
	ret := UserInfo {
		primitive.ObjectID{},
		"123",
		"123",
		"123",
		"123",
		"123",
		"123",
	}
	res.RetContent = ret
	res.SetErrorCode(Succeed)
	json.NewEncoder(w).Encode(res)
}
