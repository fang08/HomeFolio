package user

import (
	"database/sql"
	"model"
	"net/http"
	"time"

	"github.com/satori/go.uuid"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type (
	// Handler : Handler for user related requests.
	Handler struct {
		db            *sql.DB
		dbURL         string
		key           string
		notifOperator chan model.Notification
	}

	userContainer struct {
		UID      string `json:"u_id"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password,omitempty"`
		Age      int64  `json:"age,omitempty"`
		Area     string `json:"area,omitempty"`
		Bio      string `json:"bio,omitempty"`
		Token    string `json:"token,omitempty"`
		Seller   bool   `json:"seller"`
		Buyer    bool   `json:"buyer"`
	}
)

// NewHandler : Create a Handler instance.
func NewHandler(session *sql.DB, dburl string, key string, operator chan model.Notification) (h *Handler) {
	return &Handler{session, dburl, key, operator}
}

// Signup : Create an user account.
//			URL: "/api/v1/signup"
//			Method: POST
//			Return 201 Created on success, along with the user data.
//			Return 400 Bad Request if one of username, email, password is empty, or username, email already used.
func (h *Handler) Signup(c echo.Context) (err error) {
	// Bind
	userC := &userContainer{}
	if err = c.Bind(userC); err != nil {
		return err
	}

	// Validate
	if userC.Username == "" || userC.Email == "" || userC.Password == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Invalid username, email or password."}
	}

	userD := new(model.User)
	err = h.db.QueryRow("SELECT * FROM acc_user WHERE email = &var1 and password = &var2", userC.Email, userC.Password).Scan(&userD.UID,
		&userD.Email, &userD.Username, &userD.Password, &userD.Age, &userD.Area, &userD.Bio)
	if err == nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Username or email already used."}
	} else if err != sql.ErrNoRows {
		return err
	}

	// Save user
	newID, err := uuid.NewV1()
	userC.UID = newID.String()
	userC.fillD(userD)

	stmt, err := h.db.Prepare("INSERT INTO acc_user VALUES (&var1, &var2, &var3, &var4, &var5, &var6, &var7)")
	_, err = stmt.Exec(userD.UID, userD.Email, userD.Username, userD.Password, userD.Age, userD.Area, userD.Bio)
	if err != nil {
		return err
	}

	if userC.Seller {
		stmt, err := h.db.Prepare("INSERT INTO seller VALUES (&var1)")
		_, err = stmt.Exec(userD.UID)
		if err != nil {
			return err
		}
	}

	if userC.Buyer {
		stmt, err := h.db.Prepare("INSERT INTO buyer VALUES (&var1)")
		_, err = stmt.Exec(userD.UID)
		if err != nil {
			return err
		}
	}

	// Don't send password
	userC.Password = ""

	return c.JSON(http.StatusCreated, userC)
}

// Login : Login to an account associated with the email address and the password.
//		   URL: "/api/v1/login"
//		   Method: POST
//		   Return 200 OK on success, along with the user data, which now contains a token.
//		   Return 401 Unauthorized if an account associated with the email address and password is not found.
func (h *Handler) Login(c echo.Context) (err error) {
	// Bind
	userC := new(userContainer)
	if err = c.Bind(userC); err != nil {
		return err
	}

	// Find user
	userD := new(model.User)
	err = h.db.QueryRow("SELECT * FROM acc_user WHERE email = &var1 and password = &var2", userC.Email, userC.Password).Scan(&userD.UID,
		&userD.Email, &userD.Username, &userD.Password, &userD.Age, &userD.Area, &userD.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "Invalid email or password."}
		}
		return err
	}

	seller := true
	err = h.db.QueryRow("SELECT * FROM seller WHERE u_id = &var1", userD.UID).Scan(new(string))
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		seller = false
	}

	buyer := true
	err = h.db.QueryRow("SELECT * FROM buyer WHERE u_id = &var1", userD.UID).Scan(new(string))
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		buyer = false
	}

	// Pack container
	userC.fillC(userD, seller, buyer)

	// JWT

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = userD.UID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response
	userC.Token, err = token.SignedString([]byte(h.key))
	if err != nil {
		return err
	}

	// Don't send password
	userC.Password = ""

	return c.JSON(http.StatusOK, userC)
}

