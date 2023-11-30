package adapter

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// )

// func getHandler(w http.ResponseWriter, r *http.Request) {
// 	key := r.URL.Path[len("/kv/"):]
// 	value, err := Get(key)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write([]byte(value))
// }

// func setHandler(w http.ResponseWriter, r *http.Request) {
// 	var kv struct {
// 		Key   string `json:"key"`
// 		Value string `json:"value"`
// 	}
// 	err := json.NewDecoder(r.Body).Decode(&kv)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	err = Set(kv.Key, kv.Value)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write([]byte("Set successful"))
// }

// func deleteHandler(w http.ResponseWriter, r *http.Request) {
// 	key := r.URL.Path[len("/kv/"):]
// 	err := Delete(key)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.Write([]byte("Delete successful"))
// }

// func ListenHttp() {
// 	http.HandleFunc("/kv/", getHandler)
// 	http.HandleFunc("/kv", setHandler)
// 	http.HandleFunc("/kv/", deleteHandler)
// 	fmt.Println("Server is running at http://localhost:8080")
// 	http.ListenAndServe(":8080", nil)
// }
