package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/satori/go.uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/gorilla/mux"
	c "github.com/vnev/recyclr-backend/config"
	"github.com/vnev/recyclr-backend/db"
)

// Listing struct contains the listing schema in a struct format.
type Listing struct {
	ID             int     `json:"listing_id"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	ImageHash      string  `json:"img_hash"`
	MaterialType   string  `json:"material_type"`
	MaterialWeight float64 `json:"material_weight"`
	UserID         int     `json:"user_id"`
	CompanyRating  float32 `json:"company_rating"`
	Active         bool    `json:"is_active"`
	PickupDateTime string  `json:"pickup_date_time"`
	Address        string  `json:"address"`
	FrozenBy       int     `json:"frozen_by"`
	Price          float64 `json:"price"`
	Username       string  `json:"username"`
	CompanyName    string  `json:"company_name"`
}

// GetListing returns a listing from the database in JSON format, given the specific listing_id as a URL parameter.
func GetListing(w http.ResponseWriter, r *http.Request) {
	var listing Listing
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r) // Get route params
	listingID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	//fmt.Printf("id route param is %d\n", userID)
	sqlStatement := `SELECT title, user_id, description, img_hash, material_type, 
	material_weight, address, active, frozen_by FROM listings WHERE listing_id=$1`
	err = db.DBconn.QueryRow(sqlStatement, listingID).Scan(&listing.Title, &listing.UserID,
		&listing.Description, &listing.ImageHash, &listing.MaterialType, &listing.MaterialWeight,
		&listing.Address, &listing.Active, &listing.FrozenBy)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	listing.ImageHash = "https://s3.us-east-2.amazonaws.com/recyclr/images/" + listing.ImageHash
	json.NewEncoder(w).Encode(&listing)
}

// GetFrozenListings gets all frozen listings for a particular user.
func GetFrozenListings(w http.ResponseWriter, r *http.Request) {
	var frozenListings []Listing
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	userID, err := strconv.Atoi(params["user_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type attributes struct {
		IsCompany bool `json:"is_company"`
	}

	var attr attributes
	if err = json.NewDecoder(r.Body).Decode(&attr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sqlStatement := ""
	if attr.IsCompany {
		sqlStatement = `SELECT u.user_name, u2.user_name, l.listing_id, l.user_id, u2.rating, l.title, 
		l.description, l.material_type, l.material_weight, l.address, l.frozen_by, l.img_hash, l.pickup_date_time 
		FROM Listings l 
		INNER JOIN Users u ON l.user_id=u.user_id 
		INNER JOIN Users u2 ON u2.user_id=l.frozen_by 
		WHERE l.active='f' and l.frozen_by=$1`
	} else {
		sqlStatement = `SELECT u.user_name, u2.user_name, l.listing_id, l.user_id, u2.rating, l.title, 
		l.description, l.material_type, l.material_weight, l.address, l.frozen_by, l.img_hash, l.pickup_date_time 
		FROM Listings l
		INNER JOIN Users u ON u.user_id=l.user_id
		INNER JOIN Users u2 ON u2.user_id=l.frozen_by
		WHERE l.active='f' and l.user_id=$1`
	}
	rows, err := db.DBconn.Query(sqlStatement, userID)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Check your request parameters", http.StatusBadRequest)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var listing Listing
		err = rows.Scan(&listing.Username, &listing.CompanyName, &listing.ID, &listing.UserID,
			&listing.CompanyRating, &listing.Title, &listing.Description, &listing.MaterialType,
			&listing.MaterialWeight, &listing.Address, &listing.FrozenBy, &listing.ImageHash,
			&listing.PickupDateTime)
		// TODO: Error check
		listing.ImageHash = "https://s3.us-east-2.amazonaws.com/recyclr/images/" + listing.ImageHash
		//fmt.Printf("ID is %d, Type is %s\n", listing.ID, listing.MaterialType)
		frozenListings = append(frozenListings, listing)
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(frozenListings)
}

// GetListings returns all listings from the database in JSON format.
func GetListings(w http.ResponseWriter, r *http.Request) {
	var listings []Listing
	w.Header().Set("Content-Type", "application/json")
	sqlStatement := `SELECT u.user_name, l.user_id, l.listing_id, l.title, l.description, 
	l.material_type, l.material_weight, l.address, l.img_hash, l.pickup_date_time FROM Listings l 
	INNER JOIN Users u ON l.user_id=u.user_id 
	WHERE l.active='t'`
	rows, err := db.DBconn.Query(sqlStatement)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var listing Listing
		err = rows.Scan(&listing.Username, &listing.UserID, &listing.ID, &listing.Title,
			&listing.Description, &listing.MaterialType, &listing.MaterialWeight,
			&listing.Address, &listing.ImageHash, &listing.PickupDateTime)
		//fmt.Printf("ID is %d, Type is %s\n", listing.ID, listing.MaterialType)
		listing.ImageHash = "https://s3.us-east-2.amazonaws.com/recyclr/images/" + listing.ImageHash
		listings = append(listings, listing)
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(listings)
}

// CreateListing creates a new listing in the database. It expects title, description, img_hash,
// material_type, material_weight, user_id, and address. It also reads the AWS configuration to store images.
func CreateListing(w http.ResponseWriter, r *http.Request) {
	var listing Listing

	awsCreds, err := c.LoadAWSConfiguration("config.json")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error fetching aws credentials", http.StatusInternalServerError)
		return
	}

	err = r.ParseForm()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	cfg := aws.NewConfig().WithRegion("us-east-2").WithCredentials(awsCreds)
	if err != nil {
		fmt.Println("failed")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	svc := s3.New(session.New(), cfg)

	file, h, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uniqueString := h.Filename + uuid.Must(uuid.NewV4()).String()
	extensionIndex := strings.LastIndex(h.Filename, ".")
	if extensionIndex < 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hash := sha256.New()
	hash.Write([]byte(uniqueString))
	sha := hash.Sum(nil)
	hashString := hex.EncodeToString(sha)
	hashedFilename := hashString + h.Filename[extensionIndex:]
	path := "/images/" + hashedFilename
	fmt.Printf("file path: %s\n", hashString)

	params := &s3.PutObjectInput{
		Bucket:        aws.String("recyclr"),
		Key:           aws.String(path),
		Body:          file,
		ContentLength: aws.Int64(h.Size),
		ContentType:   aws.String(h.Header.Get("Content-Type")),
	}

	_, err = svc.PutObject(params)
	if err != nil {
		http.Error(w, "Error uploading image!", http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}

	listing.Description = r.FormValue("description")
	listing.Title = r.FormValue("title")
	listing.MaterialType = r.FormValue("material_type")
	listing.Address = r.FormValue("address")
	listing.PickupDateTime = r.FormValue("pickup_date_time")

	materialWeight, _ := strconv.ParseFloat(r.FormValue("material_weight"), 64)
	listing.MaterialWeight = materialWeight
	listing.ImageHash = hashedFilename

	switch listing.MaterialType {
	case "Compost":
		listing.Price = 1.00
	case "Plastic":
		listing.Price = 1.50
	case "Electronics":
		listing.Price = 1.70
	case "Rubber":
		listing.Price = 1.90
	case "Textiles":
		listing.Price = 2.00
	default:
		listing.Price = 2.30
	}

	userID, _ := strconv.Atoi(r.FormValue("user_id"))
	listing.UserID = userID

	sqlStatement := `INSERT INTO listings (title, description, img_hash, material_type, material_weight, user_id, address, pickup_date_time, price)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
					RETURNING listing_id`

	err = db.DBconn.QueryRow(sqlStatement, listing.Title, listing.Description,
		listing.ImageHash, listing.MaterialType, listing.MaterialWeight, listing.UserID,
		listing.Address, listing.PickupDateTime, listing.Price).Scan(&listing.ID)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var listingCount float64
	sqlStatement = "SELECT COUNT(*) FROM listings WHERE user_id=$1"
	err = db.DBconn.QueryRow(sqlStatement, userID).Scan(&listingCount)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// fmt.Println("New listing created with ID: ", id)
	json.NewEncoder(w).Encode(listing)
}

// FreezeListing freezes a listing for a particular company
func FreezeListing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	listingID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type attributes struct {
		CompanyID    int    `json:"company_id"`
		MaterialType string `json:"material_type"`
	}

	var attr attributes
	_ = json.NewDecoder(r.Body).Decode(&attr)

	sqlStatement := "UPDATE Listings SET active='f', frozen_by=$1 WHERE listing_id=$2"
	row, err := db.DBconn.Exec(sqlStatement, attr.CompanyID, listingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	affectedCount, err := row.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resMap := make(map[string]string)
	resMap["message"] = "Success"
	resMap["rows affected"] = strconv.FormatInt(affectedCount, 10)
	res, err := json.Marshal(resMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// UnfreezeListing unfreeze a particular listing
func UnfreezeListing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	listingID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sqlStatement := `UPDATE Listings SET active='t', frozen_by=NULL, price=NULL WHERE listing_id=$1`
	row, err := db.DBconn.Exec(sqlStatement, listingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	affectedCount, err := row.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resMap := make(map[string]string)
	resMap["message"] = "Success"
	resMap["rows affected"] = strconv.FormatInt(affectedCount, 10)
	res, err := json.Marshal(resMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// UpdateListing updates a listing in the database, given its' listing_id and other fields requesting to be changed.
func UpdateListing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var listing Listing
	if err := json.NewDecoder(r.Body).Decode(&listing); err != nil {
		fmt.Println(err)
		http.Error(w, "Bad request parameters", http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	listingID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var values []interface{}
	j := 1
	sqlStatement := "UPDATE listings SET "

	structIterator := reflect.ValueOf(listing)
	for i := 0; i < structIterator.NumField(); i++ {
		field := structIterator.Type().Field(i).Name
		val := structIterator.Field(i).Interface()

		if !reflect.DeepEqual(val, reflect.Zero(structIterator.Field(i).Type()).Interface()) {
			sqlStatement += strings.ToLower(field) + "=$" + strconv.Itoa(j) + ", "
			j++
			values = append(values, val)
		}
	}

	sqlStatement = sqlStatement[:len(sqlStatement)-2]
	sqlStatement = sqlStatement + " WHERE listing_id" + "=$" + strconv.Itoa(j)
	values = append(values, listingID)
	row, err := db.DBconn.Exec(sqlStatement, values...)
	affectedCount, err := row.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resMap := make(map[string]string)
	resMap["message"] = "Success"
	resMap["rows affected"] = strconv.FormatInt(affectedCount, 10)
	res, err := json.Marshal(resMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// DeleteListing deletes a listing from the database given its' listing_id. It will only work if
// the user sending the request has sufficient admin priveliges.
func DeleteListing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	listingID, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sqlStatement := "DELETE FROM Listings WHERE listing_id=$1"
	row, err := db.DBconn.Exec(sqlStatement, listingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	affectedCount, err := row.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resMap := make(map[string]string)
	resMap["message"] = "Success"
	resMap["rows affected"] = strconv.FormatInt(affectedCount, 10)
	res, err := json.Marshal(resMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
