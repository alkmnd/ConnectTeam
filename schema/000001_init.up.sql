CREATE TYPE user_role AS ENUM ('user', 'admin', 'plan_user');
CREATE TYPE plans AS ENUM ('basic', 'advanced', 'premium');
CREATE TYPE access AS ENUM ('super_admin', 'admin', 'user');
CREATE TYPE plan_status AS ENUM ('active', 'expired', 'on_confirm');
CREATE TYPE game_status AS ENUM ('not_started', 'in_progress', 'ended');

CREATE TABLE users
(
    id uuid DEFAULT gen_random_uuid(),
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
    id uuid DEFAULT gen_random_uuid(),
    plan_type plans,
    user_id uuid references users (id) on delete cascade,
    holder_id uuid references users (id) on delete cascade,
    expiry_date timestamp,
    duration int,
    plan_access varchar(256),
    status plan_status not null DEFAULT 'on_confirm',
    invitation_code varchar(256),
    is_trial boolean DEFAULT false
);

CREATE TABLE topics
(
    id uuid DEFAULT gen_random_uuid(),
    title varchar(256)
);

CREATE TABLE questions
(
    id uuid DEFAULT gen_random_uuid(),
    topic_id int references topics (id) on delete cascade,
    content varchar(256)
);

CREATE TABLE games
(
    id uuid DEFAULT gen_random_uuid(),
    creator_id int references users (id) on delete cascade,
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
    user_id uuid REFERENCES users (id) ON DELETE CASCADE,
    game_id uuid REFERENCES games (id) ON DELETE CASCADE,
    primary key (game_id, user_id),
    value int
);

CREATE TABLE tags
(
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    name varchar(256)
);

CREATE TABLE tags_questions
(
    tag_id uuid REFERENCES tags(id) ON DELETE CASCADE,
    question_id uuid REFERENCES questions(id) ON DELETE CASCADE,
    primary key (tag_id, question_id)
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
