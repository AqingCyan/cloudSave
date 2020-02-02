package main

import (
	"filestore-server/handler"
	"fmt"
	"net/http"
)

func main()  {
	// 静态资源处理
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/query", handler.FileQueryHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)

	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Printf("Failed to start Server, err:%s", err.Error())
	}
}
