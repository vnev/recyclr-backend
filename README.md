## recyclr-backend

Backend Repository of Project for **CS407 Fall 2018**


### Members

Cory Laker, Geoffrey Myers, Pranav Vasudha, Ryan Walden, Vedant Nevetia, Zachary Rich

### API Routes
#### Public routes
- `/user`: Create new user
    - Type: POST
    - Request Body Parameters
        - `email`: User email
        - `name`: User name
        - `passwd`: User password
        - `is_company`: Is user a company? ('f' for user, 't' for company)
    
- `/company`: Create new company
    - Type: POST
        - Request Body Parameters
            - `email`: Company email
            - `name`: Company name
            - `passwd`: Company password
            - `is_company`: Is user a company? ('f' for user, 't' for company)

- `/signin`: Sign in user/company
    - Type: POST
    - Request Body Parameters
        - `email`: User/Company email
        - `passwd`: User/Company password
    - Returns `token` response containing a JSON Web Token

- `/user/{id}/delete`: **ADMIN ROUTE**: Deletes a user from database
    - Type: GET
    - GET request accepts a user ID in the URL and deletes that user and all listings, etc. associated with that user's ID if found

- `/user/{id}/ban`: **ADMIN ROUTE**: Bans a user from Recyclr
    - Type: GET
    - GET request accepts a user ID in the URL and bans the user from Recyclr. Authentication with the same URL will fail henceforth

#### Protected routes - requires signed in user/company
- `/user/{id}`: Update user
    - Type: PUT
    - Request Body Parameters
        - `email`: New user email
        - `passwd`: New user password
        - `address`: New user address
        - `name`: New user name
    - Can send in only the parameters which need to be updated, not all

- `/user/{id}`: Get user
    - Type: GET
    - Returns relevant user information based if user with `id` is found

- `/user/progress/{id}`: Get user's recycling progress
    - Type: GET
    - Returns user's recycling progress

- `/user/delete`: Delete user's profile **NOT IMPLEMENTED**
    - Type: POST
    - Request Body Parameters

- `/user/logout`: Logout the user
    - Type: POST
    - Request Body Parameters
        - `user_id`: ID of the user

- `/companies`: Get all companies registered in DB
    - Type: GET

- `/company/{id}`: Get company account information associated with ID
    - Type: GET
    - Returns company account associated with `id`

- `/company/delete`: Delete company's profile **NOT IMPLEMENTED**
    - Type: POST
    - Request Body Parameters

- `/listings`: Get all listings
    - Type: GET
    - Returns all listings in the database

- `/listing/{id}`: Get listing associated with ID
    - Type: GET
    - Returns listing information associated with `id`

- `/listing`: Create a new listing
    - Type: POST
    - Request Body Parameters
        - `title`: string:  Listing title
        - `description`: string: Listing description
        <!-- - `img_hash`: Image hash -->
        - `material_type`: string: Type of material being recycled (hopefully will be from a constant set of strings so we can filter by a particular type of material in the database)
        - `material_weight`: float: Weight of the material being recycled
        - `user_id`: integer: ID of the user creating the listing

- `/listing/{id}/update`: Update an existing listing
    - Type: POST
    - Request parameters
        - `id`: ID of the listing being updated
    - Request Body Parameters
        - `title`: New title
        - `description`: New description
        - `material_weight`: New material weight
        - `material_type`: New material type
        - `is_active`: 'f' if order was purchased, 't' if order is still active
        - `pickup_date_time`: Date and time for pickup in case order was placed

    - Returns a status code indicating success or failure

