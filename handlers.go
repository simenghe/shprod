package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func HandleDefault(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Default Route!")
}

func HandleGetInventory(w http.ResponseWriter, r *http.Request) {
	var err error
	items, err := DB.ReadItems()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
		return
	}
	data, err := json.Marshal(items)
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
		log.Println("Handle POST ERROR", err)
		return
	}

	var item Item
	fmt.Println(string(data))
	err = json.Unmarshal(data, &item)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Handle POST ERROR", err)
		return
	}

	id, err := DB.WriteItemToInventory(item)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Handle POST ERROR", err)
		return
	}
	item.ID = id // Assign the ID once written
	err = json.NewEncoder(w).Encode(&item)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func HandlePutInventory(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
		return
	}
	_, err = DB.EditItem(item)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
		return
	}
	fmt.Printf("Received fields: %+v\n", item)
	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		fmt.Fprintln(w, err)
	}
}

// HandleDeleteInventory deletes an item by id.
func HandleDeleteInventory(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fmt.Println(r.URL)
	queryMap := r.URL.Query()
	if !queryMap.Has("id") {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "No id in URL query: ", r.URL)
		return
	}

	id := queryMap.Get("id")
	err := DB.DeleteItem(id)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, "Error in deleting item with id: ", id)
		w.WriteHeader(500)
		return
	}
	fmt.Fprintln(w, "Deleted objectid:", id)
}
