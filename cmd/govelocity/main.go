package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
)

type PushEvent struct {
	Ref string `json:"ref"`
}

func main() {
	http.HandleFunc("/webhook", handleWebhook)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	eventType := r.Header.Get("X-GitHub-Event")
	if eventType != "push" {
		http.Error(w, "Unsupported event type", http.StatusBadRequest)
		return
	}

	var pushEvent PushEvent
	if err := json.NewDecoder(r.Body).Decode(&pushEvent); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	if pushEvent.Ref != "refs/heads/master" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := deploy(); err != nil {
		http.Error(w, "Failed to deploy", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deploy() error {
	cmd := exec.Command("bash", "deploy.sh")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
