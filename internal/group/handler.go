package group

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func CreateGroupHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct{ Name string }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid body", http.StatusBadRequest)
			return
		}
		var id int
		err := db.QueryRow("INSERT INTO groups (name) VALUES ($1) RETURNING id", req.Name).Scan(&id)
		if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(Group{ID: id, Name: req.Name})
	}
}

func ListGroupsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name FROM groups")
		if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var groups []Group
		for rows.Next() {
			var g Group
			if err := rows.Scan(&g.ID, &g.Name); err != nil {
				http.Error(w, "DB error", http.StatusInternalServerError)
				return
			}
			groups = append(groups, g)
		}
		json.NewEncoder(w).Encode(groups)
	}
}

func ListGroupMembersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groupID := r.URL.Query().Get("group_id")
		rows, err := db.Query("SELECT user_id FROM group_members WHERE group_id = $1", groupID)
		if err != nil {
			http.Error(w, "DB error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var users []string
		for rows.Next() {
			var user string
			if err := rows.Scan(&user); err != nil {
				http.Error(w, "DB error", http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}
		json.NewEncoder(w).Encode(users)
	}
}
