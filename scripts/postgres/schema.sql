CREATE TABLE IF NOT EXISTS banners
(
    id         bigserial NOT NULL PRIMARY KEY,
    content    jsonb     NOT NULL,
    is_active  boolean   NOT NULL DEFAULT true,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS features
(
    id         bigserial NOT NULL PRIMARY KEY,
    name       varchar   NOT NULL UNIQUE,
    created_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tags
(
    id         bigserial NOT NULL PRIMARY KEY,
    name       varchar   NOT NULL UNIQUE,
    created_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS banner_references
(
    banner_id  bigint NOT NULL REFERENCES banners (id) ON DELETE CASCADE,
    feature_id bigint NOT NULL REFERENCES features (id) ON DELETE CASCADE,
    tag_id     bigint NOT NULL REFERENCES tags (id) ON DELETE CASCADE,
    PRIMARY KEY (feature_id, tag_id)
);

CREATE TABLE IF NOT EXISTS users
(
    id         bigserial NOT NULL PRIMARY KEY,
    username   text      NOT NULL UNIQUE,
    password   varchar   NOT NULL,
    is_admin   boolean   NOT NULL DEFAULT false,
    created_at timestamp NOT NULL DEFAULT now()
);
