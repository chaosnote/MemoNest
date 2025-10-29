
DROP PROCEDURE IF EXISTS `sp_login`;

DELIMITER $$

CREATE PROCEDURE `sp_login` (
  IN p_account VARCHAR(50),
  IN p_password VARCHAR(100),
  IN p_ip VARCHAR(45)
)
BEGIN
  DECLARE error_message VARCHAR(255);
  DECLARE debug_message VARCHAR(255);
  DECLARE member_exists INT DEFAULT 0;

  -- DECLARE EXIT HANDLER FOR SQLEXCEPTION
  -- BEGIN
  --   SET error_message = CONCAT('登入失敗，帳號 "', p_account, '" 發生錯誤於：', debug_message);
  --   SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;
  -- END;

  SET debug_message = '驗證使用者是否存在';
  SELECT COUNT('Account') INTO member_exists FROM `member` WHERE `Account` = p_account AND `Password` = p_password;

  IF member_exists = 0 THEN
    SET error_message = CONCAT('帳號 :', p_account,',密碼 :', p_password, '[ERR]帳號或密碼錯誤');
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;
  END IF;

  UPDATE `member` SET `LastIP` = p_ip, `UpdatedAt` = UTC_TIMESTAMP() WHERE `Account` = p_account;

  SELECT * FROM `member` WHERE `Account` = p_account;
END$$

DELIMITER ;