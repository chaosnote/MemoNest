DROP TABLE IF EXISTS `member` ;

CREATE TABLE IF NOT EXISTS `member` ( 
    `RowID` INT NOT NULL AUTO_INCREMENT PRIMARY KEY COMMENT '列 ID / 內部 ID', 
    `Account` VARCHAR(20) NOT NULL UNIQUE COMMENT '帳號', 
    `Password` VARCHAR(32) NOT NULL COMMENT '密碼的雜湊值', 
    `LastIP` VARCHAR(45) NULL COMMENT '最後登入 IP 地址',
    `IsEnabled` TINYINT(1) NOT NULL DEFAULT 1 COMMENT '帳號是否啟用 1=啟用, 0=停用',
    `CreatedAt` DATETIME NOT NULL COMMENT '建立時間', 
    `UpdatedAt` DATETIME NOT NULL COMMENT '更新時間', 
    INDEX idx_account (`Account`) 
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '會員資訊表';
