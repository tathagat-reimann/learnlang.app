-- Schema initialization for learnlang

BEGIN;

CREATE TABLE IF NOT EXISTS languages (
  id   TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  code TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS packs (
  id      TEXT PRIMARY KEY,
  name    TEXT NOT NULL,
  lang_id TEXT NOT NULL REFERENCES languages(id) ON DELETE RESTRICT,
  user_id TEXT NOT NULL,
  public  BOOLEAN NOT NULL DEFAULT FALSE,
  CONSTRAINT packs_unique_per_user_lang_name UNIQUE (user_id, lang_id, name)
);

CREATE TABLE IF NOT EXISTS vocabs (
  id      TEXT PRIMARY KEY,
  image   TEXT NOT NULL,
  name    TEXT NOT NULL,
  translation TEXT NOT NULL,
  pack_id TEXT NOT NULL REFERENCES packs(id) ON DELETE CASCADE,
  CONSTRAINT vocabs_unique_per_pack_name UNIQUE (pack_id, name)
);

-- Seed languages to match current in-memory list
INSERT INTO languages (id, name, code) VALUES
  ('1', 'Hindi', 'hi')
ON CONFLICT (id) DO NOTHING;

INSERT INTO languages (id, name, code) VALUES
  ('2', 'German', 'de')
ON CONFLICT (id) DO NOTHING;

COMMIT;
