package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/0x21F/goHttpTesting/models"
)

type Handler struct {
	rw          sync.RWMutex
	Peeps       map[uint]*models.Person
	FollowGraph models.Graph
	idx         uint
}

func main() {
	fmt.Println("================================================")
	fmt.Println("                 shitty app lol                 ")
	fmt.Println("================================================")

	mux := http.NewServeMux()
	hands := &Handler{}
	hands.init()

	mux.Handle("GET /person/{id}", LogPath(http.HandlerFunc(hands.getPerson)))
	mux.Handle("POST /person", LogPath(http.HandlerFunc(hands.addPerson)))
	mux.Handle("PUT /person/{id}", LogPath(http.HandlerFunc(hands.editPerson)))
	mux.Handle("DELETE /person/{id}", LogPath(http.HandlerFunc(hands.delPerson)))

	http.ListenAndServe(":8080", mux)
}

func (h *Handler) init() {
	h.Peeps = make(map[uint]*models.Person)
}

func (h *Handler) addPerson(w http.ResponseWriter, r *http.Request) {
	h.rw.Lock()
	defer h.rw.Unlock()

	person := &models.Person{}
	err := json.NewDecoder(r.Body).Decode(person)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		h.Peeps[h.idx] = person
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(r.URL.Path + "/" + fmt.Sprint(h.idx)))
		h.idx++
	}
}

func (h *Handler) delPerson(w http.ResponseWriter, r *http.Request) {
	h.rw.Lock()
	defer h.rw.Unlock()

	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}

	if h.Peeps[uint(id)] == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("doesn't exist"))
		return
	}

	h.Peeps[uint(id)] = nil

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted"))
}

func (h *Handler) editPerson(w http.ResponseWriter, r *http.Request) {
	h.rw.Lock()
	defer h.rw.Unlock()

	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}

	if h.Peeps[uint(id)] == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("doesn't exist"))
		return
	}

	person := &models.Person{}
	err = json.NewDecoder(r.Body).Decode(person)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	} else {
		person.Id = uint(id)
		h.Peeps[uint(id)] = person
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Edited\n"))
	}

}

func (h *Handler) getPerson(w http.ResponseWriter, r *http.Request) {
	h.rw.RLock()
	defer h.rw.RUnlock()

	id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}

	if h.Peeps[uint(id)] == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("doesn't exist"))
		return
	}

	res, err := json.Marshal(h.Peeps[uint(id)])

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshalling"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func LogPath(next http.Handler) http.Handler {
	res := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})

	return res
}
