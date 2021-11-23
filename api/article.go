package api

import (
	"encoding/json"
	"github.com/qiniu/qmgo/operator"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"blog_server/db"
	"io/ioutil"
	"net/http"
)

type MetaInfo struct {
	Views    uint32 `bson:"views" json:"views"`
	Comments uint32 `bson:"comments" json:"comments"`
	Likes    uint32 `bson:"likes" json:"likes"`
}

type Comment1 struct {
	Avatar     string `bson:"avatar" json:"avatar"`
	Name       string `bson:"name" json:"name"`
	Type       uint32 `bson:"type" json:"type"`
	CreateTime string `bson:"createtime" json:"createtime"`
}

type User struct {
	Name   string `bson:"name" json:"name"`
	Avatar string `bson:"avatar" json:"avatar"`
	Type   uint32 `bson:"type" json:"type"`
}

type OtherComment struct {
	Id         primitive.ObjectID `bson:"_id" json:"_id"`
	User       User               `bson:"user" json:"user"`
	CreateTime string             `bson:"createtime" json:"createtime"`
	Content    string             `bson:"content" json:"content"`
	ToUser     User               `bson:"touser" json:"touser"`
}

type Comment struct {
	Id            primitive.ObjectID `bson:"_id" json:"_id"`
	User          User               `bson:"user" json:"user"`
	CreateTime    string             `bson:"createtime" json:"createtime"`
	Content       string             `bson:"content" json:"content"`
	OtherComments []OtherComment     `bson:"othercomments" json:"othercomments"`
}

type Article struct {
	Id         primitive.ObjectID `bson:"_id" json:"_id"`
	Title      string             `bson:"title"  json:"title"`
	Author     string             `bson:"author" json:"author"`
	Desc       string             `bson:"desc" json:"desc"`
	Meta       MetaInfo           `bson:"meta" json:"meta"`
	Tags       []string           `bson:"tags" json:"tags"`
	Comments   []Comment          `bson:"comments" json:"comments"`
	LikeUsers  []string           `bson:"likeusers" json:"likeusers"`
	CreateTime string             `bson:"createtime" json:"createtime"`
	UpdateTime string             `bson:"updatetime" json:"updatetime"`
	Content    string             `bson:"content" json:"content"`
}

type ArticleDatas struct {
	Count uint32    `json:"count"`
	List  []Article `json:"list"`
}

func ArticleList(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json
	res := Response{}
	var param map[string]interface{}
	body, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(body, &param)
	log.Info(param)

	col := db.NewCollection("articles")

	list := []Article{}
	col.FindAll(bson.M{}).All(&list)
	log.Infof("one: %v", list)

	var articles = &ArticleDatas{}
	articles.Count = uint32(len(list))
	articles.List = list
	res.RetContent = articles
	res.SetErrorCode(Succeed)
	json.NewEncoder(w).Encode(res)
}

func ArticleQuery(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	res := Response{}
	var param map[string]interface{}
	body, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(body, &param)
	log.Info(param)

	id, ok := param["_id"].(string)
	if !ok {
		res.SetErrorCode(ErrParam)
		json.NewEncoder(w).Encode(res)
	}
	article := Article{}
	col := db.NewCollection("articles")
	if !primitive.IsValidObjectID(id) {
		res.SetErrorCode(ErrParam)
		json.NewEncoder(w).Encode(res)
	}
	objectid, _ := primitive.ObjectIDFromHex(id)
	ret := col.FindOne(bson.M{"_id": objectid})
	ret.One(&article)
	log.Infof("article: %v", article)
	res.RetContent = article
	res.SetErrorCode(Succeed)
	json.NewEncoder(w).Encode(res)
}

