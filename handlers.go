package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func HandleDefault(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Default Route!")
}

func HandleGetInventory(w http.ResponseWriter, r *http.Request) {
	var err error
	data := []byte("Get")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
		return
	}
	w.Write(data)
}

func HandleGetInventoryList(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Handle GET List")
}

func HandlePostInventory(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "Handle POST ERROR", err)
		return
	}

	var item Item
	fmt.Println(string(data))
	err = json.Unmarshal(data, &item)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "Handle POST ERROR", err)
		return
	}

	err = DB.WriteItemToInventory(item)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, "Handle POST ERROR", err)
		return
	}
	fmt.Fprintln(w, "Handle POST", string(data), item)
}

func HandlePatchInventory(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Handle Patch")
}

func HandleDeleteInventory(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Handle Delete")
}
