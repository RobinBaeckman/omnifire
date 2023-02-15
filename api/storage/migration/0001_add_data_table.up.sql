CREATE TABLE data (
	id serial primary key,
   body TEXT not null,
   created_at TIMESTAMP default now()
);
