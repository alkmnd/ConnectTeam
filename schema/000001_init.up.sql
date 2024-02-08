CREATE TYPE user_role AS ENUM ('user', 'admin', 'plan_user');

CREATE TYPE plans AS ENUM ('basic', 'advanced', 'premium');

CREATE TABLE users 
(
  id serial not null PRIMARY KEY,
  email varchar(256) UNIQUE,
  -- phone_number varchar(256),
  first_name varchar(256),
  second_name varchar(256),
  description varchar(256),
  password_hash varchar(256), 
  access varchar(256),
  is_verified boolean DEFAULT false,
  company_name varchar(256),
  company_info varchar(256),
  company_url varchar(256),
  company_logo varchar(256),
  profile_image varchar(256)
); 

CREATE TABLE verification_codes 
(
    user_id varchar(256) UNIQUE,
    code VARCHAR(10)
);

CREATE TABLE plans_users 
(
  plan_type plans,
  user_id int PRIMARY KEY,
  holder_id int,
  expiry_date timestamp,
  duration int,
  plan_access varchar(256), 
  confirmed boolean
);

CREATE TABLE plan_requests 
(id serial not null PRIMARY KEY
  ,
  user_id int,
  plan_type plans, 
  duration int, 
  request_date timestamp
);