// FetchUserInfo : Return user info for a specific user.
//				   URL: "/api/v1/userInfo/:uid"
//				   Method: GET
//				   Return 200 OK on success.
//				   Return 404 Not Found if the user is not in the database.
func (h *Handler) FetchUserInfo(c echo.Context) (err error) {
	// Retrieve user info from database
	userD := new(model.User)
	err = h.db.QueryRow("SELECT * FROM acc_user WHERE u_id = &var1", c.Param("uid")).Scan(&userD.UID,
		&userD.Email, &userD.Username, &userD.Password, &userD.Age, &userD.Area, &userD.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			return &echo.HTTPError{Code: http.StatusNotFound, Message: "User does not exist."}
		}
		return err
	}

	seller := true
	err = h.db.QueryRow("SELECT * FROM seller WHERE u_id = &var1", userD.UID).Scan(new(string))
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		seller = false
	}

	buyer := true
	err = h.db.QueryRow("SELECT * FROM buyer WHERE u_id = &var1", userD.UID).Scan(new(string))
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		buyer = false
	}

	userC := new(userContainer)
	userC.fillC(userD, seller, buyer)

	// Don't send password
	userC.Password = ""

	return c.JSON(http.StatusOK, userC)
}

// UpdateUserInfo : Update a user's info, uid and email are not allowed to be modified.
//					# Username cannot be empty.
//					# Request body must contain all fields, even if they are empty or not modified.
//					URL: "/api/v1/updateUserInfo"
//					Method: POST
//					Return 200 OK on success, along with the user data.
//					Return 404 Not Found if the user is not in the database.
func (h *Handler) UpdateUserInfo(c echo.Context) (err error) {
	uid := uidFromToken(c)

	userC := new(userContainer)
	if err = c.Bind(userC); err != nil {
		return err
	}

	if userC.Username == "" {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Username must not be empty."}
	}

	userD := new(model.User)
	err = h.db.QueryRow("SELECT * FROM acc_user WHERE u_id = &var1", uid).Scan(&userD.UID,
		&userD.Email, &userD.Username, &userD.Password, &userD.Age, &userD.Area, &userD.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			return &echo.HTTPError{Code: http.StatusNotFound, Message: "User does not exist."}
		}
		return err
	}

	userC.fillD(userD)

	stmt, err := h.db.Prepare("UPDATE acc_user SET username = &var1, password = &var2, age = &var3, area = &var4, bio = &var5 WHERE u_id = &var6")
	_, err = stmt.Exec(userD.Username, userD.Password, userD.Age, userD.Area, userD.Bio, uid)
	if err != nil {
		return err
	}

	// Don't send password
	userC.Password = ""

	//return c.JSON(http.StatusOK, user)
	return c.JSON(http.StatusOK, userC)
}

// FetchSellHouse : Return houses owned by a specific user.
//				   URL: "/api/v1/userInfo/sellHouses"
//				   Method: GET
//				   Return 200 OK on success.
func (h *Handler) FetchSellHouse(c echo.Context) (err error) {
	rows, err := h.db.Query("SELECT h_id, latitude, longitude FROM house WHERE u_id = &var1", c.Param("uid"))
	if err != nil {
		return err
	}
	defer rows.Close()

	type HouseC struct {
		HID       string  `json:"h_id"`
		Latitude  float32 `json:"latitude"`
		Longitude float32 `json:"longitude"`
	}

	var houses []*HouseC

	for rows.Next() {
		house := new(HouseC)
		err = rows.Scan(&house.HID, &house.Latitude, &house.Longitude)
		if err != nil {
			return err
		}

		houses = append(houses, house)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, houses)
}

// FetchLikedHouse : Return houses liked by a specific user.
//				   URL: "/api/v1/userInfo/likedHouses"
//				   Method: GET
//				   Return 200 OK on success.
func (h *Handler) FetchLikedHouse(c echo.Context) (err error) {
	rows, err := h.db.Query("SELECT house.h_id, latitude, longitude FROM house, likes WHERE house.h_id = likes.h_id and likes.u_id = &var1", c.Param("uid"))
	if err != nil {
		return err
	}
	defer rows.Close()

	type HouseC struct {
		HID       string  `json:"h_id"`
		Latitude  float32 `json:"latitude"`
		Longitude float32 `json:"longitude"`
	}

	var houses []*HouseC

	for rows.Next() {
		house := new(HouseC)
		err = rows.Scan(&house.HID, &house.Latitude, &house.Longitude)
		if err != nil {
			return err
		}

		houses = append(houses, house)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, houses)
}

// FetchViewedHouse : Return houses viewed by a specific user.
//				   URL: "/api/v1/userInfo/viewedHouses"
//				   Method: GET
//				   Return 200 OK on success.
func (h *Handler) FetchViewedHouse(c echo.Context) (err error) {
	query :=
		`SELECT house.h_id, latitude, longitude, time
		FROM viewed, house 
		WHERE house.h_id = viewed.h_id and viewed.u_id = &var1`
	rows, err := h.db.Query(query, c.Param("uid"))
	if err != nil {
		return err
	}
	defer rows.Close()

	type HouseC struct {
		HID       string  `json:"h_id"`
		Latitude  float32 `json:"latitude"`
		Longitude float32 `json:"longitude"`
		Time      string  `json:"time"`
	}

	var houses []*HouseC

	for rows.Next() {
		house := new(HouseC)
		err = rows.Scan(&house.HID, &house.Latitude, &house.Longitude, &house.Time)
		if err != nil {
			return err
		}

		houses = append(houses, house)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, houses)
}

