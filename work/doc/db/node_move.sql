DROP PROCEDURE IF EXISTS `sp_node_move` ;

DELIMITER $$
CREATE PROCEDURE IF NOT EXISTS `sp_node_move`(
  IN p_account VARCHAR(50),
  IN p_parent_id CHAR(36),
  IN p_node_id CHAR(36),
  IN p_path_name VARCHAR(255)
)
BEGIN

  CALL sp_node_del(p_account, p_node_id);

  IF p_parent_id = '00000000-0000-0000-0000-000000000000' THEN
    CALL sp_add_parent_node(p_account, p_node_id, p_path_name);
  ELSE
    CALL sp_add_child_node(p_account, p_parent_id, p_node_id, p_path_name);
  END IF;

END $$
DELIMITER ;
