CREATE TABLE IF NOT EXISTS users(
    id char(36)  PRIMARY KEY,
    name  varchar(20) unique NOT NULL,
    email  varchar(254) unique NOT NULL,
    address varchar(100),
    password varchar(60) NOT NULL,
    profile_picture varchar(255),
    identitycard_picture varchar(255),
    is_verified varchar(3) DEFAULT 'No',
    emoney_amount decimal(10,2) DEFAULT 0,
    origincity_name varchar(30),
    origincity_id int,
    originprovince_name varchar(30),
    originprovince_id int,
    emoney_updated_at timestamp NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);
