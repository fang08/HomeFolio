package house

import (
	"database/sql"
	"model"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type (
	// Handler : Handler for house related requests.
	Handler struct {
		db    *sql.DB
		dbURL string
	}
)

// NewHandler : Create a Handler instance.
func NewHandler(session *sql.DB, dburl string) (h *Handler) {
	return &Handler{session, dburl}
}

// FetchSingleHouseInfo : Return info for a specific house.
//				   URL: "/api/v1/houseInfo/:hid"
//				   Method: GET
//				   Return 200 OK on success.
//				   Return 404 Not Found if house does not exist.
func (h *Handler) FetchSingleHouseInfo(c echo.Context) (err error) {
	// Retrieve house info from database
	houseD := new(model.House)
	err = h.db.QueryRow("SELECT * FROM house WHERE h_id = &var1", c.Param("hid")).Scan(&houseD.HID, &houseD.UID,
		&houseD.BathroomCnt, &houseD.BedroomCnt, &houseD.BuildingQualityID, &houseD.LivingAreaSize, &houseD.Latitude,
		&houseD.Longitude, &houseD.LotSize, &houseD.CityID, &houseD.County, &houseD.Zip, &houseD.YearBuilt,
		&houseD.StoryNum, &houseD.Price, &houseD.Tax)
	if err != nil {
		if err == sql.ErrNoRows {
			return &echo.HTTPError{Code: http.StatusNotFound, Message: "House does not exist."}
		}
		return err
	}

	return c.JSON(http.StatusOK, houseD)
}

// FetchRegionHouseInfo : Return information about houses located in the queried area.
//				   URL: "/api/v1/houseInfo"
//				   Method: GET
//				   Return 200 OK on success.
func (h *Handler) FetchRegionHouseInfo(c echo.Context) (err error) {
	// Retrieve house info from database
	neLat, _ := strconv.ParseFloat(c.QueryParam("ne_lat"), 64)
	swLat, _ := strconv.ParseFloat(c.QueryParam("sw_lat"), 64)
	neLng, _ := strconv.ParseFloat(c.QueryParam("ne_lng"), 64)
	swLng, _ := strconv.ParseFloat(c.QueryParam("sw_lng"), 64)
	count, _ := strconv.Atoi(c.QueryParam("count"))

	query :=
		`SELECT *
		FROM (
			SELECT h_id, latitude, longitude
			FROM house 
			WHERE latitude < &var1 and latitude > &var2 and longitude < &var3 and longitude > &var4 
			ORDER BY DBMS_RANDOM.VALUE) 
		WHERE ROWNUM <= &var5`

	rows, err := h.db.Query(query, neLat, swLat, neLng, swLng, count)
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
		houseC := new(HouseC)
		err = rows.Scan(&houseC.HID, &houseC.Latitude, &houseC.Longitude)
		if err != nil {
			return err
		}

		houses = append(houses, houseC)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, houses)
}

// FetchTopLikedHouses : Return information about top liked houses in an area.
//				   URL: "/api/v1/topliked"
//				   Method: GET
//				   Return 200 OK on success.
func (h *Handler) FetchTopLikedHouses(c echo.Context) (err error) {
	// Retrieve house info from database
	neLat, _ := strconv.ParseFloat(c.QueryParam("ne_lat"), 64)
	swLat, _ := strconv.ParseFloat(c.QueryParam("sw_lat"), 64)
	neLng, _ := strconv.ParseFloat(c.QueryParam("ne_lng"), 64)
	swLng, _ := strconv.ParseFloat(c.QueryParam("sw_lng"), 64)
	count, _ := strconv.Atoi(c.QueryParam("count"))

	query :=
		`SELECT *
		FROM
			(SELECT h_id, latitude, longitude, count(*) as num
			FROM
				(SELECT h_id, latitude, longitude
				FROM house 
				WHERE latitude < &var1 and latitude > &var2 and longitude < &var3 and longitude > &var4)
				NATURAL JOIN
				likes
			GROUP BY h_id, latitude, longitude
			ORDER BY num DESC)
		WHERE ROWNUM <= &var5`

	rows, err := h.db.Query(query, neLat, swLat, neLng, swLng, count)
	if err != nil {
		return err
	}
	defer rows.Close()

	type HouseC struct {
		HID       string  `json:"h_id"`
		Latitude  float32 `json:"latitude"`
		Longitude float32 `json:"longitude"`
		Likes     int     `json:"likes"`
	}

	var houses []*HouseC

	for rows.Next() {
		houseC := new(HouseC)
		err = rows.Scan(&houseC.HID, &houseC.Latitude, &houseC.Longitude, &houseC.Likes)
		if err != nil {
			return err
		}

		houses = append(houses, houseC)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, houses)
}

// FetchTopViewedHouses : Return information about top viewed houses in an area.
//				   URL: "/api/v1/topviewed"
//				   Method: GET
//				   Return 200 OK on success.
func (h *Handler) FetchTopViewedHouses(c echo.Context) (err error) {
	// Retrieve house info from database
	neLat, _ := strconv.ParseFloat(c.QueryParam("ne_lat"), 64)
	swLat, _ := strconv.ParseFloat(c.QueryParam("sw_lat"), 64)
	neLng, _ := strconv.ParseFloat(c.QueryParam("ne_lng"), 64)
	swLng, _ := strconv.ParseFloat(c.QueryParam("sw_lng"), 64)
	count, _ := strconv.Atoi(c.QueryParam("count"))

	query :=
		`SELECT *
		FROM
			(SELECT h_id, latitude, longitude, count(*) as num
			FROM
				(SELECT h_id, latitude, longitude
				FROM house 
				WHERE latitude < &var1 and latitude > &var2 and longitude < &var3 and longitude > &var4)
				NATURAL JOIN
				viewed
			GROUP BY h_id, latitude, longitude
			ORDER BY num DESC)
		WHERE ROWNUM <= &var5`

	rows, err := h.db.Query(query, neLat, swLat, neLng, swLng, count)
	if err != nil {
		return err
	}
	defer rows.Close()

	type HouseC struct {
		HID       string  `json:"h_id"`
		Latitude  float32 `json:"latitude"`
		Longitude float32 `json:"longitude"`
		Views     int     `json:"views"`
	}

	var houses []*HouseC

	for rows.Next() {
		houseC := new(HouseC)
		err = rows.Scan(&houseC.HID, &houseC.Latitude, &houseC.Longitude, &houseC.Views)
		if err != nil {
			return err
		}

		houses = append(houses, houseC)
	}
	if err = rows.Err(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, houses)
}

// Shutdown : Gracefully shutdown user handler.
func (h *Handler) Shutdown() {
	h.db.Close()
}
