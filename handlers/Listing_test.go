package handlers

import (
	"testing"

	"github.com/vnev/recyclr-backend/db"
)

func TestGetListings(t *testing.T) {
	var listings []Listing
	//w.Header().Set("Content-Type", "application/json")
	sqlStatement := `SELECT u.user_name, l.user_id, l.listing_id, l.title, l.description, 
	l.material_type, l.material_weight, l.address, l.img_hash, l.pickup_date_time FROM Listings l 
	INNER JOIN Users u ON l.user_id=u.user_id 
	WHERE l.active='t'`
	rows, err := db.DBconn.Query(sqlStatement)
	if err != nil {
		t.Error(err.Error())
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
		t.Error(err.Error())
	}

}
