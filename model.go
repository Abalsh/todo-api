// model.go

package todo_api

import (
	"database/sql"
)

type goal struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (g *goal) getGoal(db *sql.DB) error {
	return db.QueryRow("SELECT name, price FROM goals WHERE id=$1", g.ID).Scan(&g.Name, &g.Description)
}

func (g *goal) updateGoal(db *sql.DB) error {
	_, err := db.Exec("UPDATE goals SET name=$1, description=$2 WHERE id=$3", g.Name, g.Description, g.ID)
	return err
}

func (g *goal) addGoal(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO goals(name, description) VALUES($1, $2) RETURNING id", g.Name, g.Description).Scan(&g.ID)
	if err != nil {
		return err
	}
	return nil
}

func (g *goal) deleteGoal(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM goals WHERE id=$1", g.ID)
	return err
}

func getGoals(db *sql.DB, start, count int) ([]goal, error) {
	rows, err := db.Query("SELECT id, name , description FROM goals LIMIT $1 OFFSET $2", count, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	goals := []goal{}

	for rows.Next() {
		var g goal
		if err := rows.Scan(&g.ID, &g.Name, &g.Description); err != nil {
			return nil, err
		}
		goals = append(goals, g)
	}
	return goals, nil
}
