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
		Picture  string `json:"picture,omitempty"`
		Age      int64  `json:"age,omitempty"`
		Area     string `json:"area,omitempty"`
		Bio      string `json:"bio,omitempty"`
		Token    string `json:"token,omitempty"`
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
	err = h.db.QueryRow("SELECT * FROM acc_user WHERE email = :var1 and password = :var2", userC.Email, userC.Password).Scan(&userD.UID,
		&userD.Email, &userD.Username, &userD.Password, &userD.Picture, &userD.Age, &userD.Area, &userD.Bio)
	if err == nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Username or email already used."}
	} else if err != sql.ErrNoRows {
		return err
	}

	// Save user
	newID, err := uuid.NewV1()
	userC.UID = newID.String()
	userC.fillD(userD)

	stmt, err := h.db.Prepare("INSERT INTO acc_user VALUES (:var1, :var2, :var3, :var4, :var5, :var6, :var7, :var8)")
	_, err = stmt.Exec(userD.UID, userD.Email, userD.Username, userD.Password, userD.Picture, userD.Age, userD.Area, userD.Bio)
	if err != nil {
		return err
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
	err = h.db.QueryRow("SELECT * FROM acc_user WHERE email = :var1 and password = :var2", userC.Email, userC.Password).Scan(&userD.UID,
		&userD.Email, &userD.Username, &userD.Password, &userD.Picture, &userD.Age, &userD.Area, &userD.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "Invalid email or password."}
		}
		return err
	}

	// Pack container
	userC.fillC(userD)

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
	err = h.db.QueryRow("SELECT * FROM acc_user WHERE u_id = :var1", c.Param("username")).Scan(&userD.UID,
		&userD.Email, &userD.Username, &userD.Password, &userD.Picture, &userD.Age, &userD.Area, &userD.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			return &echo.HTTPError{Code: http.StatusNotFound, Message: "User does not exist."}
		}
		return err
	}

	userC := new(userContainer)
	userC.fillC(userD)

	// Don't send password
	userC.Password = ""

	return c.JSON(http.StatusOK, userC)
}

// UpdateUserInfo : Update a user's info, uid and email are not allowed to be modified.
//					# Username cannot be empty.
//					# Image is to be encoded into base 64 string, and cannot be larger than 10 MB.
//					# Request body must contain all fields, even if they are empty or not modified.
//					URL: "/api/v1/updateUserInfo"
//					Method: POST
//					Return 200 OK on success, along with the user data.
//					Return 400 Bad Request if username is empty or the image is larger than 10 MB.
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

	// Do not accept image larger than 10 MB
	if len(userC.Picture) > 10485760 {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Image must be smaller than 10 MB."}
	}

	userD := new(model.User)
	err = h.db.QueryRow("SELECT * FROM acc_user WHERE u_id = :var1", uid).Scan(&userD.UID,
		&userD.Email, &userD.Username, &userD.Password, &userD.Picture, &userD.Age, &userD.Area, &userD.Bio)
	if err != nil {
		if err == sql.ErrNoRows {
			return &echo.HTTPError{Code: http.StatusNotFound, Message: "User does not exist."}
		}
		return err
	}

	userC.fillD(userD)

	stmt, err := h.db.Prepare("UPDATE acc_user SET username = :var1, password = :var2, pic = :var3, age = :var4, area = :var5, biography = :var6 WHERE u_id = :var7")
	_, err = stmt.Exec(userD.Username, userD.Password, userD.Picture, userD.Age, userD.Area, userD.Bio, userD.UID)
	if err != nil {
		return err
	}

	// Don't send password
	userC.Password = ""

	//return c.JSON(http.StatusOK, user)
	return c.JSON(http.StatusOK, userC)
}

func uidFromToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	return claims["uid"].(string)
}

func (userC *userContainer) fillC(userD *model.User) {
	userC.UID = userD.UID
	userC.Email = userD.Email
	userC.Username = userD.Username
	userC.Password = userD.Password
	if userD.Picture.Valid {
		userC.Picture = userD.Picture.String
	}
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
	if userC.Picture == "" {
		userD.Picture.Valid = false
	} else {
		userD.Picture.String = userC.Picture
		userD.Picture.Valid = true
	}
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