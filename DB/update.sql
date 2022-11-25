--name: fill-tables
UPDATE `matches` SET scorea=0, scoreb=2 where id=1;
UPDATE `matches` SET scorea=0, scoreb=2 where id=2;
UPDATE `matches` SET scorea=6, scoreb=2 where id=3;
UPDATE `matches` SET scorea=1, scoreb=1 where id=4;
UPDATE `matches` SET scorea=1, scoreb=2 where id=5;
UPDATE `matches` SET scorea=0, scoreb=0 where id=6;
UPDATE `matches` SET scorea=0, scoreb=0 where id=7;
UPDATE `matches` SET scorea=4, scoreb=1 where id=8;
UPDATE `matches` SET scorea=0, scoreb=0 where id=9;
UPDATE `matches` SET scorea=1, scoreb=2 where id=10;
UPDATE `matches` SET scorea=7, scoreb=0 where id=11;
UPDATE `matches` SET scorea=1, scoreb=0 where id=12;
UPDATE `matches` SET scorea=1, scoreb=0 where id=13;
UPDATE `matches` SET scorea=0, scoreb=0 where id=14;
UPDATE `matches` SET scorea=3, scoreb=2 where id=15;
UPDATE `matches` SET scorea=2, scoreb=0 where id=16;
UPDATE `config` SET stage =2 where stage=1;