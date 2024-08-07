CREATE TYPE plans AS ENUM ('basic', 'advanced', 'premium');
CREATE TYPE access AS ENUM ('super_admin', 'admin', 'user');
CREATE TYPE plan_status AS ENUM ('active', 'expired', 'not_payed');
CREATE TYPE game_status AS ENUM ('not_started', 'in_progress', 'ended', 'cancelled');
CREATE TYPE plan_access AS ENUM ('additional', 'holder');

CREATE TABLE users
(
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    email varchar(256) UNIQUE,
    first_name varchar(256),
    second_name varchar(256),
    description varchar(256),
    password_hash varchar(256),
    access access,
    company_name varchar(256),
    company_info varchar(256),
    company_url varchar(256),
    company_logo varchar(256),
    profile_image varchar(256)
);


CREATE TABLE verification_codes
(
    email varchar(256),
    code VARCHAR(10),
    PRIMARY KEY (email)
);

CREATE TABLE plan_invitation_codes
(
    holder_id uuid references users (id) on delete cascade,
    invitation_code varchar(256) unique
);

CREATE TABLE subscriptions
(
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY ,
    plan_type plans,
    holder_id uuid references users (id) on delete cascade,
    expiry_date timestamp,
    duration int,
    status plan_status not null DEFAULT 'active',
    invitation_code varchar(256),
    is_trial boolean DEFAULT false
);

CREATE TABLE subs_holders
(
    user_id uuid references users (id) on delete cascade,
    sub_id uuid references subscriptions (id) on delete cascade,
    access plan_access,
    primary key (user_id, sub_id)
);

CREATE TABLE topics
(
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    title varchar(256)
);

CREATE TABLE questions
(
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY ,
    topic_id uuid references topics (id) on delete cascade,
    content varchar(256)
);

CREATE TABLE games
(
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY ,
    creator_id uuid references users (id) on delete cascade,
    name varchar(256),
    start_date timestamp,
    invitation_code varchar(256),
    status game_status
);

CREATE TABLE games_users
(
    game_id uuid references games (id) on delete cascade,
    user_id uuid references users (id) on delete cascade,
    primary key (game_id, user_id)
);


CREATE TABLE results
(
    id serial PRIMARY KEY,
    user_id uuid REFERENCES users (id) ON DELETE CASCADE,
    game_id uuid REFERENCES games (id) ON DELETE CASCADE,
    name varchar(256),
    value int
);

CREATE TABLE tags
(
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    name varchar(256) UNIQUE
);

CREATE TABLE tags_questions
(
    tag_id uuid REFERENCES tags(id) ON DELETE CASCADE,
    question_id uuid REFERENCES questions(id) ON DELETE CASCADE,
    primary key (tag_id, question_id)
);

CREATE TABLE tags_results
(
    game_id uuid REFERENCES games (id) ON DELETE CASCADE,
    tag_id uuid REFERENCES tags(id) ON DELETE CASCADE,
    result_id int REFERENCES results(id) ON DELETE CASCADE
);
