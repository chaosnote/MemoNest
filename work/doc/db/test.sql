
DROP PROCEDURE IF EXISTS `sp_test`;

DELIMITER $$

CREATE PROCEDURE IF NOT EXISTS `sp_test` (
    IN p_account VARCHAR(50),
    IN p_NodeID CHAR(36)
)
BEGIN
    DECLARE error_message VARCHAR(255);
    DECLARE table_articles TEXT;
    DECLARE table_categories TEXT;

    -- 測試階段用
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;
    END;

    SET error_message = CONCAT('[ERR]ERROR MESSAGE');
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;

END$$

DELIMITER ;
