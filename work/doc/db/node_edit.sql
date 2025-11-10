DROP PROCEDURE IF EXISTS `sp_node_edit` ;

DELIMITER $$
CREATE PROCEDURE IF NOT EXISTS `sp_node_edit`(
  IN p_account VARCHAR(50),
  IN p_node_id CHAR(36),
  IN p_path_name VARCHAR(255)
)
BEGIN

  DECLARE table_node VARCHAR(64);

  SET table_node = CONCAT('node_', p_account);

  SET @sql = CONCAT('UPDATE `',table_node,'` SET PathName = ? WHERE NodeID = ?;');
  PREPARE stmt FROM @sql;
  EXECUTE stmt USING p_path_name, p_node_id;
  DEALLOCATE PREPARE stmt;

END $$
DELIMITER ;