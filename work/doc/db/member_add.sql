
DROP PROCEDURE IF EXISTS `sp_add_member` ;

DELIMITER $$

CREATE PROCEDURE IF NOT EXISTS `sp_add_member` (
  IN p_account VARCHAR(50),
  IN p_password VARCHAR(100),
  IN p_ip VARCHAR(45)
)
BEGIN
  -- 所有 DECLARE 必須在最前面
  DECLARE error_message TEXT;
  DECLARE table_articles TEXT;
  DECLARE table_categories TEXT;

  DECLARE EXIT HANDLER FOR SQLEXCEPTION
  BEGIN
    ROLLBACK;
    SET error_message = CONCAT('註冊失敗，帳號 "', p_account, '" 發生錯誤於：', @debug_step);
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;
  END;

  -- 初始化 debug 變數
  SET @debug_step = 'init';

  START TRANSACTION;

  SET @debug_step = 'check account exists';
  IF EXISTS (SELECT 1 FROM `member` WHERE `Account` = p_account) THEN
    SET error_message = CONCAT('帳號 "', p_account, '" 已存在，請使用其他帳號',  '[ERR]帳號已存在');
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;
  END IF;

  SET @debug_step = 'insert member';
  INSERT INTO `member` (
    `Account`, `Password`, `LastIP`, `CreatedAt`, `UpdatedAt`
  ) VALUES (
    p_account, p_password, p_ip, UTC_TIMESTAMP(), UTC_TIMESTAMP()
  );

  SET @debug_step = 'prepare table names';
  SET table_articles = CONCAT('articles_', p_account);
  SET table_categories = CONCAT('categories_', p_account);

  -- DROP articles 表（防呆）
  SET @debug_step = 'drop articles table';
  SET @sql_drop_articles = CONCAT('DROP TABLE IF EXISTS `', table_articles, '`;');
  PREPARE stmt_drop_articles FROM @sql_drop_articles;
  EXECUTE stmt_drop_articles;
  DEALLOCATE PREPARE stmt_drop_articles;

  -- 建立 articles 表
  SET @debug_step = 'create articles table';
  SET @sql_articles = CONCAT(
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
  PREPARE stmt_articles FROM @sql_articles;
  EXECUTE stmt_articles;
  DEALLOCATE PREPARE stmt_articles;

  -- DROP categories 表（防呆）
  SET @debug_step = 'drop categories table';
  SET @sql_drop_categories = CONCAT('DROP TABLE IF EXISTS `', table_categories, '`;');
  PREPARE stmt_drop_categories FROM @sql_drop_categories;
  EXECUTE stmt_drop_categories;
  DEALLOCATE PREPARE stmt_drop_categories;

  -- 建立 categories 表
  SET @debug_step = 'create categories table';
  SET @sql_categories = CONCAT(
    'CREATE TABLE `', table_categories, '` (',
    '  `RowID` INT NOT NULL AUTO_INCREMENT,',
    '  `NodeID` VARCHAR(36),',
    '  `ParentID` VARCHAR(36),',
    '  `PathName` VARCHAR(255),',
    '  `LftIdx` INT,',
    '  `RftIdx` INT,',
    '  PRIMARY KEY (`RowID`)',
    ') ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;'
  );
  PREPARE stmt_categories FROM @sql_categories;
  EXECUTE stmt_categories;
  DEALLOCATE PREPARE stmt_categories;

  SET @debug_step = 'commit';
  COMMIT;

  SELECT * from `member` where `Account` = p_account;
END$$

DELIMITER ;
