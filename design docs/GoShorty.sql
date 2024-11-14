SET TIMEZONE = 'Asia/Kolkata';

CREATE TABLE "urls" (
  "key" varchar PRIMARY KEY,
  "long_url" varchar NOT NULL,
  "custom_alias" varchar UNIQUE,
  "creation_date" timestamptz NOT NULL DEFAULT (date_trunc('second', now()::timestamptz)),
  "expiry_date" timestamptz
);

CREATE INDEX "ShortURLKeys" ON "urls" USING BTREE ("key");

CREATE INDEX "ShortURLAliases" ON "urls" USING BTREE ("custom_alias");

COMMENT ON COLUMN "urls"."key" IS 'short url, identifier for long urls';

COMMENT ON COLUMN "urls"."long_url" IS 'must have some valid url value';

COMMENT ON COLUMN "urls"."custom_alias" IS 'is optional and aliases must be unique';

COMMENT ON COLUMN "urls"."expiry_date" IS 'optional field';
