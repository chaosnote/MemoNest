DROP PROCEDURE IF EXISTS `sp_node_del` ;

DELIMITER $$
CREATE PROCEDURE IF NOT EXISTS `sp_node_del`(
  IN p_account VARCHAR(50),
  IN p_node_id CHAR(36)
)
SP:BEGIN

  DECLARE lft INT;
  DECLARE rft INT;
  DECLARE width INT;
  DECLARE table_node VARCHAR(64);
  DECLARE table_articles VARCHAR(64);

  SET table_node = CONCAT('node_', p_account);
  SET table_articles = CONCAT('articles_', p_account);

  SET @sql = CONCAT('SELECT COUNT(*) INTO @Total FROM `', table_articles,'` WHERE NodeID = ?');
  PREPARE stmt FROM @sql;
  EXECUTE stmt USING p_node_id;
  DEALLOCATE PREPARE stmt;

  IF @Total > 0 THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '[ERR]該節點有關聯文章、請先刪除文章或轉移文章節點';
  END IF;

  SET @sql = CONCAT('SELECT LftIdx, RftIdx INTO @lft, @rft FROM ', table_node, ' WHERE NodeID = ?');
  PREPARE stmt FROM @sql;
  EXECUTE stmt USING p_node_id;
  DEALLOCATE PREPARE stmt;

  -- 節點不存在，視為成功
  IF @lft IS NULL OR @rft IS NULL THEN
    LEAVE SP;
  END IF;

  SET @width = @rft - @lft + 1;

  SET @sql = CONCAT('DELETE FROM ', table_node, ' WHERE LftIdx >= ? AND RftIdx <= ?');
  PREPARE stmt FROM @sql;
  EXECUTE stmt USING @lft, @rft;
  DEALLOCATE PREPARE stmt;

  SET @sql = CONCAT('UPDATE ', table_node, ' SET RftIdx = RftIdx - ? WHERE RftIdx > ?');
  PREPARE stmt FROM @sql;
  EXECUTE stmt USING @width, @rft;
  DEALLOCATE PREPARE stmt;

  SET @sql = CONCAT('UPDATE ', table_node, ' SET LftIdx = LftIdx - ? WHERE LftIdx > ?');
  PREPARE stmt FROM @sql;
  EXECUTE stmt USING @width, @rft;
  DEALLOCATE PREPARE stmt;

END $$
DELIMITER ;