// FetchBoughtHouse : Return houses bought by a specific user.
//				   URL: "/api/v1/userInfo/boughtHouses/:uid"
//				   Method: GET
//				   Return 200 OK on success.
func (h *Handler) FetchBoughtHouse(c echo.Context) (err error) {
	query := "SELECT h_id, latitude, longitude FROM bought_house WHERE u_id = &var1"
	rows, err := h.db.Query(query, c.Param("uid"))
	if err != nil {
		return err
	}
	defer rows.Close()

	type HouseC struct {
		HID       string  `json:"h_id"`
		Latitude  float32 `json:"latitude"`
		Longitude float32 `json:"longitude"`
	}

	var houses []*HouseC

	for rows.Next() {
		house := new(HouseC)
		err = rows.Scan(&house.HID, &house.Latitude, &house.Longitude)
		if err != nil {
			return err
		}

		houses = append(houses, house)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, houses)
}

// IsPopularUser : Return whether a user is popular.
//				   URL: "/api/v1/userInfo/isPopular/:uid"
//				   Method: GET
//				   Return 200 OK on success.
func (h *Handler) IsPopularUser(c echo.Context) (err error) {
	type response struct {
		Popular bool `json:"popular"`
	}

	likeCount := 0
	viewCount := 0
	query :=
		`SELECT *
		FROM
			(SELECT count(*) as likecount
			FROM
				((SELECT h_id
				FROM house
				WHERE u_id = &var1)
				NATURAL JOIN
				likes)),
			(SELECT count(*) as viewcount
			FROM
				((SELECT h_id
				FROM house
				WHERE u_id = &var1)
				NATURAL JOIN
				viewed))`
	err = h.db.QueryRow(query, c.Param("uid"), c.Param("uid")).Scan(&likeCount, &viewCount)
	if err != nil {
		return err
	}

	res := new(response)
	if likeCount > 5 || viewCount > 20 {
		res.Popular = true
	} else {
		res.Popular = false
	}

	return c.JSON(http.StatusOK, res)
}

// IsActiveUser : Return whether a user is active.
//				   URL: "/api/v1/userInfo/isActive/:uid"
//				   Method: GET
//				   Return 200 OK on success.
func (h *Handler) IsActiveUser(c echo.Context) (err error) {
	type response struct {
		Active bool `json:"active"`
	}

	viewCount := 0
	query :=
		`SELECT count(*)
		FROM viewed
		WHERE u_id = &var1 and time > to_date(&var2,'YYYY-MM-DD')`
	err = h.db.QueryRow(query, c.Param("uid"), time.Now().AddDate(-1, 0, 0).String()[:10]).Scan(&viewCount)
	if err != nil {
		return err
	}

	res := new(response)
	if viewCount > 5 {
		res.Active = true
	} else {
		res.Active = false
	}

	return c.JSON(http.StatusOK, res)
}

func uidFromToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return claims["uid"].(string)
}

func (userC *userContainer) fillC(userD *model.User, seller bool, buyer bool) {
	userC.UID = userD.UID
	userC.Email = userD.Email
	userC.Username = userD.Username
	userC.Password = userD.Password
	userC.Seller = seller
	userC.Buyer = buyer

	if userD.Age.Valid {
		userC.Age = userD.Age.Int64
	}
	if userD.Area.Valid {
		userC.Area = userD.Area.String
	}
	if userD.Bio.Valid {
		userC.Bio = userD.Bio.String
	}
}

func (userC *userContainer) fillD(userD *model.User) {
	userD.UID = userC.UID
	userD.Email = userC.Email
	userD.Username = userC.Username
	userD.Password = userC.Password

	if userC.Age == 0 {
		userD.Age.Valid = false
	} else {
		userD.Age.Int64 = userC.Age
		userD.Age.Valid = true
	}
	if userC.Area == "" {
		userD.Area.Valid = false
	} else {
		userD.Area.String = userC.Area
		userD.Area.Valid = true
	}
	if userC.Bio == "" {
		userD.Bio.Valid = false
	} else {
		userD.Bio.String = userC.Bio
		userD.Bio.Valid = true
	}
}

// Shutdown : Gracefully shutdown user handler.
func (h *Handler) Shutdown() {
	h.db.Close()
}
