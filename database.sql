create database testdb;
use testdb;

create table filmat (
  id int auto_increment primary key,
  emri varchar(30),
  studio varchar(30),
  kategoria varchar(30),
  viti int(4),
  kohezgjatja varchar(5),
  vleresimi decima(2,1)
)
