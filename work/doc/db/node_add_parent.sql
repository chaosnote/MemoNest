DROP PROCEDURE IF EXISTS `sp_node_add_parent` ;

DELIMITER $$
CREATE PROCEDURE IF NOT EXISTS sp_node_add_parent(
  IN p_account VARCHAR(50),
  IN p_node_id CHAR(36),
  IN p_path_name VARCHAR(255)
)
BEGIN

  DECLARE max_rft INT DEFAULT 0;
  DECLARE table_node VARCHAR(64);

  SET table_node = CONCAT('node_', p_account);

  SET @sql = CONCAT('SELECT COALESCE(MAX(RftIdx), 0) INTO @max_rft FROM ', table_node);
  PREPARE stmt FROM @sql;
  EXECUTE stmt;
  DEALLOCATE PREPARE stmt;

  IF p_node_id = '' THEN
    SET p_node_id = UUID();
  END IF;

  SET @lft = @max_rft + 1;
  SET @rft = @max_rft + 2;

  SET @sql = CONCAT('INSERT INTO ', table_node, ' (NodeID, ParentID, PathName, LftIdx, RftIdx) VALUES (?, UUID(), ?, ?, ?)');
  PREPARE stmt FROM @sql;
  EXECUTE stmt USING p_node_id, p_path_name, @lft, @rft;
  DEALLOCATE PREPARE stmt;

  SET @sql = CONCAT('SELECT `RowID` FROM `',table_node,'` WHERE `NodeID` = ? ') ;
  PREPARE stmt FROM @sql;
  EXECUTE stmt USING p_node_id;
  DEALLOCATE PREPARE stmt;

END $$
DELIMITER ;