-- Create "accounts" table
CREATE TABLE "public"."accounts" (
  "id" uuid NOT NULL,
  "name" character varying(120) NOT NULL,
  "currency" character(3) NOT NULL,
  "created_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "transactions" table
CREATE TABLE "public"."transactions" (
  "id" uuid NOT NULL,
  "description" character varying(255) NOT NULL,
  "created_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create "entries" table
CREATE TABLE "public"."entries" (
  "id" uuid NOT NULL,
  "transaction_id" uuid NOT NULL,
  "account_id" uuid NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_entries_account" FOREIGN KEY ("account_id") REFERENCES "public"."accounts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_entries_transaction" FOREIGN KEY ("transaction_id") REFERENCES "public"."transactions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_entries_account_id" to table: "entries"
CREATE INDEX "idx_entries_account_id" ON "public"."entries" ("account_id");
-- Create index "idx_entries_transaction_id" to table: "entries"
CREATE INDEX "idx_entries_transaction_id" ON "public"."entries" ("transaction_id");
