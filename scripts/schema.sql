CREATE TABLE colorkey (
  color varchar(30) DEFAULT NULL,
  level varchar(30) DEFAULT NULL
);

CREATE TABLE fruitinfo (
  id int NOT NULL AUTO_INCREMENT,
  name varchar(255) NOT NULL,
  color varchar(30) DEFAULT NULL,
  taste varchar(30) DEFAULT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE fruits (
  id int NOT NULL AUTO_INCREMENT,
  name varchar(255) NOT NULL,
  quantity int DEFAULT NULL,
  PRIMARY KEY (id)
);
