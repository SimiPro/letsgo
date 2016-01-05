package main

import (
	"net/http"
	"mime"
)

func init() {
	mime.AddExtensionType(".dart","application/dart")

	fs := http.FileServer(http.Dir("app/build/web"))
	http.Handle("/", fs)
	//router.Run(":8080")
	//http.ListenAndServe(":8080", nil)
}
