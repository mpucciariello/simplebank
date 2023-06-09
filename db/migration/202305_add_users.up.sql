CREATE TABLE "users" (
    "username" varchar PRIMARY KEY,
    "hashed_password" varchar NOT NULL,
    "full_name" varchar NOT NULL,
    "password_changed_at" varchar NOT NULL,
    "created_at" timestamp DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");
