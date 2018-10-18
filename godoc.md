

> recyclr-backend





- - -


# config
`import "/Users/zachrich/go/src/github.com/vnev/recyclr-backend/config"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package config reads our configuration for our database and AWS.




## <a name="pkg-index">Index</a>
* [func LoadAWSConfiguration(file string) (*credentials.Credentials, error)](#LoadAWSConfiguration)
* [type AWSCredentials](#AWSCredentials)
* [type Config](#Config)
  * [func LoadConfiguration(file string) Config](#LoadConfiguration)


#### <a name="pkg-files">Package files</a>
[config.go](/src/target/config.go) 





## <a name="LoadAWSConfiguration">func</a> [LoadAWSConfiguration](/src/target/config.go?s=1161:1233#L42)
``` go
func LoadAWSConfiguration(file string) (*credentials.Credentials, error)
```
LoadAWSConfiguration loads AWS Configuration and returns the credentials.




## <a name="AWSCredentials">type</a> [AWSCredentials](/src/target/config.go?s=608:732#L23)
``` go
type AWSCredentials struct {
    AWSAccessKeyID  string `json:"accessKeyId"`
    SecretAccessKey string `json:"secretAccessKey"`
}

```
AWSCredentials is the struct which holds our access information for AWS










## <a name="Config">type</a> [Config](/src/target/config.go?s=317:531#L14)
``` go
type Config struct {
    DBHost       string `json:"dbhost"`
    DBPass       string `json:"dbpass"`
    DBUser       string `json:"dbuser"`
    DBName       string `json:"dbname"`
    StripeSecret string `json:"stripe_secret"`
}

