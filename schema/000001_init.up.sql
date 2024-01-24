CREATE TYPE user_role AS ENUM ('user', 'admin', 'plan_user');

CREATE TABLE users 
(
  id serial not null PRIMARY KEY,
  email varchar(256) UNIQUE,
  phone_number varchar(256),
  first_name varchar(256),
  second_name varchar(256),
  password_hash varchar(256), 
  access varchar(256),
  is_verified boolean DEFAULT false
); 

