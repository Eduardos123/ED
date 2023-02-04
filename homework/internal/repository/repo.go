package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"test/internal/models"

	"github.com/go-chi/chi/v5"
)

type Repo struct {
	db *sql.DB
}

func CreateRepo(db *sql.DB) *Repo {
	return &Repo{
		db: db,
	}
}


func (y *Repo) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Server2")
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var usr models.User
	err = json.Unmarshal(b, &usr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	id, err := y.AddUser(usr)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	w.Write([]byte(strconv.Itoa(id)))
}

func (y *Repo) AddUser(usr models.User) (int, error) {
	_, err := y.db.Exec("INSERT user(name, age) VALUES (?, ?)", usr.Name, usr.Age)
	if err != nil {

		return 0, err
	}
	res, err := y.db.Query("SELECT id FROM user WHERE name = ? AND age = ?", usr.Name, usr.Age)
	if err != nil {

		return 0, err
	}
	var id int
	for res.Next() {
		err = res.Scan(&id)
		if err != nil {
			return 0, err
		}
	}
	return id, nil
}

func (y *Repo) UpdateAge(w http.ResponseWriter, r *http.Request) {
	var age models.Age
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = json.Unmarshal(b, &age)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, errors.New("Неверный ID пользователя").Error(), 500)
		return
	}
	id1, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	ageI, err := strconv.Atoi(age.Age)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, err = y.db.Query("UPDATE  user set age = ? WHERE id = ?", ageI, id1)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("Возраст пользователя успешно обновлен!"))
}

func (y *Repo) MakeFriends(w http.ResponseWriter, r *http.Request) {
	var f models.Frds
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = json.Unmarshal(b, &f)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_, err = y.db.Query("INSERT INTO friends (id_user,id_friend) VALUES (?,?)", f.SourceID, f.TargetID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_, err = y.db.Query("INSERT INTO friends (id_user,id_friend) VALUES (?,?)", f.TargetID, f.SourceID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(f.SourceID + " и " + f.TargetID + " теперь друзья!"))
}

func (y *Repo) GetFriends(w http.ResponseWriter, r *http.Request) {
	var friendsNames string
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, errors.New("Такого пользователя нет.").Error(), 500)
		return
	}
	id1, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res, err := y.db.Query("SELECT user.name FROM friends,user WHERE user.id=friends.id_friend AND friends.id_user = ? GROUP BY user.name", id1)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var name string

	for res.Next() {
		err = res.Scan(&name)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		friendsNames = friendsNames + name + ", "
	}
	if friendsNames == "" {
		w.WriteHeader(200)
		w.Write([]byte("У пользователя нет друзей"))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(friendsNames))
}

func (y *Repo) DeleteUser(w http.ResponseWriter, r *http.Request) {

	var name string
	id := chi.URLParam(r, "id")

	if id == "" {
		http.Error(w, errors.New("Такого пользователя нет.").Error(), 500)
		return
	}
	id1, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_, err = y.db.Query("DELETE FROM friends WHERE id_user = ? OR id_friend = ?", id1, id1)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	res, err := y.db.Query("SELECT name FROM user WHERE id=?", id1)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	for res.Next() {
		err = res.Scan(&name)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	if name == "" {
		w.WriteHeader(200)
		w.Write([]byte("Нет такого пользователя"))
		return
	}
	_, err = y.db.Query("DELETE FROM user WHERE id = ?", id1)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(name))
}
