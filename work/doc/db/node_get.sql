DROP PROCEDURE IF EXISTS `sp_node_get` ;

DELIMITER $$
CREATE PROCEDURE IF NOT EXISTS `sp_node_get`(
  IN p_account VARCHAR(50),
  IN p_node_id CHAR(36)
)
BEGIN

  DECLARE table_node VARCHAR(64);

  SET table_node = CONCAT('node_', p_account);

  SET @sql = CONCAT('SELECT * FROM `', table_node,'` WHERE `NodeID` = ? ORDER BY LftIdx ASC');
  PREPARE stmt FROM @sql;
  EXECUTE stmt USING p_node_id;
  DEALLOCATE PREPARE stmt;

END $$
DELIMITER ;

--

DROP PROCEDURE IF EXISTS `sp_node_list` ;

DELIMITER $$
CREATE PROCEDURE IF NOT EXISTS `sp_node_list`(
  IN p_account VARCHAR(50)
)
BEGIN

  DECLARE table_node VARCHAR(64);

  SET table_node = CONCAT('node_', p_account);

  SET @sql = CONCAT('SELECT * FROM `', table_node,'` ORDER BY LftIdx ASC');
  PREPARE stmt FROM @sql;
  EXECUTE stmt ;
  DEALLOCATE PREPARE stmt;

END $$
DELIMITER ;