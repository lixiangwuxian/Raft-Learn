package adapter

import (
	"net/http"

	"lxtend.com/m/logger"
	"lxtend.com/m/store"
)

var lOnMsg func(command string)
var storage *store.InMemoryLogStore

func userHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	command := query.Get("command")
	logger.Glogger.Info("user command:%s", command)
	lOnMsg(command)
}

func getAllCommands(w http.ResponseWriter, r *http.Request) {
	var commands string
	for _, command := range storage.GetSince(0) {
		commands += command.Command + "\n"
	}
	w.Write([]byte(commands))
}

func ListenHttp(port string, onMsg func(command string), logStore *store.InMemoryLogStore) {
	http.HandleFunc("/", userHandler)
	http.HandleFunc("/commands", getAllCommands)
	lOnMsg = onMsg
	storage = logStore
	logger.Glogger.Info("Server is running at http://localhost:" + port)
	go http.ListenAndServe(":"+port, nil)
}
