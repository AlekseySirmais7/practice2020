DROP DATABASE myService;

CREATE DATABASE myService
    WITH
    OWNER = postgres
    CONNECTION LIMIT = -1
    ;


CREATE COLLATION posix (LOCALE = 'POSIX');

CREATE unlogged TABLE Users (
	login text PRIMARY KEY NOT NULL COLLATE posix UNIQUE,
	email text NOT NULL UNIQUE,
	password bytea NOT NULL
);

CREATE UNIQUE INDEX Users_email ON Users ( LOWER(email), LOWER(login) );


CREATE unlogged TABLE Posts (
    id serial NOT NULL PRIMARY KEY,
    title text NOT NULL UNIQUE, --- пусть навзвания не повторяются
    ttext text NOT NULL,
	author text NOT NULL,
	created timestamp with time zone default NULL
);

CREATE INDEX Posts_created ON Posts (created);


INSERT INTO Posts (title, ttext, author, created) VALUES ( 'Крупнейшие ИТ-компании предупредили об угрозе из-за дела о картеле',
 'Ассоциация предприятий компьютерных и информационных технологий (АПКИТ, среди ее членов — 1 C, ABBYY, Acer, IBM, Лаборатория Касперского и др.) в открытом письме заявила, что ситуация с уголовным преследованием ИТ-предпринимателей воспринимается как угроза рынку',
 'Alex_Sir', now() - interval '3 hour' ) ;

INSERT INTO Posts (title, ttext, author, created) VALUES ( 'Продам гараж',
 'Продается отличный гараж. Есть все удобства, можно жить. 8 800 555 35 **',
 'Nick', now() - interval '1 hour' ) ;