func ArticleDelete(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	res := Response{}
	var param map[string]interface{}
	body, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(body, &param)

	id, ok := param["_id"].(string)
	if !ok {
		res.SetErrorCode(ErrParam)
		json.NewEncoder(w).Encode(res)
	}
	col := db.NewCollection("articles")
	if !primitive.IsValidObjectID(id) {
		res.SetErrorCode(ErrParam)
		json.NewEncoder(w).Encode(res)
	}
	objectid, _ := primitive.ObjectIDFromHex(id)
	err := col.DeleteOne(bson.M{"_id": objectid})
	if err != nil {
		res.SetErrorCode(ErrDeleteFail, err)
		json.NewEncoder(w).Encode(res)
		return
	}
	res.SetErrorCode(Succeed)
	json.NewEncoder(w).Encode(res)
}

// 参数不为空就更新
func ArticleUpdate(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	res := Response{}
	var param map[string]interface{}
	body, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(body, &param)
	log.Info(string(body))
	id, ok := param["_id"].(string)
	if !ok {
		res.SetErrorCode(ErrParam)
		json.NewEncoder(w).Encode(res)
		return
	}
	log.Info(param)
	data, condition := getUpdateParam(param)
	col := db.NewCollection("articles")
	if !primitive.IsValidObjectID(id) {
		res.SetErrorCode(ErrParam)
		json.NewEncoder(w).Encode(res)
		return
	}
	objectid, _ := primitive.ObjectIDFromHex(id)

	condition["_id"] = objectid
	log.Info(data)
	log.Info(condition)
	err := col.UpdateOne(condition, data)
	if err != nil {
		res.SetErrorCode(ErrUpdateFail, err)
		json.NewEncoder(w).Encode(res)
		return
	}
	res.SetErrorCode(Succeed)
	json.NewEncoder(w).Encode(res)

}

func ArticleCreate(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")             //返回数据格式是json

	res := Response{}
	var param map[string]interface{}
	body, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(body, &param)
	log.Info(param)
	if len(param) == 0 {
		res.SetErrorCode(ErrParam)
		json.NewEncoder(w).Encode(res)
		return
	}

	var artcile = &Article{
		Title:      getCreateParam(param, "title").(string),
		Author:     getCreateParam(param, "author").(string),
		Desc:       getCreateParam(param, "desc").(string),
		CreateTime: getCreateParam(param, "createtime").(string),
		UpdateTime: getCreateParam(param, "updatetime").(string),
		Content:    getCreateParam(param, "content").(string),
		Comments:   make([]Comment, 0, 0),
		Meta:       MetaInfo{},
	}
	tags := make([]string, 0, 5)
	if temp, ok := param["tags"].([]interface{}); ok {
		for _, v := range temp {
			tags = append(tags, v.(string))
		}
	}

	artcile.Tags = tags
	artcile.Id = primitive.NewObjectID()
	log.Info(artcile)

	col := db.NewCollection("articles")
	ret, err := col.InsertOne(artcile)
	if err != nil {
		log.Info(ret, err)
		res.SetErrorCode(ErrInsertFail, err)
		json.NewEncoder(w).Encode(res)
	}
	res.SetErrorCode(Succeed)
	json.NewEncoder(w).Encode(res)

}

