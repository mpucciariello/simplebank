CREATE TABLE "sessions"
(
    "id"            uuid PRIMARY KEY,
    "username"      varchar        NOT NULL,
    "refresh_token" varchar        NOT NULL,
    "user_agent"    varchar        NOT NULL,
    "client_ip"     varchar UNIQUE NOT NULL,
    "is_blocked"    boolean        NOT NULL DEFAULT FALSE,
    "expires_at"    timestamp               DEFAULT (now()),
    "created_at"    timestamp               DEFAULT (now())
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
