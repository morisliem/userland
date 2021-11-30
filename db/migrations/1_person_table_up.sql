CREATE TABLE IF NOT EXISTS "Person"{
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
} WITHOUT OIDS;