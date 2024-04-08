CREATE TABLE IF NOT EXISTS features
(
    id         serial    NOT NULL PRIMARY KEY,
    name       varchar   NOT NULL UNIQUE,
    created_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS banners
(
    id         bigserial NOT NULL PRIMARY KEY,
    feature_id int       NOT NULL REFERENCES features (id) ON DELETE CASCADE,
    content    jsonb     NOT NULL,
    is_active  boolean   NOT NULL DEFAULT true,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tags
(
    id         serial    NOT NULL PRIMARY KEY,
    name       varchar   NOT NULL UNIQUE,
    created_at timestamp NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS banner_tags
(
    banner_id bigint NOT NULL REFERENCES banners (id) ON DELETE CASCADE,
    tag_id    int    NOT NULL REFERENCES tags (id) ON DELETE CASCADE,
    PRIMARY KEY (banner_id, tag_id)
);

CREATE TABLE IF NOT EXISTS users
(
    id         bigserial NOT NULL PRIMARY KEY,
    username   text      NOT NULL UNIQUE,
    password   varchar   NOT NULL,
    is_admin   boolean   NOT NULL DEFAULT false,
    created_at timestamp NOT NULL DEFAULT now()
);
