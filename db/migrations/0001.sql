-- +migrate Up

CREATE TABLE "Question" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "content" text,
  "amount_answer" int,
  "url" varchar,
  "vote" int
);

CREATE TABLE "Answer" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "content" text,
  "vote" int,
  "question_id" int
);

-- CREATE TABLE "Document" (
--   "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
--   "content" text,
--   "vote" int,
--   "type" int,
--   "question_id" int
-- ) 

ALTER TABLE "Answer" ADD FOREIGN KEY ("question_id") REFERENCES "Question" ("id");

-- ALTER TABLE "Document" ADD FOREIGN KEY ("question_id") REFERENCES "Question" ("id")
