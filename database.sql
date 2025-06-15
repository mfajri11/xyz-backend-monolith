CREATE TABLE `user` (
	`id` BIGINT NOT NULL AUTO_INCREMENT,
	`national_id` VARCHAR(16) UNIQUE,
	`full_name` VARCHAR(255) NOT NULL,
	`legal_name` VARCHAR(255),
	`birth_of_place` VARCHAR(255),
	`birth_of_date` DATE,
	`salary` DECIMAL,
	`nation_id_photo` BLOB,
	`user_photo` BLOB,
	`is_nid_valid` BOOLEAN DEFAULT FALSE
	`is_photo_valid` BOOLEAN DEFAULT FALSE
	`created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`created_by` VARCHAR(255) NOT NULL,
	`updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	`updated_by` BIGINT,
	PRIMARY KEY(`id`)
);

CREATE INDEX `user_national_id_idx`
ON `user` (`national_id`);

-- DROP TABLE `user`


CREATE TABLE `loan` (
	`id` BIGINT NOT NULL AUTO_INCREMENT UNIQUE,
	`user_id` BIGINT NOT NULL,
	`contract_number` VARCHAR(255) NOT NULL,
	`otr_amount` DECIMAL NOT NULL,
	`principal_amount` DECIMAL NOT NULL,
	`asset_name` VARCHAR(255) NOT NULL,
	`loan_type_id` TINYINT NOT NULL,
	`limit_type_id` TINYINT NOT NULL,
	`status` ENUM('ACTIVE', 'INACTIVE', 'REJECTED'),
	`start_date` DATETIME,
	`interest_rate` DECIMAL DEFAULT 0,
	PRIMARY KEY(`id`)
);
CREATE INDEX loan_user_id_idx ON loan(user_id)


-- DROP TABLE loan_payment
CREATE TABLE `loan_payment` (
	`id` BIGINT NOT NULL AUTO_INCREMENT UNIQUE,
	`loan_id` BIGINT NOT NULL,
	`amount` DECIMAL,
	`date` TIMESTAMP NOT NULL,
	`channel` VARCHAR(255) NOT NULL,
	PRIMARY KEY(`id`)
);

CREATE INDEX loan_payment_loan_id_user_id_idx ON loan_payment(loan_id)


CREATE TABLE `limit_type` (
	`id` TINYINT NOT NULL AUTO_INCREMENT UNIQUE,
	`amount` DECIMAL,
	`term` TINYINT NOT NULL,
	PRIMARY KEY(`id`)
);


CREATE TABLE `loan_type` (
	`id` TINYINT NOT NULL AUTO_INCREMENT UNIQUE,
	`name` ENUM('CAR', 'BIKE', 'WHITE_GOODS'),
	PRIMARY KEY(`id`)
);


ALTER TABLE `loan`
ADD FOREIGN KEY(`user_id`) REFERENCES `user`(`id`)
ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE `loan_payment`
ADD FOREIGN KEY(`loan_id`) REFERENCES `loan`(`id`)
ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE `loan`
ADD FOREIGN KEY(`loan_type_id`) REFERENCES `loan_type`(`id`)
ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE `loan`
ADD FOREIGN KEY(`limit_type_id`) REFERENCES `limit_type`(`id`)
ON UPDATE NO ACTION ON DELETE NO ACTION;



SELECT l.id, l.user_id, l.contract_number, l.otr_amount, l.principal_amount, l.asset_name, lot.name , lit.amount, lit.term , l.status, l.start_date, l.interest_rate
FROM xyz.loan l 
JOIN limit_type lit ON l.limit_type_id = lit.id
JOIN loan_type lot ON  l.loan_type_id  = lot.id
WHERE contract_number = '1'