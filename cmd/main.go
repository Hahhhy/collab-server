package main

import (
	"TO/internal/service"
	"TO/internal/storage"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	//1.初始化数据库连接
	//加载驱动
	//验证 DSN 格式
	//返回一个 *sql.DB 对象（连接池的句柄）
	db, err := sql.Open("postgres", "your-dsn-here")
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	defer db.Close()

	//2。初始化storage层
	store := storage.NewStorage(db)

	//3.初始化广播服务
	broadcaster := service.NewBroadcastService()

	//4.初始化协作服务
	collab := service.NewCollabService(*store, broadcaster)

	//5.启动HTTP服务器，注册处理函数
	http.HandleFunc("/collab", func(w http.ResponseWriter, r *http.Request) {
		//处理HTTP请求
		var req service.EditRequest
		//解析请求参数
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//调用协作服务处理请求
		resp, err := collab.HandleEdit(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//返回处理结果
		json.NewEncoder(w).Encode(resp)
		http.Error(w, "ok", http.StatusOK)
		http.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {})
		http.ListenAndServe(":8080", nil)
	})

}
