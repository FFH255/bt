CREATE TABLE IF NOT EXISTS statistics (
     id              SERIAL PRIMARY KEY,
     user_id         INTEGER NOT NULL REFERENCES users(id),
     wpm             DOUBLE PRECISION NOT NULL,
     cpm             DOUBLE PRECISION NOT NULL,
     accuracy        DOUBLE PRECISION NOT NULL,
     duration        BIGINT NOT NULL,
     played_at       TIMESTAMPTZ NOT NULL,
     language        VARCHAR(128) NOT NULL,
     mode            VARCHAR(16) NOT NULL,
     sub_mode        VARCHAR(16) NOT NULL,
     is_punctuation  BOOLEAN NOT NULL,
     uncompleted_tests_count INTEGER NOT NULL,
     uncompleted_tests_total_duration BIGINT NOT NULL,
     idempotency_key VARCHAR(255) NOT NULL
);
