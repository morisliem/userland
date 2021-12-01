CREATE TABLE IF NOT EXISTS "Person" (
	"id" UUID NOT NULL,
	"fullname" VARCHAR (128) NOT NULL,
	"password" VARCHAR (128) NOT NULL,
	"email" VARCHAR (128) UNIQUE NOT NULL,
	"location" VARCHAR (256),
	"bio" VARCHAR (256),
	"web" VARCHAR(128),
	"picture" VARCHAR (128),
	"created_at" TIMESTAMP NOT NULL,
	"is_active" BOOLEAN NOT NULL,
	PRIMARY KEY ("id")
) WITHOUT OIDS;

CREATE TABLE IF NOT EXISTS "Session" (
	"id" UUID NOT NULL,
	"user_id" UUID NOT NULL,
	"ip_address" VARCHAR (16) NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	"updated_at" TIMESTAMP,
	"device" VARCHAR (20) NOT NULL,
	PRIMARY KEY ("id"),
	CONSTRAINT "fk_user_id" FOREIGN KEY ("user_id") REFERENCES "Person" ("id") ON DELETE CASCADE ON UPDATE CASCADE
) WITHOUT OIDS;

CREATE TABLE IF NOT EXISTS "User_password" (
	"id" UUID NOT NULL,
	"password" VARCHAR (128) NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	CONSTRAINT "fk_id" FOREIGN KEY ("id") REFERENCES "Person" ("id") ON DELETE CASCADE ON UPDATE CASCADE
) WITHOUT OIDS;

CREATE TABLE IF NOT EXISTS "Login_audit" (
	"username" UUID NOT NULL,
	"ip_address" VARCHAR (16) NOT NULL,
	"created_at" TIMESTAMP NOT NULL,
	CONSTRAINT "fk_username" FOREIGN KEY ("username") REFERENCES "Person" ("id") ON DELETE CASCADE ON UPDATE CASCADE
) WITHOUT OIDS;