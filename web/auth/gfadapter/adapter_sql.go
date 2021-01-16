package gfadapter

import "github.com/gogf/gf/database/gdb"

// CreateSqlite3Table 创建Sqlite3表结构sql
func CreateSqlite3Table(tableName string) (sql string) {
	return gdb.FormatSqlWithArgs(`
CREATE TABLE IF NOT EXISTS ?(
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
CREATE INDEX IF NOT EXISTS ? ON ? (p_type, v0, v1);`,
		[]interface{}{tableName, "idx_" + tableName, tableName})

}

// CreateMysqlTable 创建Mysql表结构sql
func CreateMysqlTable(tableName string) (sql string) {
	return gdb.FormatSqlWithArgs(`
CREATE TABLE IF NOT EXISTS ?(
    p_type VARCHAR(32)  DEFAULT '' NOT NULL,
    v0     VARCHAR(255) DEFAULT '' NOT NULL,
    v1     VARCHAR(255) DEFAULT '' NOT NULL,
    v2     VARCHAR(255) DEFAULT '' NOT NULL,
    v3     VARCHAR(255) DEFAULT '' NOT NULL,
    v4     VARCHAR(255) DEFAULT '' NOT NULL,
    v5     VARCHAR(255) DEFAULT '' NOT NULL,
    INDEX idx_? (p_type, v0, v1)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;`,
		[]interface{}{tableName, "idx_" + tableName})
}

// CreatePgsqlTable 创建Pgsql表结构sql
func CreatePgsqlTable(tableName string) (sql string) {
	return gdb.FormatSqlWithArgs(`
CREATE TABLE IF NOT EXISTS ?(
    p_type VARCHAR(32)  DEFAULT '' NOT NULL,
    v0     VARCHAR(255) DEFAULT '' NOT NULL,
    v1     VARCHAR(255) DEFAULT '' NOT NULL,
    v2     VARCHAR(255) DEFAULT '' NOT NULL,
    v3     VARCHAR(255) DEFAULT '' NOT NULL,
    v4     VARCHAR(255) DEFAULT '' NOT NULL,
    v5     VARCHAR(255) DEFAULT '' NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_? ON ? (p_type, v0, v1);`,
		[]interface{}{tableName, "idx_" + tableName})

}

// DropTable 删除指定表
func DropTable(tableName string) (sql string) {
	return gdb.FormatSqlWithArgs("DROP TABLE IF EXISTS ?;",
		[]interface{}{tableName})
}
