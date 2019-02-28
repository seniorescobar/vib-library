CREATE TABLE member(
    id SERIAL PRIMARY KEY NOT NULL,
    first_name VARCHAR(128) NOT NULL,
    last_name VARCHAR(128) NOT NULL
);

CREATE TABLE book(
    id SERIAL PRIMARY KEY NOT NULL,
    title VARCHAR(128) NOT NULL,
    author VARCHAR(128) NOT NULL,
    year INT NOT NULL
);

CREATE TABLE rent(
    id SERIAL  PRIMARY KEY NOT NULL,
    member_id INT REFERENCES member(id),
    book_id INT REFERENCES book(id),
    time TIMESTAMP NOT NULL,
    time_return TIMESTAMP
);

CREATE TABLE availability(
    book_id INT PRIMARY KEY NOT NULL REFERENCES book(id),
    amount INT NOT NULL
);
