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
    tag_ids    int[]     NOT NULL,
    content    jsonb     NOT NULL,
    is_active  boolean   NOT NULL DEFAULT true,
    created_at timestamp NOT NULL DEFAULT now(),
    updated_at timestamp NOT NULL DEFAULT now()
);

-- CREATE OR REPLACE FUNCTION check_unique_feature_tag_combination(_feature_id int, _tag_ids int)
--     RETURNS BOOLEAN AS
-- $$
-- BEGIN
--     RETURN NOT EXISTS(
--             SELECT id
--             FROM banners
--             WHERE feature_id = _feature_id
--               AND _tag_ids = ANY (tag_ids)
--         );
-- END;
-- $$ LANGUAGE plpgsql;
--
-- ALTER TABLE banners
--     ADD CONSTRAINT check_unique_feature_tag_combination_constraint
--         CHECK (check_unique_feature_tag_combination(feature_id, tag_ids));

-- CREATE UNIQUE INDEX banners_feature_id_tag_ids_idx
--     ON banners (feature_id, unnest(tag_ids));

CREATE TABLE IF NOT EXISTS users
(
    id         bigserial NOT NULL PRIMARY KEY,
    username   text      NOT NULL UNIQUE,
    password   varchar   NOT NULL,
    is_admin   boolean   NOT NULL DEFAULT false,
    created_at timestamp NOT NULL DEFAULT now()
);
