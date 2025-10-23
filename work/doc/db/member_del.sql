
DROP PROCEDURE IF EXISTS `sp_del_member`;

DELIMITER $$

CREATE PROCEDURE IF NOT EXISTS `sp_del_member` (
  IN p_account VARCHAR(50)
)
BEGIN
  DECLARE error_message TEXT;
  DECLARE table_articles TEXT;
  DECLARE table_categories TEXT;

  DECLARE EXIT HANDLER FOR SQLEXCEPTION
  BEGIN
    ROLLBACK;
    -- SET error_message = CONCAT('刪除失敗，帳號 "', p_account, '" 發生錯誤於：', @debug_step);
    -- SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;
  END;

  SET @debug_step = 'init';
  START TRANSACTION;

  SET @debug_step = 'prepare table names';
  SET table_articles = CONCAT('articles_', p_account);
  SET table_categories = CONCAT('categories_', p_account);

  -- DROP articles 表
  SET @debug_step = 'drop articles table';
  SET @sql_drop_articles = CONCAT('DROP TABLE IF EXISTS `', table_articles, '`;');
  PREPARE stmt_drop_articles FROM @sql_drop_articles;
  EXECUTE stmt_drop_articles;
  DEALLOCATE PREPARE stmt_drop_articles;

  -- DROP categories 表
  SET @debug_step = 'drop categories table';
  SET @sql_drop_categories = CONCAT('DROP TABLE IF EXISTS `', table_categories, '`;');
  PREPARE stmt_drop_categories FROM @sql_drop_categories;
  EXECUTE stmt_drop_categories;
  DEALLOCATE PREPARE stmt_drop_categories;

  -- DELETE member row
  SET @debug_step = 'delete member';
  DELETE FROM `member` WHERE `Account` = p_account;

  SET @debug_step = 'commit';
  COMMIT;
END$$

DELIMITER ;
