package mysql

import (
	"database/sql"
	"fmt"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/ysicing/ergo/pkg/util/log"
)

type DB struct {
	Dsn    string
	client *sql.DB
}

func NewDB(dsn string) (*DB, error) {
	var db DB
	db.Dsn = dsn
	dbclient, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Flog.Errorf("create sql client err: %v", err)
		return nil, err
	}
	dbclient.SetConnMaxLifetime(time.Minute * 3)
	dbclient.SetMaxOpenConns(10)
	dbclient.SetMaxIdleConns(10)
	db.client = dbclient
	return &db, nil
}

func (db *DB) Create(dbname, dbuser, dbpass string) error {
	_, err := db.client.Exec("CREATE DATABASE IF NOT EXISTS " + dbname + ";")
	if err != nil {
		log.Flog.Errorf("create db %v err: %v", dbname, err)
		return fmt.Errorf("创建数据库失败")
	}
	_, err = db.client.Exec("use " + dbname)
	if err != nil {
		log.Flog.Errorf("use db %v err: %v", dbname, err)
		return fmt.Errorf("创建数据库失败")
	}
	_, err = db.client.Exec("CREATE USER '" + dbuser + "'@'%' IDENTIFIED BY '" + dbpass + "';")
	if err != nil {
		log.Flog.Errorf("create user %v err: %v", dbuser, err)
		return fmt.Errorf("创建用户失败")
	}
	grantCmd := fmt.Sprintf("GRANT ALL ON %s.* TO '%s'@'%%'", dbname, dbuser)
	_, err = db.client.Exec(grantCmd)
	if err != nil {
		log.Flog.Errorf("grant user %v err: %v", dbuser, err)
		return fmt.Errorf("授权失败")
	}
	_, err = db.client.Exec("flush privileges;")
	if err != nil {
		log.Flog.Errorf("flush privileges err: %v", err)
		return fmt.Errorf("刷新权限失败")
	}
	return nil
}

func (db *DB) Exec(sql string) error {
	res, err := db.client.Exec(sql)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Flog.Debugf("exec sql: %s, affected: %v", sql, affected)
	return nil
}

func (db *DB) Drop(dbname, dbuser string) error {
	if len(dbuser) != 0 {
		// 移除权限
		revokeCmd := fmt.Sprintf("REVOKE ALL ON %s.* FROM '%s'@'%%';", dbname, dbuser)
		_, err := db.client.Exec(revokeCmd)
		if err != nil {
			log.Flog.Warnf("revoke user %v err: %v, sql: %v", dbuser, err, revokeCmd)
		}
		log.Flog.Debugf("revoke user %v", dbuser)
		// 删除用户
		dropUserCmd := fmt.Sprintf("DROP USER IF EXISTS \"%v\";", dbuser)
		_, err = db.client.Exec(dropUserCmd)
		if err != nil {
			log.Flog.Errorf("delete user %v err: %v, sql: %v", dbuser, err, dropUserCmd)
			return err
		}
		log.Flog.Debugf("delete user %v", dbuser)
	}
	// 删除数据库
	dropDBCmd := fmt.Sprintf("DROP DATABASE IF EXISTS %v;", dbname)
	_, err := db.client.Exec(dropDBCmd)
	if err != nil {
		log.Flog.Errorf("delete db %v err: %v, sql: %v", dbname, err, dropDBCmd)
		return err
	}
	log.Flog.Debugf("delete db %v", dbname)
	_, err = db.client.Exec("flush privileges;")
	if err != nil {
		log.Flog.Errorf("flush privileges err: %v", err)
		return err
	}
	log.Flog.Debugf("刷新权限")
	return nil
}

func (db *DB) Ping() error {
	return db.client.Ping()
}

type DBCfg struct {
	Name string `json:"name"`
}

func (db *DB) Show() ([]DBCfg, error) {
	res, err := db.client.Query("SELECT schema_name as `database` FROM information_schema.schemata;")
	if err != nil {
		log.Flog.Errorf("query db err: %v", err)
		return nil, err
	}
	dbs := make([]DBCfg, 0)
	for res.Next() {
		var dbname string
		if err := res.Scan(&dbname); err != nil {
			log.Flog.Errorf("scan err: ", err)
			continue
		}
		dbs = append(dbs, DBCfg{
			Name: dbname,
		})
	}
	return dbs, nil
}

func (db *DB) Close() {
	if err := db.client.Close(); err != nil {
		log.Flog.Errorf("close db err: %v", err)
	}
}