
-- USERS
INSERT INTO rs_user(user_id, rf_id, password, first_name, last_name, email)
VALUES('tuommii', null, '$2a$10$hbTUJYT0kwBkhMHKweV36eWFVw2CqL7v0FVV14UKCl1xbP1rBFzpG', 'Miikka', 'Tuominen', 'm@t.fi');

INSERT INTO rs_user(user_id, rf_id, password, first_name, last_name, email)
VALUES('raccoon', ' D9 16 D1 B9', '$2a$10$hbTUJYT0kwBkhMHKweV36eWFVw2CqL7v0FVV14UKCl1xbP1rBFzpG',
'Raccoon', 'Hustler', 'h@t.fi');

INSERT INTO rs_user(user_id, rf_id, password, first_name, last_name, email)
VALUES('turtle', '', '$2a$10$hbTUJYT0kwBkhMHKweV36eWFVw2CqL7v0FVV14UKCl1xbP1rBFzpG',
'Mr', 'Turtle', 'm@d.fi');

INSERT INTO rs_user(user_id, rf_id, password, first_name, last_name, email)
VALUES('bstinson', null, '$2a$10$hbTUJYT0kwBkhMHKweV36eWFVw2CqL7v0FVV14UKCl1xbP1rBFzpG', 'Barney', 'Stinson', 'b@s.fi');

INSERT INTO rs_user(user_id, rf_id, password, first_name, last_name, email)
VALUES('jack', null, '$2a$10$hbTUJYT0kwBkhMHKweV36eWFVw2CqL7v0FVV14UKCl1xbP1rBFzpG', 'Jack', 'Bauer', 'j@b.fi');

INSERT INTO rs_user(user_id, rf_id, password, first_name, last_name, email)
VALUES('ebachman', null, '$2a$10$hbTUJYT0kwBkhMHKweV36eWFVw2CqL7v0FVV14UKCl1xbP1rBFzpG', 'Erlich', 'Bachman', 'e@b.fi');

INSERT INTO rs_user(user_id, rf_id, password, first_name, last_name, email)
VALUES('bhobbes', null, '$2a$10$hbTUJYT0kwBkhMHKweV36eWFVw2CqL7v0FVV14UKCl1xbP1rBFzpG', 'Bobby', 'Hobbes', 'bob@h.fi');

INSERT INTO rs_user(user_id, rf_id, password, first_name, last_name, email)
VALUES('mflinkman', null, '$2a$10$hbTUJYT0kwBkhMHKweV36eWFVw2CqL7v0FVV14UKCl1xbP1rBFzpG', 'Marshall', 'Flinkman', 'm@flinkman.fi');


-- SPORTS
INSERT INTO rs_sport(sport_name) VALUES('Chess');
INSERT INTO rs_sport(sport_name) VALUES('Table Tennis');
INSERT INTO rs_sport(sport_name) VALUES('Madden 21');
INSERT INTO rs_sport(sport_name) VALUES('FIFA 21');


-- GAMES
INSERT INTO rs_game(sport_name) VALUES('Chess');
INSERT INTO rs_game(sport_name) VALUES('Table Tennis');
INSERT INTO rs_game(sport_name) VALUES('Madden 21');
INSERT INTO rs_game(sport_name) VALUES('Chess');
INSERT INTO rs_game(sport_name) VALUES('Chess');
INSERT INTO rs_game(sport_name) VALUES('Madden 21');
INSERT INTO rs_game(sport_name) VALUES('Chess');
INSERT INTO rs_game(sport_name) VALUES('Ping Pong');
INSERT INTO rs_game(sport_name) VALUES('Chess');
INSERT INTO rs_game(sport_name) VALUES('Ping Pong');
INSERT INTO rs_game(sport_name) VALUES('Ping Pong');
INSERT INTO rs_game(sport_name) VALUES('Ping Pong');
INSERT INTO rs_game(sport_name) VALUES('Chess');

-- ADD PLAYERS TO GAMES
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(1, 'tuommii', 0);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(1, 'bhobbes', 1);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(2, 'bhobbes', 1);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(2, 'jack', 0);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(3, 'jack', 2);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(3, 'bhobbes', 2);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(4, 'jack', 1);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(4, 'tuommii', 0);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(5, 'ebachman', 1);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(5, 'mflinkman', 0);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(6, 'ebachman', 1);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(6, 'mflinkman', 0);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(7, 'mflinkman', 1);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(7, 'ebachman', 0);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(8, 'jack', 1);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(8, 'raccoon', 0);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(9, 'jack', 1);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(9, 'raccoon', 0);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(10, 'jack', 1);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(10, 'raccoon', 0);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(11, 'jack', 1);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(11, 'raccoon', 0);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(12, 'jack', 0);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(12, 'raccoon', 1);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(13, 'tuommii', 1);
INSERT INTO rs_game_players(game_id, user_id, result) VALUES(13, 'bhobbes', 0);


-- UNKNOWN TAGS
-- INSERT INTO rs_audit_trail(event, subject) VALUES('unknown_tag', 'A1 A2');
-- INSERT INTO rs_audit_trail(event, subject) VALUES('unknown_tag', 'A1 A3');
-- INSERT INTO rs_audit_trail(event, subject) VALUES('unknown_tag', 'A1 A4');
-- INSERT INTO rs_audit_trail(event, subject) VALUES('unknown_tag', 'A1 A5');
-- INSERT INTO rs_audit_trail(event, subject) VALUES('unknown_tag', 'A1 A6');
-- INSERT INTO rs_audit_trail(event, subject) VALUES('unknown_tag', 'A1 A7');
