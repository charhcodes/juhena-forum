// package main

// import (
// 	"net/http"
// 	"fmt"
// )

// func main(){
// 	fmt.Println("Juhena")
// 	http.HandleFunc("/", home)
// 	// fmt.Println("listening on port :8080")
// 	http.ListenAndServe(":8080", nil)
// }

// func home(w http.ResponseWriter, r *http.Request){
// 	w.Header().Set("Content-Type", "text/html")
// 	w.Write([]byte("<h2>This is my home page</h2>")) }