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

-- SP
DROP PROCEDURE IF EXISTS `insert_article`;

DELIMITER $$
CREATE PROCEDURE IF NOT EXISTS `insert_article` (
    IN p_Title VARCHAR(255),
    IN p_Content TEXT,
    IN p_UpdateDt DATETIME,
    IN p_CreatedDt DATETIME,
    IN p_NodeID CHAR(36)
)
BEGIN
    DECLARE error_message VARCHAR(255);
    -- 驗證 NodeID 是否存在
    IF EXISTS (SELECT 1 FROM categories WHERE NodeID = p_NodeID) THEN
        INSERT INTO articles (Title, Content, UpdateDt, CreatedDt, NodeID)
        VALUES (p_Title, p_Content, p_UpdateDt, p_CreatedDt, p_NodeID);
    ELSE
        SET error_message = CONCAT('NodeID(', p_NodeID, ') 不存在於 `categories`');
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message ;
    END IF;
END $$
DELIMITER ;

-- CALL insert_article('AAA', '123', NOW(), NOW(), '002198c1-1efe-4111-8c1c-d74d293b5823') ;