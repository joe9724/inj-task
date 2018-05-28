package main

import (
	"net/http"
	"log"
	_"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
	"encoding/json"
	"inj-task/utils"
	"inj-task/model"
)


func main() {
	//http.HandleFunc("/task/status", Status)
	http.HandleFunc("/task/do", Do)
	//http.HandleFunc("/task/today",GetTaskToday)
	http.HandleFunc("/task/list",GetTaskList)
	err := http.ListenAndServe(":1067", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	fmt.Println("start server at 0.0.0.0:1067")

}

func GetTaskToday (w http.ResponseWriter, req *http.Request) {
	db,err := utils.OpenConnection()
	if err!=nil{
		fmt.Println(err.Error())
	}
	defer db.Close()

	var msg string
	var code int64
	var response utils.Response
	//签到,先检查是否已签到
	var list []model.TaskModel
	db.Raw("select id,user_id,`status`,DATE_FORMAT(create_time,'%Y-%m-%d %h:%m:%s') as create_time,type,name,value,value_unit,target_id,event_id from btk_Task where user_id=? and to_days(create_time) = to_days(now()",utils.DeUserID(req.URL.Query().Get("euid"))).Find(&list)
	if len(list)==0{ //当天还没有签到，可以正常签到
		//db.Exec("insert into btk_Check(user_id) values(?)",utils.DeUserID(req.URL.Query().Get("euid")))
		msg = "今天没有做任务"
		code = 200
	}else{
		msg = "今天任务记录"
		code = 300
		response.Data = list

	}
	response.Msg = msg
	response.Code = code

	//fmt.Println("start is",start)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", "A Go Web Server")
	data, err := json.Marshal(response)
	if err != nil {
		log.Fatal("err get data: ", err)
	}

	w.Write(data)

}

func GetTaskList (w http.ResponseWriter, req *http.Request) {
	db,err := utils.OpenConnection()
	if err!=nil{
		fmt.Println(err.Error())
	}
	defer db.Close()

	var msg string
	var code int64
	var response utils.Response
	//签到,先检查是否已签到
	var list []model.TaskModel
	db.Raw("SELECT btk_Task_Set.name,btk_Task_Set.type,btk_Task_Set.value,btk_Task_Set.value_unit,btk_Task_Set.target_id,btk_Task_Set.event_id,DATE_FORMAT(btk_Task.create_time,'%Y-%m-%d %h:%m:%s') as create_time,btk_Task.user_id FROM `btk_Task_Set` left join btk_Task on btk_Task_Set.`name` = btk_Task.`name` and btk_Task.user_id = ?  and to_days(btk_Task.create_time) = to_days(now())",utils.DeUserID(req.URL.Query().Get("euid"))).Find(&list)
	response.Data = list
	response.Msg = msg
	response.Code = code

	//fmt.Println("start is",start)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", "A Go Web Server")
	data, err := json.Marshal(response)
	if err != nil {
		log.Fatal("err get data: ", err)
	}

	w.Write(data)

}

func Status(w http.ResponseWriter, req *http.Request) {
	db,err := utils.OpenConnection()
	if err!=nil{
		fmt.Println(err.Error())
	}
	defer db.Close()
	var msg string
	var code int64
	var response utils.Response
	//签到,先检查是否已签到
	var list []model.TaskModel
	db.Raw("select id,user_id,create_time,`status` from btk_Check where user_id=? and  to_days(create_time) = to_days(now())",utils.DeUserID(req.URL.Query().Get("euid"))).Find(&list)
	if len(list)==0{ //当天还没有签到，可以正常签到
		//db.Exec("insert into btk_Check(user_id) values(?)",utils.DeUserID(req.URL.Query().Get("euid")))
		msg = "任务未完成"
		code = 200
	}else{
		msg = "任务已完成"
		code = 300
	}
	response.Msg = msg
	response.Code = code

	//fmt.Println("start is",start)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", "A Go Web Server")
	data, err := json.Marshal(response)
	if err != nil {
		log.Fatal("err get data: ", err)
	}

	w.Write(data)


}

func Do(w http.ResponseWriter, req *http.Request) {
	db,err := utils.OpenConnection()
	if err!=nil{
		fmt.Println(err.Error())
	}
	defer db.Close()
	var msg string
	var code int64
	var response utils.Response
	//签到,先检查是否已签到
	var list []model.TaskModel
	db.Raw("select id,user_id,create_time,`status`,name from btk_Task where user_id=? and to_days(create_time) = to_days(now()) and name=?",utils.DeUserID(req.URL.Query().Get("euid")),req.URL.Query().Get("name")).Find(&list)
	fmt.Println("select id,user_id,create_time,`status` from btk_Task where user_id=? and to_days(create_time) = to_days(now()) and name=?",utils.DeUserID(req.URL.Query().Get("euid")),req.URL.Query().Get("name"))
	if len(list)==0{ //当天还没有签到，可以正常签到
		db.Exec("insert into btk_Task(user_id,name) values(?,?)",utils.DeUserID(req.URL.Query().Get("euid")),req.URL.Query().Get("name"))
		msg = "成功"
		code = 200
	}else{
		msg = "已做过了,亲"
		code = 300
	}
	response.Msg = msg
	response.Code = code

	//fmt.Println("start is",start)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", "A Go Web Server")
	data, err := json.Marshal(response)
	if err != nil {
		log.Fatal("err get data: ", err)
	}

	w.Write(data)
}