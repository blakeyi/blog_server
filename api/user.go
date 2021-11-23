package api

import (
	"blog_server/db"
	"encoding/json"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
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
		return
	}
	pwd, ok := param["password"].(string)
	if !ok {
		res.SetErrorCode(ErrParam)
		json.NewEncoder(w).Encode(res)
		return
	}
	log.Info(email, pwd)
	user := UserInfo{}
	col := db.NewCollection("users")
	ret := col.FindOne(bson.M{"email": email})
	if cnt, _ := ret.Count(); cnt  == 0  {
		res.SetErrorCode(ErrParam, errors.New("用户不存在"))
		json.NewEncoder(w).Encode(res)
		return
	}
	ret.One(&user)
	if user.Password != pwd {
		res.SetErrorCode(ErrParam, errors.New("密码不正确"))
		json.NewEncoder(w).Encode(res)
		return
	}
	res.RetContent = user
	res.SetErrorCode(Succeed)
	json.NewEncoder(w).Encode(res)
}

func UserRegister(w http.ResponseWriter, req *http.Request) {
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

	name, ok := param["name"].(string)
	if !ok {
		res.SetErrorCode(ErrParam)
		json.NewEncoder(w).Encode(res)
	}
	phone, ok := param["phone"].(string)
	if !ok {
		res.SetErrorCode(ErrParam)
		json.NewEncoder(w).Encode(res)
	}
	desc, ok := param["desc"].(string)
	if !ok {
		res.SetErrorCode(ErrParam)
		json.NewEncoder(w).Encode(res)
	}

	id := primitive.NewObjectID()
	data := UserInfo{
		id,
		name,
		email,
		pwd,
		phone,
		desc,
		"../assets/user.png",
	}
	// 检查是否已经注册
	if isUserExist(data) {
		res.SetErrorCode(ErrInsertFail, errors.New("你已经注册过了，无需重复注册"))
		json.NewEncoder(w).Encode(res)
		return
	}
	col := db.NewCollection("users")
	ret, err := col.InsertOne(data)
	if err != nil {
		log.Info(ret, err)
		res.SetErrorCode(ErrInsertFail, err)
		json.NewEncoder(w).Encode(res)
		return
	}
	res.SetErrorCode(Succeed)
	json.NewEncoder(w).Encode(res)
}

func getAllUser() []*UserInfo {
	col := db.NewCollection("users")
	list := []*UserInfo{}
	col.FindAll(bson.M{}).All(&list)
	log.Infof("one: %v", list)
	return list
}

func isUserExist(info UserInfo) bool {
	list := getAllUser()
	for _, v := range list {
		if v.Email == info.Email {
			return true
		}
	}
	return false
}
