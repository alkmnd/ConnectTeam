CREATE TYPE user_role AS ENUM ('user', 'admin', 'plan_user');
CREATE TYPE plans AS ENUM ('basic', 'advanced', 'premium');
CREATE TYPE access AS ENUM ('super_admin', 'admin', 'user');
CREATE TYPE status AS ENUM ('active', 'expired', 'on_confirm');
CREATE TYPE game_status AS ENUM ('not_started', 'in_process', 'ended');

CREATE TABLE users
(
    id serial not null PRIMARY KEY,
    email varchar(256) UNIQUE,
    first_name varchar(256),
    second_name varchar(256),
    description varchar(256),
    password_hash varchar(256),
    access access,
    is_verified boolean DEFAULT false,
    company_name varchar(256),
    company_info varchar(256),
    company_url varchar(256),
    company_logo varchar(256),
    profile_image varchar(256)
);

CREATE TABLE verification_codes
(
    user_id int REFERENCES users (id) ON DELETE CASCADE,
    code VARCHAR(10),
    PRIMARY KEY (user_id)
);

CREATE TABLE plan_invitation_codes
(
    holder_id int references users (id) on delete cascade,
    invitation_code varchar(256) unique
);

CREATE TABLE subscriptions
(
    id serial not null PRIMARY KEY,
    plan_type plans,
    user_id int references users (id) on delete cascade,
    holder_id int references users (id) on delete cascade,
    expiry_date timestamp,
    duration int,
    plan_access varchar(256),
    status status not null DEFAULT 'on_confirm',
    invitation_code varchar(256),
    is_trial boolean DEFAULT false
);

CREATE TABLE topics
(
    id serial not null PRIMARY KEY,
    title varchar(256)
);

CREATE TABLE questions
(
    id serial not null PRIMARY KEY,
    topic_id int references topics (id) on delete cascade,
    content varchar(256)
);

CREATE TABLE games
(
    id serial not null PRIMARY KEY,
    creator_id int references users (id) on delete cascade,
    name string,
    start_date timestamp,
    invitation_code string,
    status game_status
);

CREATE TABLE games_users
(
    game_id int references games (id) on delete cascade,
    user_id int references users (id) on delete cascade,
);


CREATE OR REPLACE FUNCTION update_subscription()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.expiry_date < now() AND NEW.status != 'on_confirm' THEN
    NEW.status := 'expired';
END IF;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_subscription_expiry
    AFTER INSERT OR UPDATE ON subscriptions
                        FOR EACH ROW
                        EXECUTE FUNCTION update_subscription();