// 返回 要更新的数据和额外的定位条件
// 目前只有comment需要返回额外的条件
func getUpdateParam(param map[string]interface{}) (bson.M, bson.M)  {
	data := bson.M{}
	condtion := bson.M{}
	operation, ok := param["operation"].(string)
	if !ok {
		return bson.M{}, condtion
	}
	var list = []string{
		"title", "author", "desc", "createtime", "updatetime", "content",
	}
	for _, l := range list {
		if temp, ok := param[l].(string); ok {
			data[l] = temp
		}
	}
	if temp, ok := param["meta"].(interface{}); ok {
		bytes, err := json.Marshal(temp)
		if err != nil {
			log.Error(err)
		}
		log.Info(string(bytes))
		meta := &MetaInfo{}
		json.Unmarshal(bytes, meta)
		log.Info(meta)
		return updateMeta(*meta) , condtion
	}

	if temp, ok := param["tags"].([]interface{}); ok {
		tags := make([]string, 0, 5)
		for _, tag := range temp {
			tags = append(tags, tag.(string))
		}
		if operation == "update" {
			data["tags"] = tags
		} else {
			return updateArray(tags, operation, "tags"), condtion
		}

	}

	if temp, ok := param["likeusers"].([]interface{}); ok {
		likeusers := make([]string, 0, 5)
		for _, likeuser := range temp {
			likeusers = append(likeusers, likeuser.(string))
		}
		if operation == "update" {
			data["likeusers"] = likeusers
		} else {
			return updateArray(likeusers, operation, "likeusers"), condtion
		}
	}

	if temp, ok := param["comments"].([]interface{}); ok {
		comments := make([]*Comment, 0, 5)
		for _, comment := range temp {
			bytes, err := json.Marshal(comment)
			if err != nil {
				log.Error(err)
			}
			log.Info(string(bytes))
			c := &Comment{}
			json.Unmarshal(bytes, c)
			log.Info(c)
			comments = append(comments, c)
		}

		if operation == "update" {
			data["comments"] = comments
		} else {
			return updateComment(comments[0], operation)
		}
	}

	update := bson.M{
		operator.Set: data,
	}
	return update, condtion
}

func getCreateParam(param map[string]interface{}, filed string) interface{} {
	if _, ok := param[filed]; !ok {
		return ""
	}
	return param[filed]
}

func updateMeta(meta MetaInfo) bson.M {
	log.Info(meta)
	ret := bson.M{}
	if meta.Views != 0 {
		ret["meta.views"] = meta.Views
	}
	if meta.Comments != 0 {
		ret["meta.comments"] = meta.Comments
	}
	if meta.Likes != 0 {
		ret["meta.likes"] = meta.Likes

	}
	return bson.M{operator.Inc: ret}
}

// 更新数组
// operation: [add, del]
// property:[tags, comments, likeusers]

func updateArray(arrs interface{}, operation string, property string) bson.M {
	if operation == "add" {
		return bson.M{operator.Push: bson.M{property: bson.M{operator.Each: arrs}}}
	} else if operation == "del" {
		return bson.M{operator.Pull: bson.M{property: bson.M{operator.In: arrs}}}
	}
	return bson.M{}
}

// 包括增删，不支持修改
func updateComment(comment *Comment, operation string) (bson.M,  bson.M) {
	// add 分两种，
	//1.增加正常评论
	//2.回复别人评论
	condition := bson.M{}
	log.Info(comment.Id.IsZero())
	if operation == "add" {
		if comment.Id.IsZero() { // 第一种
			comment.Id = primitive.NewObjectID()
			log.Info(comment.Id)
			comment.OtherComments = make([]OtherComment, 0, 0) // 清空防止输入非法内容
			return bson.M{operator.Push: bson.M{"comments": comment}}, condition
		} else { // 第二种
			objectid, _ := primitive.ObjectIDFromHex(comment.Id.Hex())
			id := primitive.NewObjectID()
			otherComment := comment.OtherComments[0]
			otherComment.Id = id
			// 注意数组需要由于不知道位置，需要先加$占位，查询条件不需要占位符
			return bson.M{operator.Push: bson.M{"comments.$.othercomments": otherComment}}, bson.M{"comments._id": objectid}
		}

	} else if operation == "del" {
		log.Info(comment)
		if comment.Id.IsZero() { // 为空不删除
			return bson.M{}, condition
		}
		if len(comment.OtherComments) > 0 {
			objectid, _ := primitive.ObjectIDFromHex(comment.Id.Hex())
			id, _ := primitive.ObjectIDFromHex(comment.OtherComments[0].Id.Hex())
			return bson.M{operator.Pull: bson.M{"comments.$.othercomments": bson.M{"_id":id}}}, bson.M{"comments._id": objectid}
		} else {
			id, _ := primitive.ObjectIDFromHex(comment.Id.Hex())
			return bson.M{operator.Pull: bson.M{"comments": bson.M{"_id":id}}}, condition
		}

	}
	return bson.M{}, condition
}
