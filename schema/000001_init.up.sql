CREATE TABLE users 
(
  id serial not null PRIMARY KEY,
  email varchar(256),
  phone_number varchar(256),
  first_name varchar(256),
  second_name varchar(256),
  password_hash varchar(256)
);