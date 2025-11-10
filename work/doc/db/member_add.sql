
DROP PROCEDURE IF EXISTS `sp_add_member` ;

DELIMITER $$

CREATE PROCEDURE IF NOT EXISTS `sp_add_member` (
  IN p_account VARCHAR(50),
  IN p_password VARCHAR(100),
  IN p_ip VARCHAR(45)
)
BEGIN
  -- 所有 DECLARE 必須在最前面
  DECLARE error_message VARCHAR(255);
  DECLARE debug_message VARCHAR(255);
  DECLARE table_articles TEXT;
  DECLARE table_node TEXT;

  DECLARE EXIT HANDLER FOR SQLEXCEPTION
  BEGIN
    ROLLBACK;
    -- SET error_message = CONCAT('註冊失敗，帳號 "', p_account, '" 發生錯誤於：', error_message);
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;
  END;

  START TRANSACTION;

  SET table_articles = CONCAT('articles_', p_account);
  SET table_node = CONCAT('node_', p_account);

  -- 帳號驗證
  SET error_message = '帳號驗證';
  IF EXISTS (SELECT 1 FROM `member` WHERE `Account` = p_account) THEN
    SET error_message = CONCAT('帳號 "', p_account, '" 已存在',  '[ERR]帳號已存在，請使用其他帳號');
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;
  END IF;

  -- 加入資料
  SET error_message = '加入資料';
  INSERT INTO `member` ( `Account`, `Password`, `LastIP`, `CreatedAt`, `UpdatedAt`) 
  VALUES ( p_account, p_password, p_ip, UTC_TIMESTAMP(), UTC_TIMESTAMP() );

  -- 建立使用者文章表單
  SET error_message = '建立使用者文章表單';
  SET @query = CONCAT(
    'CREATE TABLE `', table_articles, '` (',
    '  `RowID` INT NOT NULL AUTO_INCREMENT,',
    '  `Title` TEXT,',
    '  `Content` LONGTEXT,',
    '  `NodeID` VARCHAR(36),',
    '  `CreatedDt` DATETIME,',
    '  `UpdateDt` DATETIME,',
    '  PRIMARY KEY (`RowID`)',
    ') ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;'
  );
  PREPARE stmt FROM @query;
  EXECUTE stmt;
  DEALLOCATE PREPARE stmt;

  -- 建立使用者節點表單
  SET error_message = '建立使用者節點表單';
  SET @query = CONCAT(
    'CREATE TABLE `', table_node, '` (',
    '  `RowID` INT NOT NULL AUTO_INCREMENT,',
    '  `NodeID` VARCHAR(36),',
    '  `ParentID` VARCHAR(36),',
    '  `PathName` VARCHAR(255),',
    '  `LftIdx` INT,',
    '  `RftIdx` INT,',
    '  PRIMARY KEY (`RowID`)',
    ') ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;'
  );
  PREPARE stmt FROM @query;
  EXECUTE stmt;
  DEALLOCATE PREPARE stmt;

  COMMIT;

  SELECT * from `member` where `Account` = p_account;
END$$

DELIMITER ;
