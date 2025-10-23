DROP TABLE IF EXISTS `articles` ;

CREATE TABLE IF NOT EXISTS `articles` (
    `RowID` INT AUTO_INCREMENT,
    `Title` VARCHAR(255) NOT NULL,
    `Content` TEXT,
    `UpdateDt` DATETIME NOT NULL,
    `CreatedDt` DATETIME NOT NULL,
    `NodeID` CHAR(36) NOT NULL,
    PRIMARY KEY (`RowID`),
    KEY `k_node_id` (`NodeID`)
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4;