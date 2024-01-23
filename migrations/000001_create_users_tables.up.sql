CREATE TABLE user_profile
(
    id                UUID         NOT NULL PRIMARY KEY,
    nickname          VARCHAR(255) NOT NULL,
    date_registration TIMESTAMP,
    photo_link        VARCHAR(255),
    coins             INT
);

CREATE TABLE user_data
(
    id               SERIAL PRIMARY KEY,
    profile_id       UUID UNIQUE REFERENCES user_profile (id) ON DELETE CASCADE,
    password_encoded VARCHAR(255),
    password_salt    VARCHAR(255)
);

CREATE TABLE user_contacts
(
    id                 SERIAL PRIMARY KEY,
    profile_id         UUID UNIQUE REFERENCES user_profile (id) ON DELETE CASCADE,
    email              VARCHAR(255),
    email_subscription BOOLEAN
);

CREATE TABLE refresh_sessions
(
    id               SERIAL PRIMARY KEY,
    profile_id       UUID REFERENCES user_profile (id) ON DELETE CASCADE,
    refresh_token_id UUID                     NOT NULL,
    issued_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    expires_in       TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE email_verification
(
    id    SERIAL PRIMARY KEY,
    email VARCHAR(255),
    code  INT
);