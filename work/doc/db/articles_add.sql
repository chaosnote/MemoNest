
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
    DECLARE debug_message VARCHAR(255);
    DECLARE table_articles TEXT;
    DECLARE table_categories TEXT;

    -- 測試階段用
    -- DECLARE EXIT HANDLER FOR SQLEXCEPTION
    -- BEGIN
    --     SET error_message = CONCAT('新增文章失敗，帳號 "', p_account, '" 發生錯誤於：', debug_message);
    --     SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;
    -- END;

    SET table_articles = CONCAT('articles_', p_account);
    SET table_categories = CONCAT('node_', p_account);

    -- 驗證是否有使用者文章表單
    SET debug_message = '驗證是否有使用者文章表單';
    SET @table_exists = 0 ;
    SET @query = CONCAT(
        'SELECT COUNT(*) INTO @table_exists FROM INFORMATION_SCHEMA.TABLES ',
        'WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = "', table_articles, '"'
    );
    PREPARE stmt FROM @query;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    IF @table_exists = 0 THEN
        SET error_message = CONCAT('資料表 `', table_articles, '` 不存在[ERR]無指定資料');
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;
    END IF;

    -- 驗證是否有使用者節點表單
    SET debug_message = '驗證是否有使用者節點表單';
    SET @table_exists = 0 ;
    SET @query = CONCAT(
        'SELECT COUNT(*) INTO @table_exists FROM INFORMATION_SCHEMA.TABLES ',
        'WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = "', table_categories, '"'
    );
    PREPARE stmt FROM @query;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    IF @table_exists = 0 THEN
        SET error_message = CONCAT('資料表 `', table_articles, '` 不存在[ERR]無指定資料');
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;
    END IF;

    -- 驗證是否有指定的節點
    SET debug_message = '驗證是否有指定的節點';
    SET @node_exists = 0;
    SET @query = CONCAT(
        'SELECT 1 INTO @node_exists FROM `', table_categories, '` WHERE `NodeID` = ?'
    );
    PREPARE stmt FROM @query;    
    EXECUTE stmt USING p_NodeID;
    DEALLOCATE PREPARE stmt;

    IF @node_exists = 0 THEN
        SET error_message = CONCAT('NodeID(', p_NodeID, ') 不存在於 `', table_categories, '`[ERR]無指定節點');
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = error_message;
    END IF;

    -- 插入
    SET debug_message = '插入文章';
    SET @query = CONCAT(
        'INSERT INTO `', table_articles, '` (`Title`, `Content`, `UpdateDt`, `CreatedDt`, `NodeID`) VALUES (',
        QUOTE(p_Title), ', ',
        QUOTE(p_Content), ', ',
        QUOTE(p_UpdateDt), ', ',
        QUOTE(p_CreatedDt), ', ',
        QUOTE(p_NodeID), ')'
    );
    PREPARE stmt FROM @query;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    -- 查詢總筆數
    SET debug_message = '查詢總筆數';
    SET @query = CONCAT('SELECT COUNT(`RowID`) FROM ', table_articles);
    PREPARE stmt FROM @query;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END$$

DELIMITER ;
