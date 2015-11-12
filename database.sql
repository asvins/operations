CREATE TABLE packs (
	owner varchar(30) CONSTRAINT subscription_pk PRIMARY KEY,
	supervisor TEXT,
	from_date TIMESTAMP WITH TIME ZONE,
	to_date TIMESTAMP WITH TIME ZONE,
	delivery_status TEXT,
	tracking_code TEXT
	pack_type TEXT,
	pack_hash TEXT
);
