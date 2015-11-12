CREATE TABLE packs (
  id SERIAL PRIMARY KEY,
	owner varchar(30),
	supervisor TEXT,
	from_date TIMESTAMP WITH TIME ZONE,
	to_date TIMESTAMP WITH TIME ZONE,
	delivery_status TEXT,
	tracking_code TEXT,
	pack_type TEXT,
	status INTEGER,
	pack_hash TEXT
);
