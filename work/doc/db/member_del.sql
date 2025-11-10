
DROP PROCEDURE IF EXISTS `sp_del_member`;

DELIMITER $$

CREATE PROCEDURE IF NOT EXISTS `sp_del_member` (
  IN p_account VARCHAR(50)
)
BEGIN
  DECLARE error_message VARCHAR(255);
  DECLARE debug_message VARCHAR(255);
  DECLARE table_articles TEXT;
  DECLARE table_node TEXT;

  DECLARE EXIT HANDLER FOR SQLEXCEPTION
  BEGIN
    ROLLBACK;
    -- SET error_message = CONCAT('刪除失敗，帳號 "', p_account, '" 發生錯誤於：', debug_message);
    -- SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;
  END;

  START TRANSACTION;

  SET table_articles = CONCAT('articles_', p_account);
  SET table_node = CONCAT('node_', p_account);

  -- DROP articles 表
  SET debug_message = '刪除使用者文章表單';
  SET @query = CONCAT('DROP TABLE IF EXISTS `', table_articles, '`;');
  PREPARE stmt FROM @query;
  EXECUTE stmt;
  DEALLOCATE PREPARE stmt;

  -- DROP categories 表
  SET debug_message = '刪除使用者節點表單';
  SET @query = CONCAT('DROP TABLE IF EXISTS `', table_node, '`;');
  PREPARE stmt FROM @query;
  EXECUTE stmt;
  DEALLOCATE PREPARE stmt;

  -- DELETE member row
  SET debug_message = 'delete member';
  DELETE FROM `member` WHERE `Account` = p_account;

  COMMIT;
END$$

DELIMITER ;
