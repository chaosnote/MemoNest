
DROP PROCEDURE IF EXISTS `sp_login`;

DELIMITER $$

CREATE PROCEDURE `sp_login` (
  IN p_account VARCHAR(50),
  IN p_password VARCHAR(100),
  IN p_ip VARCHAR(45)
)
BEGIN
  DECLARE error_message TEXT;
  DECLARE member_exists INT DEFAULT 0;

  DECLARE EXIT HANDLER FOR SQLEXCEPTION
  BEGIN
    SET error_message = CONCAT('登入失敗，帳號 "', p_account, '" 發生錯誤於：', @debug_step);
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;
  END;

  SET @debug_step = 'check member';
  SELECT COUNT(*) INTO member_exists
  FROM `member`
  WHERE `Account` = p_account AND `Password` = p_password;

  IF member_exists = 0 THEN
    SET error_message = CONCAT('帳號或密碼錯誤：', p_account, '[ERR]帳號或密碼錯誤');
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;
  END IF;

  SET @debug_step = 'update last_ip';
  UPDATE `member`
  SET `LastIP` = p_ip,
      `UpdatedAt` = UTC_TIMESTAMP()
  WHERE `Account` = p_account;

  SET @debug_step = 'return member';
  SELECT * FROM `member` WHERE `Account` = p_account;
END$$

DELIMITER ;