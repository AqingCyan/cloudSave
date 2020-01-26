package main

import (
	"filestore-server/handler"
	"fmt"
	"net/http"
)

func main()  {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Printf("Failed to start Server, err:%s", err.Error())
	}
}
