-- used for password salting + hashing god bless postgres
CREATE EXTENSION pgcrypto;

CREATE TABLE Users (
	user_id 		INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
	address 		TEXT NOT NULL,
	email 			TEXT NOT NULL,
	user_name 		TEXT NOT NULL,
	last_signed_in 	TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_DATE,
	is_company 		BOOLEAN NOT NULL,
	rating 			REAL NOT NULL DEFAULT 0,
	joined_on 		TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_DATE,
	passwd 			TEXT NOT NULL,
	token			TEXT NOT NULL DEFAULT '0'
);

CREATE TABLE Listings (
	title					TEXT NOT NULL,
	description				TEXT NOT NULL,
	img_hash				TEXT NOT NULL,
	listing_id 				INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
	material_type 			TEXT NOT NULL,
	material_weight 		REAL NOT NULL,
	user_id 				INT NOT NULL,
	active 					BOOLEAN NOT NULL DEFAULT 't',
	pickup_date_time		TIMESTAMP WITH TIME ZONE,
	FOREIGN KEY (user_id) REFERENCES Users(user_id) ON UPDATE CASCADE ON DELETE CASADE
);

CREATE TABLE Orders (
	order_id 	INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
	user_id 	INT NOT NULL,
	company_id 	INT NOT NULL,
	total 		REAL NOT NULL,
	confirmed 	BOOLEAN NOT NULL DEFAULT 't',
	FOREIGN KEY (user_id) REFERENCES Users(user_id) ON UPDATE CASCADE ON DELETE CASADE, 
	FOREIGN KEY (company_id) REFERENCES Users(user_id) ON UPDATE CASCADE ON DELETE CASADE
);

CREATE TABLE Messages (
	message_id 		INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
	message_time 	TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_DATE,
	to_user 		INT NOT NULL,
	from_user 		INT NOT NULL,
	message_content TEXT NOT NULL,
	FOREIGN KEY (to_user) REFERENCES Users(user_id) ON UPDATE CASCADE ON DELETE CASADE,
	FOREIGN KEY (from_user) REFERENCES Users(user_id) ON UPDATE CASCADE ON DELETE CASADE
);