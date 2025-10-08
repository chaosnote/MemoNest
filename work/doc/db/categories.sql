DROP TABLE IF EXISTS `categories` ;

CREATE TABLE IF NOT EXISTS `categories` (
    `RowID` INT AUTO_INCREMENT,
    `NodeID` CHAR(36) NOT NULL,
    `ParentID` CHAR(36) NOT NULL,
    `PathName` VARCHAR(255) NOT NULL,
    `LftIdx` INT NOT NULL,
    `RftIdx` INT NOT NULL,
    PRIMARY KEY (`RowID`),
    UNIQUE KEY `idx_node_id` (`NodeID`),
    KEY `idx_parent_id` (`ParentID`)
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4;