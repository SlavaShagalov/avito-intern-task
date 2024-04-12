INSERT INTO users(username, password, is_admin)
VALUES ('admin', '$2a$10$F/YJirOprcmureYhionTPuiSBR8TC94SXQzxojoL1yjb2yt6SU2Qe', true),
       ('user', '$2a$10$F/YJirOprcmureYhionTPuiSBR8TC94SXQzxojoL1yjb2yt6SU2Qe', false);

INSERT INTO features(name)
VALUES ('feature_1'),
       ('feature_2'),
       ('feature_3'),
       ('feature_4'),
       ('feature_5');

INSERT INTO banners(feature_id, tag_ids, content)
VALUES (1, ARRAY [1, 2, 3]::integer[], '{
  "title": "banner_1",
  "info": "banner_1 info"
}'),
       (2, ARRAY [4, 5]::integer[], '{
         "title": "banner_2",
         "info": "banner_2 info"
       }'),
       (1, ARRAY [2]::integer[], '{
         "title": "banner_3",
         "info": "banner_3 info"
       }'),
       (1, ARRAY []::integer[], '{
         "title": "banner_4",
         "info": "banner_4 info"
       }');