```
Config is the struct which will hold our various configuration settings such as
the database username and password, and our Stripe secret.







### <a name="LoadConfiguration">func</a> [LoadConfiguration](/src/target/config.go?s=824:866#L29)
``` go
func LoadConfiguration(file string) Config
```
LoadConfiguration loads the info in from a JSON file and returns the resulting struct.









- - -


# db
`import "/Users/zachrich/go/src/github.com/vnev/recyclr-backend/db"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package db implements the connection to our PostgreSQL server.




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [func ConnectToDB()](#ConnectToDB)


#### <a name="pkg-files">Package files</a>
[db.go](/src/target/db.go) 



## <a name="pkg-variables">Variables</a>
``` go
var DBconn *sql.DB
```
DBconn is the main database connection object, used globally.



## <a name="ConnectToDB">func</a> [ConnectToDB](/src/target/db.go?s=395:413#L17)
``` go
func ConnectToDB()
```
ConnectToDB opens a connection to the database, and keeps it open while the server is running.








- - -


# handlers
`import "/Users/zachrich/go/src/github.com/vnev/recyclr-backend/handlers"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
Package handlers contains all of our handlers for our HTTP routes on our API.




## <a name="pkg-index">Index</a>
* [func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc](#AuthMiddleware)
* [func AuthenticateUser(w http.ResponseWriter, r *http.Request)](#AuthenticateUser)
* [func BanUser(w http.ResponseWriter, r *http.Request)](#BanUser)
* [func CreateCompany(w http.ResponseWriter, r *http.Request)](#CreateCompany)
* [func CreateListing(w http.ResponseWriter, r *http.Request)](#CreateListing)
* [func CreateOrder(w http.ResponseWriter, r *http.Request)](#CreateOrder)
* [func CreateTimeslot(w http.ResponseWriter, r *http.Request)](#CreateTimeslot)
* [func CreateUser(w http.ResponseWriter, r *http.Request)](#CreateUser)
* [func DeleteCompany(w http.ResponseWriter, r *http.Request)](#DeleteCompany)
* [func DeleteListing(w http.ResponseWriter, r *http.Request)](#DeleteListing)
* [func DeleteOrder(w http.ResponseWriter, r *http.Request)](#DeleteOrder)
* [func DeleteTimeslot(w http.ResponseWriter, r *http.Request)](#DeleteTimeslot)
* [func DeleteUser(w http.ResponseWriter, r *http.Request)](#DeleteUser)
* [func GetCompanies(w http.ResponseWriter, r *http.Request)](#GetCompanies)
* [func GetListing(w http.ResponseWriter, r *http.Request)](#GetListing)
* [func GetListings(w http.ResponseWriter, r *http.Request)](#GetListings)
* [func GetOrder(w http.ResponseWriter, r *http.Request)](#GetOrder)
* [func GetOrders(w http.ResponseWriter, r *http.Request)](#GetOrders)
* [func GetProgress(w http.ResponseWriter, r *http.Request)](#GetProgress)
* [func GetTimeslot(w http.ResponseWriter, r *http.Request)](#GetTimeslot)
* [func GetTimeslots(w http.ResponseWriter, r *http.Request)](#GetTimeslots)
* [func GetTransactions(w http.ResponseWriter, r *http.Request)](#GetTransactions)
* [func GetUser(w http.ResponseWriter, r *http.Request)](#GetUser)
* [func LogoutUser(w http.ResponseWriter, r *http.Request)](#LogoutUser)
* [func StripePayment(w http.ResponseWriter, r *http.Request)](#StripePayment)
* [func UpdateCompany(w http.ResponseWriter, r *http.Request)](#UpdateCompany)
* [func UpdateListing(w http.ResponseWriter, r *http.Request)](#UpdateListing)
* [func UpdateOrder(w http.ResponseWriter, r *http.Request)](#UpdateOrder)
* [func UpdateTimeslot(w http.ResponseWriter, r *http.Request)](#UpdateTimeslot)
* [func UpdateUser(w http.ResponseWriter, r *http.Request)](#UpdateUser)
* [type Listing](#Listing)
* [type Order](#Order)
* [type Timeslot](#Timeslot)
* [type User](#User)


#### <a name="pkg-files">Package files</a>
[Authentication.go](/src/target/Authentication.go) [Company.go](/src/target/Company.go) [Listing.go](/src/target/Listing.go) [Order.go](/src/target/Order.go) [Stripe.go](/src/target/Stripe.go) [Timeslot.go](/src/target/Timeslot.go) [User.go](/src/target/User.go) 





## <a name="AuthMiddleware">func</a> [AuthMiddleware](/src/target/Authentication.go?s=318:377#L15)
``` go
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc
```
AuthMiddleware is the authentication middleware, ensuring that users are authenticated
before accessing protected routes.



## <a name="AuthenticateUser">func</a> [AuthenticateUser](/src/target/User.go?s=8066:8127#L263)
``` go
func AuthenticateUser(w http.ResponseWriter, r *http.Request)
```
AuthenticateUser generates a JWT for the user and returns it in JSON format.



## <a name="BanUser">func</a> [BanUser](/src/target/User.go?s=7322:7374#L235)
``` go
func BanUser(w http.ResponseWriter, r *http.Request)
```
BanUser bans a specific user given their user_id as a URL parameter.



## <a name="CreateCompany">func</a> [CreateCompany](/src/target/Company.go?s=1145:1203#L42)
``` go
func CreateCompany(w http.ResponseWriter, r *http.Request)
```
CreateCompany creates a new company in the database, and returns the newly created company in JSON format.
In the request body, it expects an address, email, user_name, is_company, and password.



## <a name="CreateListing">func</a> [CreateListing](/src/target/Listing.go?s=3078:3136#L92)
``` go
func CreateListing(w http.ResponseWriter, r *http.Request)
```
CreateListing creates a new listing in the database. It expects title, description, img_hash,
material_type, material_weight, user_id, and zipcode. It also reads the AWS configuration to store images.



## <a name="CreateOrder">func</a> [CreateOrder](/src/target/Order.go?s=2492:2548#L86)
``` go
func CreateOrder(w http.ResponseWriter, r *http.Request)
```
CreateOrder creates a new listing in the database. It expects user_id, company_id, total, and confirmed.



## <a name="CreateTimeslot">func</a> [CreateTimeslot](/src/target/Timeslot.go?s=2565:2624#L86)
``` go
func CreateTimeslot(w http.ResponseWriter, r *http.Request)
```
CreateTimeslot creates a new timeslot in the database. It expects user_id, day, start_time, and end_time.



## <a name="CreateUser">func</a> [CreateUser](/src/target/User.go?s=1647:1702#L56)
``` go
func CreateUser(w http.ResponseWriter, r *http.Request)
```
CreateUser creates a new user in the database. It expects address, email, user_name, is_company, and passwd fields.



## <a name="DeleteCompany">func</a> [DeleteCompany](/src/target/Company.go?s=2376:2434#L69)
``` go
func DeleteCompany(w http.ResponseWriter, r *http.Request)
```
DeleteCompany deletes a company from the database. It expects the user_id, and will only work if
the user sending the request has sufficient admin priveliges.



## <a name="DeleteListing">func</a> [DeleteListing](/src/target/Listing.go?s=7529:7587#L236)
``` go
func DeleteListing(w http.ResponseWriter, r *http.Request)
```
DeleteListing deletes a listing from the database given its' listing_id. It will only work if
the user sending the request has sufficient admin priveliges.



## <a name="DeleteOrder">func</a> [DeleteOrder](/src/target/Order.go?s=5249:5305#L170)
``` go
func DeleteOrder(w http.ResponseWriter, r *http.Request)
```
DeleteOrder deletes an order from the database given its' order_id. It will only work if
the user sending the request has sufficient admin priveliges.



## <a name="DeleteTimeslot">func</a> [DeleteTimeslot](/src/target/Timeslot.go?s=5377:5436#L170)
``` go
func DeleteTimeslot(w http.ResponseWriter, r *http.Request)
```
DeleteTimeslot deletes a timeslot from the database given its' time_id. It will only work if
the user sending the request has sufficient admin priveliges.



## <a name="DeleteUser">func</a> [DeleteUser](/src/target/User.go?s=10686:10741#L354)
``` go
func DeleteUser(w http.ResponseWriter, r *http.Request)
```
DeleteUser deletes a listing from the database given their user_id. It will only work if
the user sending the request has sufficient admin priveliges.



## <a name="GetCompanies">func</a> [GetCompanies](/src/target/Company.go?s=220:277#L13)
``` go
func GetCompanies(w http.ResponseWriter, r *http.Request)
```
GetCompanies returns all companies from the database in the JSON format. It does not require
any parameters.



## <a name="GetListing">func</a> [GetListing](/src/target/Listing.go?s=1034:1089#L39)
``` go
func GetListing(w http.ResponseWriter, r *http.Request)
```
GetListing returns a listing from the database in JSON format, given the specific listing_id as a URL parameter.



## <a name="GetListings">func</a> [GetListings](/src/target/Listing.go?s=1954:2010#L62)
``` go
func GetListings(w http.ResponseWriter, r *http.Request)
```
GetListings returns all listings from the database in JSON format.



## <a name="GetOrder">func</a> [GetOrder](/src/target/Order.go?s=530:583#L25)
``` go
func GetOrder(w http.ResponseWriter, r *http.Request)
```
GetOrder returns an order from the database in JSON format, given the specific order_id as a URL parameter.



## <a name="GetOrders">func</a> [GetOrders](/src/target/Order.go?s=1383:1437#L49)
``` go
func GetOrders(w http.ResponseWriter, r *http.Request)
```
GetOrders returns all orders from the database in JSON format for a specific given user,
given their user_id as a URL parameter.



## <a name="GetProgress">func</a> [GetProgress](/src/target/User.go?s=11547:11603#L383)
``` go
func GetProgress(w http.ResponseWriter, r *http.Request)
```
GetProgress gets the progress of a user's listings, returning it in JSON format given their user_id as a URL parameter.



## <a name="GetTimeslot">func</a> [GetTimeslot](/src/target/Timeslot.go?s=586:642#L26)
``` go
func GetTimeslot(w http.ResponseWriter, r *http.Request)
```
GetTimeslot returns a timeslot from the database in JSON format, given the specific time_id as a URL parameter.
This probably isn't needed at all.



## <a name="GetTimeslots">func</a> [GetTimeslots](/src/target/Timeslot.go?s=1424:1481#L49)
``` go
func GetTimeslots(w http.ResponseWriter, r *http.Request)
```
GetTimeslots returns all timeslots from the database for a specific user or company in JSON format.



## <a name="GetTransactions">func</a> [GetTransactions](/src/target/User.go?s=12733:12793#L415)
``` go
func GetTransactions(w http.ResponseWriter, r *http.Request)
```
GetTransactions returns all orders for a company in JSON format, given their user_id as a URL parameter.



## <a name="GetUser">func</a> [GetUser](/src/target/User.go?s=781:833#L32)
``` go
func GetUser(w http.ResponseWriter, r *http.Request)
```
GetUser returns a user from the database in JSON format, given the specific user_id as a URL parameter.
TODO: maybe just return a newly defined struct without password field.



## <a name="LogoutUser">func</a> [LogoutUser](/src/target/User.go?s=9881:9936#L325)
``` go
func LogoutUser(w http.ResponseWriter, r *http.Request)
```
LogoutUser logs a user out, setting their JWT to 0. It expects the user_id to be sent in the request body.



## <a name="StripePayment">func</a> [StripePayment](/src/target/Stripe.go?s=281:339#L15)
``` go
func StripePayment(w http.ResponseWriter, r *http.Request)
```
StripePayment handles a payment from Stripe, given the payment token in a form in the request body.



## <a name="UpdateCompany">func</a> [UpdateCompany](/src/target/Company.go?s=2146:2204#L63)
``` go
func UpdateCompany(w http.ResponseWriter, r *http.Request)
```
UpdateCompany updates a company in the database, and will return the number of rows updated. It expects
the user_id, along with all the fields that are requesting to be changed with their new information.



## <a name="UpdateListing">func</a> [UpdateListing](/src/target/Listing.go?s=5787:5845#L179)
``` go
func UpdateListing(w http.ResponseWriter, r *http.Request)
```
UpdateListing updates a listing in the database, given its' listing_id and other fields requesting to be changed.



## <a name="UpdateOrder">func</a> [UpdateOrder](/src/target/Order.go?s=3530:3586#L113)
``` go
func UpdateOrder(w http.ResponseWriter, r *http.Request)
```
UpdateOrder updates an order in the database, given its' order_id and other fields requesting to be changed.



## <a name="UpdateTimeslot">func</a> [UpdateTimeslot](/src/target/Timeslot.go?s=3631:3690#L113)
``` go
func UpdateTimeslot(w http.ResponseWriter, r *http.Request)
```
UpdateTimeslot updates a timeslot in the database, given its' time_id and other fields requesting to be changed.



## <a name="UpdateUser">func</a> [UpdateUser](/src/target/User.go?s=2521:2576#L77)
``` go
func UpdateUser(w http.ResponseWriter, r *http.Request)
```
UpdateUser updates a user in the database, given its' user_id and other fields requesting to be changed.




## <a name="Listing">type</a> [Listing](/src/target/Listing.go?s=452:916#L25)
``` go
type Listing struct {
    ID             int     `json:"listing_id"`
    Title          string  `json:"title"`
    Description    string  `json:"description"`
    ImageHash      string  `json:"img_hash"`
    MaterialType   string  `json:"material_type"`
    MaterialWeight float64 `json:"material_weight"`
    UserID         int     `json:"user_id"`
    Active         bool    `json:"is_active"`
    PickupDateTime string  `json:"pickup_date_time"`
    Zipcode        int     `json:"zipcode"`
}

```
Listing struct contains the listing schema in a struct format.










## <a name="Order">type</a> [Order](/src/target/Order.go?s=227:417#L16)
``` go
type Order struct {
    ID        int  `json:"order_id"`
    UserID    int  `json:"user_id"`
    CompanyID int  `json:"company_id"`
    Total     int  `json:"total"`
    Confirmed bool `json:"confirmed"`
}

```
Order struct contains the order DB schema in a struct format










## <a name="Timeslot">type</a> [Timeslot](/src/target/Timeslot.go?s=232:431#L16)
``` go
type Timeslot struct {
    ID        int    `json:"time_id"`
    UserID    int    `json:"user_id"`
    Day       string `json:"day"`
    StartTime string `json:"start_time"`
    EndTime   string `json:"end_time"`
}

```
Timeslot struct contains the timeslot schema in a struct format.










## <a name="User">type</a> [User](/src/target/User.go?s=267:598#L18)
``` go
type User struct {
    ID        int    `json:"user_id"`
    Address   string `json:"address"`
    Email     string `json:"email"`
    Name      string `json:"name"`
    IsCompany bool   `json:"is_company"`
    Rating    int    `json:"rating"`
    JoinedOn  string `json:"joined_on"`
    Password  string `json:"passwd"`
    Token     string `json:"token"`
}

```
User struct contains the user schema in a struct format.














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
