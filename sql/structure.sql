
CREATE TABLE IF NOT EXISTS users (
		id			BIGSERIAL PRIMARY KEY,
		name		TEXT NOT NULL,
		phone		TEXT NOT NULL UNIQUE,
        login       TEXT NOT NULL,
        password    TEXT NOT NULL,
		active		BOOLEAN NOT NULL DEFAULT TRUE,
		created		TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

    DROP TABLE users;
    DROP TABLE users_tokens;
    DROP TABLE  user_courses;
    DROP TABLE courses;

CREATE TABLE IF NOT EXISTS managers (
		id			BIGSERIAL PRIMARY KEY,
		name		TEXT NOT NULL,
        phone		TEXT NOT NULL UNIQUE,
		login       TEXT NOT NULL UNIQUE,
        password    TEXT NOT NULL,
        is_admin   BOOLEAN   NOT NULL DEFAULT FALSE,
        active		BOOLEAN NOT NULL DEFAULT TRUE,
		created		TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS users_tokens (
    token       TEXT not NULL UNIQUE,
    user_id BIGINT NOT NULL REFERENCES users,
    expire      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created		TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS courses (
    id BIGSERIAL PRIMARY KEY,
    manager_id BIGINT NOT NULL REFERENCES managers,
    course_name TEXT NOT NULL UNIQUE,
    price INTEGER NOT NULL CHECK (price > 0) DEFAULT 1,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created timestamp not null DEFAULT CURRENT_timestamp 
);


DROP TABLE user_courses;



CREATE TABLE IF NOT EXISTS user_courses (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users,
    courses_id BIGINT NOT NULL REFERENCES courses,
    course_name TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0) DEFAULT 1,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created timestamp not null DEFAULT CURRENT_timestamp 
);



