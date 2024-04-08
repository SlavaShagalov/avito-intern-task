INSERT INTO users(username, password, is_admin)
VALUES ('admin', '$2a$10$F/YJirOprcmureYhionTPuiSBR8TC94SXQzxojoL1yjb2yt6SU2Qe', true),
       ('user', '$2a$10$F/YJirOprcmureYhionTPuiSBR8TC94SXQzxojoL1yjb2yt6SU2Qe', false);

INSERT INTO features(name)
VALUES ('feature_1'),
       ('feature_2'),
       ('feature_3'),
       ('feature_4'),
       ('feature_5');

INSERT INTO tags(name)
VALUES ('tag_1'),
       ('tag_2'),
       ('tag_3'),
       ('tag_4'),
       ('tag_5');

INSERT INTO banners(feature_id, content)
VALUES (1, '{
  "title": "some_title",
  "text": "some_text",
  "url": "some_url"
}'),
       (2, '{
         "title": "some_title_2",
         "info": "some_text"
       }');

INSERT INTO banner_tags(banner_id, tag_id)
VALUES (1, 1),
       (1, 2),
       (1, 3),
       (2, 4),
       (2, 5);
