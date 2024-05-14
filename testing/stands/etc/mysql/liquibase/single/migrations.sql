--liquibase formatted sql

--changeset agurinov:SCHEMA1
CREATE TABLE `foobar`
(
	`uuid` binary(16) NOT NULL,
	PRIMARY KEY (`uuid`)
) ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;
--rollback DROP TABLE `foobar`;

--changeset agurinov:SCHEMA2
CREATE TABLE `barbaz`
(
	`uuid` binary(16) NOT NULL,
	PRIMARY KEY (`uuid`)
) ENGINE = InnoDB
DEFAULT CHARSET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;
--rollback DROP TABLE `barbaz`;
