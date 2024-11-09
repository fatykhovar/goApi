package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)




// type ApiError struct {
// 	Error string `json:"error"`
// }
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, err.Error())
		}
	}
}


type APIServer struct {
	listenAddr string
	store      Storage
}

func (s *APIServer) Run(){
	router := mux.NewRouter()
    router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
    router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleAccount))

	http.ListenAndServe(s.listenAddr, router)
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccount(w, r)
	case "POST":
		return s.handleNewAccount(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	default:
		return fmt.Errorf("method not allowed %s", r.Method)
	}
	//return APIServer.NewRoute().Path(path).HandlerFunc(f)
	
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	
	id, err := getID(r)

	if err != nil {
		return err
	}

	account, err := s.store.GetAccountByID(id)
	
	if err != nil {
		return err
	}
	log.Println("got account")
	
	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleNewAccount(w http.ResponseWriter, r *http.Request) error {
	acc1 := NewAccount("Regina", "Fatykhova")
	if err := s.store.CreateAccount(acc1); err != nil {
		return err
	}
	// log.Println("new account")
	return WriteJSON(w, http.StatusOK, acc1)
}


func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)

	if err != nil {
		return err
	}
	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}
	// log.Println("new account")
	return WriteJSON(w, http.StatusOK, id)
}

func getID(r *http.Request) (int, error){
	vars := mux.Vars(r)
	category := vars["id"]
	id, err := strconv.Atoi(category)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", category)
	}
	return id, nil
}
