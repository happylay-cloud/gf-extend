package gfadapter

import (
	"fmt"
)

// CreateSqlite3Table 创建Sqlite3表结构sql
func CreateSqlite3Table(tableName string) (sql string) {
	return fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s(
    p_type VARCHAR(32)  DEFAULT '' NOT NULL,
    v0     VARCHAR(255) DEFAULT '' NOT NULL,
    v1     VARCHAR(255) DEFAULT '' NOT NULL,
    v2     VARCHAR(255) DEFAULT '' NOT NULL,
    v3     VARCHAR(255) DEFAULT '' NOT NULL,
    v4     VARCHAR(255) DEFAULT '' NOT NULL,
    v5     VARCHAR(255) DEFAULT '' NOT NULL,
    CHECK (TYPEOF("p_type") = "text" AND
           LENGTH("p_type") <= 32),
    CHECK (TYPEOF("v0") = "text" AND
           LENGTH("v0") <= 255),
    CHECK (TYPEOF("v1") = "text" AND
           LENGTH("v1") <= 255),
    CHECK (TYPEOF("v2") = "text" AND
           LENGTH("v2") <= 255),
    CHECK (TYPEOF("v3") = "text" AND
           LENGTH("v3") <= 255),
    CHECK (TYPEOF("v4") = "text" AND
           LENGTH("v4") <= 255),
    CHECK (TYPEOF("v5") = "text" AND
           LENGTH("v5") <= 255)
);
CREATE INDEX IF NOT EXISTS %s ON %s (p_type, v0, v1);`,
		tableName, "idx_"+tableName, tableName)
}

// CreateMysqlTable 创建Mysql表结构sql
func CreateMysqlTable(tableName string) (sql string) {
	return fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s (
    p_type VARCHAR(32)  NOT NULL DEFAULT '' COMMENT '类型，p代表策略，g代表角色',
    v0     VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'sub',
    v1     VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'obj',
    v2     VARCHAR(255) NOT NULL DEFAULT '' COMMENT 'act',
    v3     VARCHAR(255) NOT NULL DEFAULT '',
    v4     VARCHAR(255) NOT NULL DEFAULT '',
    v5     VARCHAR(255) NOT NULL DEFAULT '',
    INDEX %s (p_type, v0, v1)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;`,
		tableName, "idx_"+tableName)
}

// CreatePgsqlTable 创建Pgsql表结构sql
func CreatePgsqlTable(tableName string) (sql string) {
	return fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s(
    p_type VARCHAR(32)  DEFAULT '' NOT NULL,
    v0     VARCHAR(255) DEFAULT '' NOT NULL,
    v1     VARCHAR(255) DEFAULT '' NOT NULL,
    v2     VARCHAR(255) DEFAULT '' NOT NULL,
    v3     VARCHAR(255) DEFAULT '' NOT NULL,
    v4     VARCHAR(255) DEFAULT '' NOT NULL,
    v5     VARCHAR(255) DEFAULT '' NOT NULL
);
CREATE INDEX IF NOT EXISTS %s ON %s (p_type, v0, v1);`,
		tableName, "idx_"+tableName, tableName)
}

// DropTable 删除指定表
func DropTable(tableName string) (sql string) {
	return fmt.Sprintf(`DROP TABLE IF EXISTS %s ;`, tableName)
}
