-- name: create-tables
CREATE TABLE `teams` (
	`id` INT NOT NULL ,
	`name` varchar(45) UNIQUE,
	`active` INT(1) NOT NULL DEFAULT '1',
	`groupid` varchar(2) NOT NULL,
	`point` INT(1) NOT NULL DEFAULT '0',
	`win` INT NOT NULL DEFAULT '0',
	`drawn` INT NOT NULL DEFAULT '0',
	`lose` INT NOT NULL DEFAULT '0',
	`goalfor` INT NOT NULL DEFAULT '0',
	`goalagainst` INT NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`)
);

CREATE TABLE `matches` (
	`id` INT NOT NULL ,
	`stage` INT NOT NULL DEFAULT '0',
	`date` DATETIME NOT NULL,
	`teama` INT NOT NULL,
	`scorea` INT NOT NULL DEFAULT '0',
	`pena` INT NOT NULL DEFAULT '0',
	`teamb` INT NOT NULL,
	`scoreb` INT NOT NULL DEFAULT '0',
	`penb` INT NOT NULL DEFAULT '0',
	`stadium` INT NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`)
	FOREIGN KEY (`stage`)   REFERENCES `stage`(`id`)  
	FOREIGN KEY (`teama`)   REFERENCES `teams`(`id`)  
	FOREIGN KEY (`stadium`) REFERENCES `stadium`(`id`)
	FOREIGN KEY (`teamb`)   REFERENCES `teams`(`id`)  
);

CREATE TABLE `users` (
	`id` INT NOT NULL ,
	`token` varchar(80) DEFAULT NULL,
	`firstname` varchar(45) NOT NULL,
	`lastname` varchar(45) NOT NULL,
	`password` varchar(80) NOT NULL,
	`email` varchar(80) NOT NULL,
	`score` INT NOT NULL DEFAULT '0',
	`champion` INT NOT NULL DEFAULT '0',
	`champchange` INT NOT NULL DEFAULT '0',
	`position` INT NOT NULL DEFAULT '0',
	`positionbefore` INT NOT NULL DEFAULT '0',
	`malus` INT NOT NULL DEFAULT '0',
	`bonus` INT NOT NULL DEFAULT '0',
	`ngoodscores` INT NOT NULL DEFAULT '0',
	`ngoodwinner` INT NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`)
	FOREIGN KEY (`champion`) REFERENCES `teams`(`id`)
);

CREATE TABLE `stadium` (
	`id` INT NOT NULL,
	`name` varchar(45) NOT NULL,
	PRIMARY KEY (`id`)
);

CREATE TABLE `stage` (
	`id` INT NOT NULL ,
	`name` varchar(45) NOT NULL,
	`active` INT NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`)
);

CREATE TABLE `config` (
	`lastmatch` INT NOT NULL DEFAULT '0',
	`refresh` INT NOT NULL DEFAULT '0',
	`stage` INT NOT NULL,
	FOREIGN KEY (`lastmatch`) REFERENCES `matches`(`id`)
	FOREIGN KEY (`stage`) REFERENCES `stage`(`id`)
);

CREATE TABLE `pronostics` (
	`id` INT NOT NULL ,
	`user` INT NOT NULL,
	`match` INT NOT NULL,
	`scorea` INT NOT NULL DEFAULT '0',
	`scoreb` INT NOT NULL DEFAULT '0',
	`pena` INT NOT NULL DEFAULT '0',
	`penb` INT NOT NULL DEFAULT '0',
	PRIMARY KEY (`id`)
	FOREIGN KEY (`user`) REFERENCES `users`(`id`) 
	FOREIGN KEY (`match`) REFERENCES `matches`(`id`)
);

CREATE TABLE `stats` (
	`nombrebut` INT NOT NULL DEFAULT '0',
	`goodprono` INT NOT NULL DEFAULT '0',
	`goodwinner` INT NOT NULL DEFAULT '0',
	`redcard` INT NOT NULL DEFAULT '0',
	`yellowcard` INT NOT NULL DEFAULT '0'
);