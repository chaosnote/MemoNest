DROP PROCEDURE IF EXISTS `sp_node_add_child` ;

DELIMITER $$
CREATE PROCEDURE IF NOT EXISTS `sp_node_add_child`(
  IN p_account VARCHAR(50),
  IN p_parent_id CHAR(36),
  IN p_node_id CHAR(36),
  IN p_path_name VARCHAR(255)
)
BEGIN

  DECLARE parent_rft INT;
  DECLARE table_node VARCHAR(64);

  SET table_node = CONCAT('node_', p_account);

  SET @sql = CONCAT('SELECT RftIdx INTO @parent_rft FROM ', table_node, ' WHERE NodeID = ?');
  PREPARE stmt FROM @sql;
  EXECUTE stmt USING p_parent_id;
  DEALLOCATE PREPARE stmt;

  IF @parent_rft IS NULL THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '[ERR]未找到父節點';
  END IF;

  SET @sql = CONCAT('UPDATE ', table_node, ' SET RftIdx = RftIdx + 2 WHERE RftIdx >= ?');
  PREPARE stmt FROM @sql;
  EXECUTE stmt USING @parent_rft;
  DEALLOCATE PREPARE stmt;

  SET @sql = CONCAT('UPDATE ', table_node, ' SET LftIdx = LftIdx + 2 WHERE LftIdx >= ?');
  PREPARE stmt FROM @sql;
  EXECUTE stmt USING @parent_rft;
  DEALLOCATE PREPARE stmt;

  IF p_node_id = '' THEN
    SET p_node_id = UUID();
  END IF;

  SET @lft = @parent_rft;
  SET @rft = @parent_rft + 1;

  SET @sql = CONCAT('INSERT INTO ', table_node, ' (NodeID, ParentID, PathName, LftIdx, RftIdx) VALUES (?, ?, ?, ?, ?)');
  PREPARE stmt FROM @sql;
  EXECUTE stmt USING p_node_id, p_parent_id, p_path_name, @lft, @rft;
  DEALLOCATE PREPARE stmt;

  SET @sql = CONCAT('SELECT * FROM `',table_node,'` WHERE `NodeID` = ? ') ;
  PREPARE stmt FROM @sql;
  EXECUTE stmt USING p_node_id;
  DEALLOCATE PREPARE stmt;

END $$
DELIMITER ;