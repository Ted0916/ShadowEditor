package system

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"

	shadow "github.com/tengge1/shadoweditor"
	"github.com/tengge1/shadoweditor/context"
	"github.com/tengge1/shadoweditor/helper"
	"github.com/tengge1/shadoweditor/model"
	"github.com/tengge1/shadoweditor/model/system"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	role := Role{}
	context.Mux.UsingContext().Handle(http.MethodGet, "/api/Role/List", role.List)
	context.Mux.UsingContext().Handle(http.MethodPost, "/api/Role/Add", role.Add)
	context.Mux.UsingContext().Handle(http.MethodPost, "/api/Role/Edit", role.Edit)
	context.Mux.UsingContext().Handle(http.MethodPost, "/api/Role/Delete", role.Delete)
}

// Role 角色控制器
type Role struct {
}

// List 获取列表
func (Role) List(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	pageSize, err := strconv.Atoi(r.FormValue("pageSize"))
	if err != nil {
		pageSize = 20
	}
	pageNum, err := strconv.Atoi(r.FormValue("pageNum"))
	if err != nil {
		pageNum = 1
	}
	keyword := strings.TrimSpace(r.FormValue("keyword"))

	db, err := context.Mongo()
	if err != nil {
		helper.WriteJSON(w, model.Result{
			Code: 300,
			Msg:  err.Error(),
		})
		return
	}

	filter := bson.M{
		"$ne": bson.M{
			"Status": -1,
		},
	}

	if keyword != "" {
		filter1 := bson.M{
			"$regex": bson.M{
				"Name": "/" + keyword + "/i",
			},
		}
		filter = bson.M{
			"$and": bson.A{
				filter,
				filter1,
			},
		}
	}

	skip := int64(pageSize * (pageNum - 1))
	limit := int64(pageSize)
	opts := options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
		Sort: bson.M{
			"ID": -1,
		},
	}

	total, _ := db.Count(shadow.RoleCollectionName, filter)

	docs := bson.A{}
	_ = db.FindMany(shadow.RoleCollectionName, filter, &docs, &opts)

	rows := []system.Role{}

	for _, doc := range docs {
		data := doc.(primitive.D).Map()
		rows = append(rows, system.Role{
			ID:          data["ID"].(string),
			Name:        data["Name"].(string),
			CreateTime:  data["CreateTime"].(time.Time),
			UpdateTime:  data["UpdateTime"].(time.Time),
			Description: data["Description"].(string),
			Status:      data["Status"].(int),
		})
	}

	helper.WriteJSON(w, model.Result{
		Code: 200,
		Msg:  "Get Successfully!",
		Data: map[string]interface{}{
			"total": total,
			"rows":  rows,
		},
	})
}

// Add 添加
func (Role) Add(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	parentID := r.FormValue("ParentID")
	name := strings.Trim(r.FormValue("Name"), " ")
	adminID := r.FormValue("AdminID")

	if name == "" {
		helper.WriteJSON(w, model.Result{
			Code: 300,
			Msg:  "Name is not allowed to be empty.",
		})
		return
	}

	db, err := context.Mongo()
	if err != nil {
		helper.WriteJSON(w, model.Result{
			Code: 300,
			Msg:  err.Error(),
		})
		return
	}

	doc := bson.M{
		"ID":       primitive.NewObjectID(),
		"ParentID": parentID,
		"Name":     name,
		"AdminID":  adminID,
		"Status":   0,
	}

	db.InsertOne(shadow.DepartmentCollectionName, doc)

	helper.WriteJSON(w, model.Result{
		Code: 200,
		Msg:  "Saved successfully!",
	})
}

// Edit 编辑
func (Role) Edit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("ID")
	parentID := r.FormValue("ParentID")
	name := r.FormValue("Name")
	adminID := r.FormValue("AdminID")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		helper.WriteJSON(w, model.Result{
			Code: 300,
			Msg:  "ID is not allowed.",
		})
		return
	}

	if strings.Trim(name, " ") == "" {
		helper.WriteJSON(w, model.Result{
			Code: 300,
			Msg:  "Name is not allowed to be empty.",
		})
		return
	}

	db, err := context.Mongo()
	if err != nil {
		helper.WriteJSON(w, model.Result{
			Code: 300,
			Msg:  err.Error(),
		})
		return
	}

	filter := bson.M{
		"ID": objectID,
	}

	update1 := bson.M{
		"$set": bson.M{
			"ParentID": parentID,
		},
	}
	update2 := bson.M{
		"$set": bson.M{
			"Name": name,
		},
	}
	update3 := bson.M{
		"$set": bson.M{
			"AdminID": adminID,
		},
	}
	update := bson.A{update1, update2, update3}

	_, err = db.UpdateOne(shadow.DepartmentCollectionName, filter, update)
	if err != nil {
		helper.WriteJSON(w, model.Result{
			Code: 300,
			Msg:  err.Error(),
		})
		return
	}

	helper.WriteJSON(w, model.Result{
		Code: 200,
		Msg:  "Saved successfully!",
	})
}

// Delete 删除
func (Role) Delete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("ID")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		helper.WriteJSON(w, model.Result{
			Code: 300,
			Msg:  "ID is not allowed.",
		})
		return
	}

	db, err := context.Mongo()
	if err != nil {
		helper.WriteJSON(w, model.Result{
			Code: 300,
			Msg:  err.Error(),
		})
		return
	}

	filter := bson.M{
		"ID": objectID,
	}

	var doc = bson.M{}
	find, err := db.FindOne(shadow.DepartmentCollectionName, filter, &doc)
	if err != nil {
		helper.WriteJSON(w, model.Result{
			Code: 300,
			Msg:  err.Error(),
		})
		return
	}

	if !find {
		helper.WriteJSON(w, model.Result{
			Code: 300,
			Msg:  "The department is not existed.",
		})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"Status": -1,
		},
	}

	db.UpdateOne(shadow.DepartmentCollectionName, filter, update)

	helper.WriteJSON(w, model.Result{
		Code: 200,
		Msg:  "Delete successfully!",
	})
}
