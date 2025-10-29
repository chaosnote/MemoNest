
DROP PROCEDURE IF EXISTS `sp_test`;

DELIMITER $$

CREATE PROCEDURE IF NOT EXISTS `sp_test` (
    IN p_account VARCHAR(50)
)
BEGIN
    DECLARE error_message VARCHAR(255);    
    DECLARE member_exists INT DEFAULT 0;

    -- 測試階段用
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;
    END;

    -- SET error_message = CONCAT('[ERR]ERROR MESSAGE');
    -- SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;

    SELECT COUNT(`Account`) INTO member_exists FROM `member` WHERE `Account` = p_account;
    SELECT member_exists as `Exist`;
END$$

DELIMITER ;

CALL `sp_test`('chris') ;