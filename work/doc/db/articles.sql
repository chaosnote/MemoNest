CREATE TABLE `articles` (
    `ArticleID` INT AUTO_INCREMENT,
    `Title` VARCHAR(255) NOT NULL,
    `Content` TEXT,
    `UpdateDt` DATETIME NOT NULL,
    `CreatedDt` DATETIME NOT NULL,
    `NodeID` CHAR(36) NOT NULL,
    PRIMARY KEY (`ArticleID`),
    KEY `idx_node_id` (`NodeID`),
    FOREIGN KEY (`NodeID`) REFERENCES `categories`(`NodeID`)
) ENGINE=INNODB DEFAULT CHARSET=utf8mb4;

DROP TABLE `articles` ;