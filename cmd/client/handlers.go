package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
)

const (
	authCookieName = "auth"
)

type user struct {
	UserFullName string `db:"user_full_name"`
	Email        string `db:"email"`
	Password     string `db:"password"`
}

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authClaims struct {
	UserFullName string `json:"user_full_name"`
	jwt.StandardClaims
}

type customer struct {
	CustomerId    int       `db:"customer_id"`
	AccountNumber string    `db:"account_number"`
	LastName      string    `db:"last_name"`
	FirstName     string    `db:"first_name"`
	MiddleName    string    `db:"middle_name"`
	DateOfBirth   time.Time `db:"date_of_birth"`
	TaxId         string    `db:"tax_id"`
	Status        string    `db:"status"`
}

type dashboardData struct {
	UserFullName string
	Customers    []*customer
	Dob          []*string
}

type userStatusUpdate struct {
	CustomerId string `json:"customerId"`
	Status     string `json:"status"`
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func dashboard(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userFullName, err := authByCookie(w, r)
		if err != nil {
			return
		}
		customerData, err := selectCustomers(db, userFullName)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		var dobStrings []*string
		for _, cust := range customerData {
			dobString := cust.DateOfBirth.Format("02-01-2006")
			dobStrings = append(dobStrings, &dobString)
		}

		data := dashboardData{
			UserFullName: userFullName,
			Customers:    customerData,
			Dob:          dobStrings, // Сохраняем строки с датами в новом поле
		}

		ts, err := template.ParseFiles("pages/dashboard.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		log.Println("Request completed successfully")
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("pages/login.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	log.Println("Request completed successfully")
}

func logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    authCookieName,
			Path:    "/",
			Expires: time.Now().AddDate(0, 0, -1),
		})
		w.WriteHeader(http.StatusOK)
		log.Println("Request completed successfully")
	}
}

func authentication(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		var req userRequest
		err = json.Unmarshal(reqData, &req)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		user, err := CheckUser(db, req.Email, req.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Incorrect password or email", http.StatusUnauthorized)
				return
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				log.Println(err)
				return
			}
		}
		const tokenExpireDuration = time.Hour * 24
		claims := &authClaims{
			UserFullName: user.UserFullName,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(tokenExpireDuration).Unix(),
				IssuedAt:  time.Now().Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("FDsh43nd650zDfkdjKDSfd45DdfJSHdsj42"))
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    authCookieName,
			Value:   tokenString,
			Path:    "/",
			Expires: time.Now().AddDate(0, 0, 1),
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func authByCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie(authCookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			log.Println(err)
			return "", err
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err)
		return "", err
	}
	tokenString := cookie.Value
	claims := &authClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("FDsh43nd650zDfkdjKDSfd45DdfJSHdsj42"), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			log.Println("Invalid token:", err)
			return "", err
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println("Error parsing token:", err)
		return "", err
	}
	if !token.Valid {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		log.Println("Invalid token:", err)
		return "", errors.New("invalid token")
	}
	return claims.UserFullName, nil
}

func updateCustomerStatus(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := authByCookie(w, r)
		if err != nil {
			return
		}
		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		var req userStatusUpdate
		err = json.Unmarshal(reqData, &req)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		customerId, err := strconv.Atoi(req.CustomerId)
		if err != nil {
			http.Error(w, "Invalid customerId", http.StatusBadRequest)
			return
		}
		err = updateCustomerStatusInDB(db, customerId, req.Status)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		w.Header().Set("Content-Type", "application/json")
	}
}

func updateCustomerStatusInDB(db *sqlx.DB, customerID int, newStatus string) error {
	const query = `
		UPDATE customer
		SET status = ?
		WHERE customer_id = ?
	`
	_, err := db.Exec(query, newStatus, customerID)
	if err != nil {
		return err
	}
	return nil
}

func CheckUser(db *sqlx.DB, email, password string) (user, error) {
	const query = `
	SELECT
	    user_full_name,
		email,
		password
	FROM
		user
	WHERE
		email = ? AND password = ?
`
	var u user
	err := db.Get(&u, query, email, password)
	if err != nil {
		return user{}, err
	}
	return u, nil
}

func selectCustomers(db *sqlx.DB, userFullName string) ([]*customer, error) {
	const query = `
        SELECT
            customer_id,
            account_number,
            last_name,
            first_name,
            middle_name,
            date_of_birth,
            tax_id,
            status
        FROM
            customer
        WHERE user_full_name = ?
    `
	var customers []*customer
	err := db.Select(&customers, query, userFullName)
	if err != nil {
		return nil, err
	}
	return customers, nil
}
