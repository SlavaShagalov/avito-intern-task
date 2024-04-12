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

INSERT INTO banners(content, is_active)
VALUES ('{
  "title": "banner_1",
  "info": "banner_1 info"
}', true),
       ('{
         "title": "banner_2",
         "info": "banner_2 info"
       }', true),
       ('{
         "title": "banner_3",
         "info": "banner_3 info"
       }', false);

INSERT INTO banner_references(banner_id, feature_id, tag_id)
VALUES (1, 1, 1),
       (1, 1, 2),
       (1, 1, 3),
       (2, 2, 4),
       (2, 2, 5),
       (3, 1, 4);
