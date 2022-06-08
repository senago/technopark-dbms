-- citext is a case-insensitive character string type
CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE IF NOT EXISTS users (
  id bigserial,
  nickname citext COLLATE "ucs_basic" NOT NULL UNIQUE PRIMARY KEY,
  fullname text NOT NULL,
  about text,
  email citext NOT NULL UNIQUE
);