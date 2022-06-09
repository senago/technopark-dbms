-- citext is a case-insensitive character string type
CREATE EXTENSION IF NOT EXISTS citext;

CREATE UNLOGGED TABLE IF NOT EXISTS users (
  id bigserial,
  nickname citext COLLATE "ucs_basic" NOT NULL PRIMARY KEY,
  fullname text NOT NULL,
  about text,
  email citext NOT NULL UNIQUE
);

CREATE UNLOGGED TABLE IF NOT EXISTS forums (
  id bigserial,
  title text NOT NULL,
  "user" citext COLLATE "ucs_basic" NOT NULL REFERENCES users (nickname),
  slug citext NOT NULL PRIMARY KEY,
  posts bigint DEFAULT 0,
  threads bigint DEFAULT 0
);

CREATE UNLOGGED TABLE IF NOT EXISTS threads (
  id bigserial PRIMARY KEY,
  title text NOT NULL,
  author citext COLLATE "ucs_basic" NOT NULL REFERENCES users (nickname),
  forum citext NOT NULL REFERENCES forums (slug),
  message text NOT NULL,
  votes integer DEFAULT 0,
  slug citext NOT NULL,
  created timestamp with time zone DEFAULT now()
);

CREATE UNLOGGED TABLE IF NOT EXISTS posts (
  id bigserial NOT NULL PRIMARY KEY,
  parent integer DEFAULT 0,
  author citext COLLATE "ucs_basic" NOT NULL REFERENCES users (nickname),
  message text NOT NULL,
  is_edited boolean DEFAULT FALSE,
  forum citext NOT NULL REFERENCES forums (slug),
  thread integer REFERENCES threads (id),
  created timestamp with time zone DEFAULT now(),
  path bigint [] DEFAULT ARRAY [] :: INTEGER []
);

CREATE UNLOGGED TABLE IF NOT EXISTS forum_users (
  nickname citext COLLATE "ucs_basic" NOT NULL REFERENCES users (nickname),
  fullname text NOT NULL,
  about text,
  email citext NOT NULL,
  forum citext NOT NULL REFERENCES forums (slug),
  CONSTRAINT forum_users_key UNIQUE (nickname, forum)
);

CREATE UNLOGGED TABLE IF NOT EXISTS votes (
  nickname citext COLLATE "ucs_basic" NOT NULL REFERENCES users (nickname),
  thread int NOT NULL REFERENCES threads (id),
  voice integer NOT NULL
);

-- Functions and Triggers

CREATE OR REPLACE FUNCTION set_threads_votes() RETURNS TRIGGER AS $$
  BEGIN
    UPDATE threads SET votes = votes + NEW.voice WHERE id = NEW.thread_id;
    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER insert_votes AFTER INSERT ON votes FOR EACH ROW EXECUTE PROCEDURE set_threads_votes();


CREATE OR REPLACE FUNCTION update_threads_votes() RETURNS TRIGGER AS $$
  BEGIN
    UPDATE threads SET votes = votes + NEW.voice - OLD.voice WHERE id = NEW.thread_id;
    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_votes AFTER UPDATE ON votes FOR EACH ROW EXECUTE PROCEDURE update_threads_votes();


CREATE OR REPLACE FUNCTION update_post_path() RETURNS TRIGGER AS $$
  BEGIN
    new.path = (SELECT path FROM posts WHERE id = new.parent) || new.id;
    RETURN new;
  END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_path BEFORE INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE update_post_path();


CREATE OR REPLACE FUNCTION count_forum_threads() RETURNS TRIGGER AS $$
  BEGIN
    UPDATE forums SET threads = forums.threads + 1 WHERE slug = NEW.forum;
    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_count_threads AFTER INSERT ON threads FOR EACH ROW EXECUTE PROCEDURE count_forum_threads();


CREATE OR REPLACE FUNCTION count_forum_posts() RETURNS TRIGGER AS $$
  BEGIN
    UPDATE forums SET posts = forums.posts + 1 WHERE slug = NEW.forum;
    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_count_posts AFTER INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE count_forum_posts();


CREATE OR REPLACE FUNCTION update_forum_user() RETURNS TRIGGER AS $$
DECLARE
    nickname citext;
    fullname text;
    about    text;
    email    citext;
  BEGIN
    SELECT u.nickname, u.fullname, u.about, u.email FROM users u WHERE u.nickname = NEW.author
    INTO nickname, fullname, about, email;

    INSERT INTO forum_users (nickname, fullname, about, email, forum)
    VALUES (nickname, fullname, about, email, NEW.forum)
    ON CONFLICT do nothing;

    RETURN NEW;
  END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_forum_users_on_post AFTER INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE update_forum_user();
CREATE OR REPLACE TRIGGER update_forum_users_on_thread AFTER INSERT ON threads FOR EACH ROW EXECUTE PROCEDURE update_forum_user();

-- Vacuum for better performance
VACUUM ANALYSE;