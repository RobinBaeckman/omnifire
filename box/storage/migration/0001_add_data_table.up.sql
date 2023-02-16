CREATE TABLE box (
	id serial primary key,
   body TEXT not null,
   created TIMESTAMP default now()
);
