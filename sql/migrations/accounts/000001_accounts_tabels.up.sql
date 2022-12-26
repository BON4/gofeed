CREATE TYPE account_role AS ENUM('basic', 'admin');

CREATE TABLE Accounts (
	username text not null unique PRIMARY KEY,
  	password bytea not null,
  	email text unique not null,
	role account_role not null default 'basic'
);

CREATE TABLE Users (
       uuid text not null unique PRIMARY KEY,
       ip text,
       os text,
       browser text,
       account_username text not null
);

ALTER TABLE Users
	ADD FOREIGN KEY (account_username) REFERENCES Accounts (username);
