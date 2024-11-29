CREATE TABLE IF NOT EXISTS users(
    id uuid  PRIMARY KEY,
    name  varchar(20) unique NOT NULL,
    email  varchar(254) unique NOT NULL,
    address varchar(100),
    password varchar(60) NOT NULL,
    profile_picture varchar(255),
    identitycard_picture varchar(255),
    emoney_amount decimal(10,2) DEFAULT 0,
    origincity_name varchar(255),
    originprovince_name varchar(255),
    emoney_updated_at timestamp NOT NULL,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL
);
