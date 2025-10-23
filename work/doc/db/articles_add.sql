-- SP
DROP PROCEDURE IF EXISTS `sp_add_article`;

DELIMITER $$

CREATE PROCEDURE IF NOT EXISTS `sp_add_article` (
    IN p_account VARCHAR(50),
    IN p_Title VARCHAR(255),
    IN p_Content TEXT,
    IN p_UpdateDt DATETIME,
    IN p_CreatedDt DATETIME,
    IN p_NodeID CHAR(36)
)
BEGIN
    DECLARE error_message VARCHAR(255);
    DECLARE table_articles TEXT;
    DECLARE table_categories TEXT;
    DECLARE node_exists INT DEFAULT 0;

    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        SET error_message = CONCAT('新增文章失敗，帳號 "', p_account, '" 發生錯誤於：', @debug_step);
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;
    END;

    SET @debug_step = 'prepare table names';
    SET table_articles = CONCAT('articles_', p_account);
    SET table_categories = CONCAT('categories_', p_account);

    SET @debug_step = 'check node exists';
    SET @sql_check_node = CONCAT(
        'SELECT COUNT(*) INTO @node_exists FROM `', table_categories, '` WHERE `NodeID` = ?'
    );
    PREPARE stmt_check_node FROM @sql_check_node;
    SET @node_exists = 0;
    EXECUTE stmt_check_node USING p_NodeID;
    DEALLOCATE PREPARE stmt_check_node;

    IF @node_exists = 0 THEN
        SET error_message = CONCAT('NodeID(', p_NodeID, ') 不存在於 `', table_categories, '`[ERR]無指定節點');
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;
    END IF;

    SET @debug_step = 'insert article';
    SET @sql_insert_article = CONCAT(
        'INSERT INTO `', table_articles, '` (`Title`, `Content`, `UpdateDt`, `CreatedDt`, `NodeID`) VALUES (',
        QUOTE(p_Title), ', ',
        QUOTE(p_Content), ', ',
        QUOTE(p_UpdateDt), ', ',
        QUOTE(p_CreatedDt), ', ',
        QUOTE(p_NodeID), ')'
    );

    -- SELECT @sql_insert_article AS debug_sql;
    
    PREPARE stmt_insert_article FROM @sql_insert_article;
    EXECUTE stmt_insert_article;
    DEALLOCATE PREPARE stmt_insert_article;

    SET @debug_step = 'return rowid';
    SELECT LAST_INSERT_ID() AS RowID;
END$$

DELIMITER ;

-- CALL sp_add_article('AAA', '123', NOW(), NOW(), '002198c1-1efe-4111-8c1c-d74d293b5823') ;