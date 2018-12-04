
> recyclr-backend

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
* [func ConnectToDB(path string)](#ConnectToDB)


#### <a name="pkg-files">Package files</a>
[db.go](/src/target/db.go) 



## <a name="pkg-variables">Variables</a>
``` go
var DBconn *sql.DB
```
DBconn is the main database connection object, used globally.



## <a name="ConnectToDB">func</a> [ConnectToDB](/src/target/db.go?s=395:424#L17)
``` go
func ConnectToDB(path string)
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
* [func CreateInvoice(w http.ResponseWriter, r *http.Request)](#CreateInvoice)
* [func CreateListing(w http.ResponseWriter, r *http.Request)](#CreateListing)
* [func CreateOrder(w http.ResponseWriter, r *http.Request)](#CreateOrder)
* [func CreateTimeslot(w http.ResponseWriter, r *http.Request)](#CreateTimeslot)
* [func CreateUser(w http.ResponseWriter, r *http.Request)](#CreateUser)
* [func DeductUserPoints(w http.ResponseWriter, r *http.Request)](#DeductUserPoints)
* [func DeleteListing(w http.ResponseWriter, r *http.Request)](#DeleteListing)
* [func DeleteUser(w http.ResponseWriter, r *http.Request)](#DeleteUser)
* [func FreezeListing(w http.ResponseWriter, r *http.Request)](#FreezeListing)
* [func GetCompanies(w http.ResponseWriter, r *http.Request)](#GetCompanies)
* [func GetFrozenListings(w http.ResponseWriter, r *http.Request)](#GetFrozenListings)
* [func GetInvoices(w http.ResponseWriter, r *http.Request)](#GetInvoices)
* [func GetListing(w http.ResponseWriter, r *http.Request)](#GetListing)
* [func GetListings(w http.ResponseWriter, r *http.Request)](#GetListings)
* [func GetMessages(w http.ResponseWriter, r *http.Request)](#GetMessages)
* [func GetOrder(w http.ResponseWriter, r *http.Request)](#GetOrder)
* [func GetOrders(w http.ResponseWriter, r *http.Request)](#GetOrders)
* [func GetProgress(w http.ResponseWriter, r *http.Request)](#GetProgress)
* [func GetTimeslot(w http.ResponseWriter, r *http.Request)](#GetTimeslot)
* [func GetTimeslots(w http.ResponseWriter, r *http.Request)](#GetTimeslots)
* [func GetTransactions(w http.ResponseWriter, r *http.Request)](#GetTransactions)
* [func GetUser(w http.ResponseWriter, r *http.Request)](#GetUser)
* [func LogoutUser(w http.ResponseWriter, r *http.Request)](#LogoutUser)
* [func PutMessage(w http.ResponseWriter, r *http.Request)](#PutMessage)
* [func StripePayment(w http.ResponseWriter, r *http.Request)](#StripePayment)
* [func UnfreezeListing(w http.ResponseWriter, r *http.Request)](#UnfreezeListing)
* [func UpdateListing(w http.ResponseWriter, r *http.Request)](#UpdateListing)
* [func UpdateOrder(w http.ResponseWriter, r *http.Request)](#UpdateOrder)
* [func UpdateRating(w http.ResponseWriter, r *http.Request)](#UpdateRating)
* [func UpdateTimeslot(w http.ResponseWriter, r *http.Request)](#UpdateTimeslot)
* [func UpdateUser(w http.ResponseWriter, r *http.Request)](#UpdateUser)
* [type Invoice](#Invoice)
* [type Listing](#Listing)
* [type Message](#Message)
* [type Order](#Order)
* [type Timeslot](#Timeslot)
* [type User](#User)


#### <a name="pkg-files">Package files</a>
[Authentication.go](/src/target/Authentication.go) [Company.go](/src/target/Company.go) [Invoice.go](/src/target/Invoice.go) [Listing.go](/src/target/Listing.go) [Message.go](/src/target/Message.go) [Order.go](/src/target/Order.go) [Stripe.go](/src/target/Stripe.go) [Timeslot.go](/src/target/Timeslot.go) [User.go](/src/target/User.go) 





## <a name="AuthMiddleware">func</a> [AuthMiddleware](/src/target/Authentication.go?s=318:377#L15)
``` go
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc
```
AuthMiddleware is the authentication middleware, ensuring that users are authenticated
before accessing protected routes.



## <a name="AuthenticateUser">func</a> [AuthenticateUser](/src/target/User.go?s=8712:8773#L278)
``` go
func AuthenticateUser(w http.ResponseWriter, r *http.Request)
```
AuthenticateUser generates a JWT for the user and returns it in JSON format.



## <a name="BanUser">func</a> [BanUser](/src/target/User.go?s=7968:8020#L250)
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



## <a name="CreateInvoice">func</a> [CreateInvoice](/src/target/Invoice.go?s=694:752#L27)
``` go
func CreateInvoice(w http.ResponseWriter, r *http.Request)
```
CreateInvoice creates a new invoice and stores it into the database
requires: Listing ID passed into request body



## <a name="CreateListing">func</a> [CreateListing](/src/target/Listing.go?s=5983:6041#L174)
``` go
func CreateListing(w http.ResponseWriter, r *http.Request)
```
CreateListing creates a new listing in the database. It expects title, description, img_hash,
material_type, material_weight, user_id, and address. It also reads the AWS configuration to store images.



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



## <a name="CreateUser">func</a> [CreateUser](/src/target/User.go?s=1818:1873#L59)
``` go
func CreateUser(w http.ResponseWriter, r *http.Request)
```
CreateUser creates a new user in the database. It expects address, email, user_name, is_company, and passwd fields.



## <a name="DeductUserPoints">func</a> [DeductUserPoints](/src/target/User.go?s=14577:14638#L463)
``` go
func DeductUserPoints(w http.ResponseWriter, r *http.Request)
```
DeductUserPoints deducts user points and applies new price to listing



## <a name="DeleteListing">func</a> [DeleteListing](/src/target/Listing.go?s=13377:13435#L426)
``` go
func DeleteListing(w http.ResponseWriter, r *http.Request)
```
DeleteListing deletes a listing from the database given its' listing_id. It will only work if
the user sending the request has sufficient admin priveliges.



## <a name="DeleteUser">func</a> [DeleteUser](/src/target/User.go?s=11332:11387#L369)
``` go
func DeleteUser(w http.ResponseWriter, r *http.Request)
```
DeleteUser deletes a listing from the database given their user_id. It will only work if
the user sending the request has sufficient admin priveliges.



## <a name="FreezeListing">func</a> [FreezeListing](/src/target/Listing.go?s=9384:9442#L290)
``` go
func FreezeListing(w http.ResponseWriter, r *http.Request)
```
FreezeListing freezes a listing for a particular company



## <a name="GetCompanies">func</a> [GetCompanies](/src/target/Company.go?s=220:277#L13)
``` go
func GetCompanies(w http.ResponseWriter, r *http.Request)
```
GetCompanies returns all companies from the database in the JSON format. It does not require
any parameters.



## <a name="GetFrozenListings">func</a> [GetFrozenListings](/src/target/Listing.go?s=2279:2341#L70)
``` go
func GetFrozenListings(w http.ResponseWriter, r *http.Request)
```
GetFrozenListings gets all frozen listings for a particular user.



## <a name="GetInvoices">func</a> [GetInvoices](/src/target/Invoice.go?s=2792:2848#L87)
``` go
func GetInvoices(w http.ResponseWriter, r *http.Request)
```
GetInvoices returns the status and listing ID associated with
the invoice identified by invoice_id (passed into request body)



## <a name="GetListing">func</a> [GetListing](/src/target/Listing.go?s=1204:1259#L43)
``` go
func GetListing(w http.ResponseWriter, r *http.Request)
```
GetListing returns a listing from the database in JSON format, given the specific listing_id as a URL parameter.



## <a name="GetListings">func</a> [GetListings](/src/target/Listing.go?s=4567:4623#L137)
``` go
func GetListings(w http.ResponseWriter, r *http.Request)
```
GetListings returns all listings from the database in JSON format.



## <a name="GetMessages">func</a> [GetMessages](/src/target/Message.go?s=506:562#L23)
``` go
func GetMessages(w http.ResponseWriter, r *http.Request)
```
GetMessages returns all messages between user and company for a particular listing



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



## <a name="GetProgress">func</a> [GetProgress](/src/target/User.go?s=12193:12249#L398)
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



## <a name="GetTransactions">func</a> [GetTransactions](/src/target/User.go?s=13504:13564#L431)
``` go
func GetTransactions(w http.ResponseWriter, r *http.Request)
```
GetTransactions returns all orders for a company in JSON format, given their user_id as a URL parameter.



## <a name="GetUser">func</a> [GetUser](/src/target/User.go?s=892:944#L35)
``` go
func GetUser(w http.ResponseWriter, r *http.Request)
```
GetUser returns a user from the database in JSON format, given the specific user_id as a URL parameter.
TODO: maybe just return a newly defined struct without password field.



## <a name="LogoutUser">func</a> [LogoutUser](/src/target/User.go?s=10527:10582#L340)
``` go
func LogoutUser(w http.ResponseWriter, r *http.Request)
```
LogoutUser logs a user out, setting their JWT to 0. It expects the user_id to be sent in the request body.



## <a name="PutMessage">func</a> [PutMessage](/src/target/Message.go?s=1705:1760#L63)
``` go
func PutMessage(w http.ResponseWriter, r *http.Request)
```
PutMessage adds a new message between user and company for particular listing



## <a name="StripePayment">func</a> [StripePayment](/src/target/Stripe.go?s=331:389#L18)
``` go
func StripePayment(w http.ResponseWriter, r *http.Request)
```
StripePayment handles a payment from Stripe, given the payment token in a form in the request body.



## <a name="UnfreezeListing">func</a> [UnfreezeListing](/src/target/Listing.go?s=10566:10626#L333)
``` go
func UnfreezeListing(w http.ResponseWriter, r *http.Request)
```
UnfreezeListing unfreeze a particular listing



## <a name="UpdateListing">func</a> [UpdateListing](/src/target/Listing.go?s=11635:11693#L369)
``` go
func UpdateListing(w http.ResponseWriter, r *http.Request)
```
UpdateListing updates a listing in the database, given its' listing_id and other fields requesting to be changed.



## <a name="UpdateOrder">func</a> [UpdateOrder](/src/target/Order.go?s=3530:3586#L113)
``` go
func UpdateOrder(w http.ResponseWriter, r *http.Request)
```
UpdateOrder updates an order in the database, given its' order_id and other fields requesting to be changed.



## <a name="UpdateRating">func</a> [UpdateRating](/src/target/User.go?s=15841:15898#L512)
``` go
func UpdateRating(w http.ResponseWriter, r *http.Request)
```
UpdateRating updates a user's rating in the database



## <a name="UpdateTimeslot">func</a> [UpdateTimeslot](/src/target/Timeslot.go?s=3631:3690#L113)
``` go
func UpdateTimeslot(w http.ResponseWriter, r *http.Request)
```
UpdateTimeslot updates a timeslot in the database, given its' time_id and other fields requesting to be changed.



## <a name="UpdateUser">func</a> [UpdateUser](/src/target/User.go?s=2751:2806#L81)
``` go
func UpdateUser(w http.ResponseWriter, r *http.Request)
```
UpdateUser updates a user in the database, given its' user_id and other fields requesting to be changed.




## <a name="Invoice">type</a> [Invoice](/src/target/Invoice.go?s=204:574#L14)
``` go
type Invoice struct {
    ID              int     `json:"invoice_id"`
    Status          bool    `json:"invoice_status"`
    Price           float64 `json:"price"`
    CreatedAt       string  `json:"created_at"`
    CompanyName     string  `json:"company_name"`
    InvoiceDateTime string  `json:"invoice_date_time"`
    UserName        string  `json:"user_name"`
    ForListing      Listing
}

```
Invoice struct to hold information pertaining to an invoice










## <a name="Listing">type</a> [Listing](/src/target/Listing.go?s=452:1086#L25)
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
    Address        string  `json:"address"`
    FrozenBy       int     `json:"frozen_by"`
    Price          float64 `json:"price"`
    Username       string  `json:"username"`
    CompanyName    string  `json:"company_name"`
}

```
Listing struct contains the listing schema in a struct format.










## <a name="Message">type</a> [Message](/src/target/Message.go?s=163:418#L13)
``` go
type Message struct {
    ID        int    `json:"message_id"`
    Timestamp string `json:"message_time"`
    FromUser  int    `json:"from_user"`
    ToUser    int    `json:"to_user"`
    ListingID int    `json:"for_listing"`
    Content   string `json:"message_content"`
}

```
Message struct to store message information










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










## <a name="User">type</a> [User](/src/target/User.go?s=267:709#L18)
``` go
type User struct {
    ID        int     `json:"user_id"`
    Address   string  `json:"address"`
    Email     string  `json:"email"`
    Name      string  `json:"name"`
    IsCompany bool    `json:"is_company"`
    Rating    float32 `json:"rating"`
    JoinedOn  string  `json:"joined_on"`
    Password  string  `json:"passwd"`
    Token     string  `json:"token"`
    City      string  `json:"city"`
    State     string  `json:"state"`
    Points    int     `json:"points"`
}

```
User struct contains the user schema in a struct format.














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
