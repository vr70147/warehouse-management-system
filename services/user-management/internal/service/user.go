package service

import (
	"context"
	"net/http"
	"os"
	"user-management/internal/model"

	"github.com/jackc/pgx/v4"
)

func CreateUser(conn *pgx.Conn, user model.User) error {
	sql := `INSERT INTO users (name, email, age, role, phone, street, city, password)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := conn.Exec(context.Background(), sql, user.Name, user.Email, user.Age, user.Role, user.Phone, user.Street, user.City, user.Password)
	return err
}

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func LoginUser(w http.ResponseWriter, r *http.Request) {

}
