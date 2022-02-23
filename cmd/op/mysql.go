package op

import (
	"fmt"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/environ"
	"github.com/ergoapi/util/exhash"
	"github.com/ergoapi/util/ztime"
	"github.com/spf13/cobra"
	"github.com/ysicing/ergo/pkg/ergo/ops/mysql"
)

type mysqlOption struct {
	Dsn  string
	Host string
	Port int
	User string
	Pass string
	Name string
}

func (opt *mysqlOption) Check() {
	opt.Dsn = environ.GetEnv("MYSQL_DSN", opt.Dsn)
	opt.Host = environ.GetEnv("MYSQL_HOST", opt.Host)
	opt.Port = environ.GetEnvAsInt("MYSQL_PORT", opt.Port)
	opt.User = environ.GetEnv("MYSQL_USER", opt.User)
	opt.Pass = environ.GetEnv("MYSQL_PASS", opt.Pass)
	opt.Name = environ.GetEnv("MYSQL_NAME", opt.Name)
	if len(opt.Dsn) == 0 {
		opt.Dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", opt.User, opt.Pass, opt.Host, opt.Port, opt.Name)
	}
}

func MysqlCmd() *cobra.Command {
	opt := &mysqlOption{}
	mysqlcmd := &cobra.Command{
		Use:     "mysql",
		Short:   "ping",
		Version: "3.1.0",
		Args:    cobra.MinimumNArgs(1),
		Run: func(cobraCmd *cobra.Command, args []string) {
			zlog := log.GetInstance()
			opt.Check()
			db, err := mysql.NewDB(opt.Dsn)
			if err != nil {
				zlog.Errorf("create db client err: %v", err)
				return
			}
			defer db.Close()
			zlog.Debugf("mysql dsn: %s", opt.Dsn)
			switch args[0] {
			case "ping":
				if err := db.Ping(); err != nil {
					zlog.Errorf("db ping err: %v", err)
					return
				}
				zlog.Info("db ping done")
			case "exec":
				if err := db.Exec(args[1]); err != nil {
					zlog.Errorf("db exec err: %v", err)
					return
				}
				zlog.Info("db exec done")
			case "new":
				user := exhash.MD5(ztime.NowFormat())
				pass := exhash.MD5(ztime.NowFormat())
				name := exhash.MD5(ztime.NowFormat())
				if len(args) == 2 {
					name = args[1]
				}
				if len(args) == 3 {
					name = args[1]
					user = args[2]
				}
				if len(args) >= 4 {
					name = args[1]
					user = args[2]
					pass = args[3]
				}
				if err := db.Create(name, user, pass); err != nil {
					zlog.Errorf("create new db err: %v", err)
					return
				}
				zlog.Infof("create new db [user: %s,pass: %s,name: %s] done")
			case "drop":
				name := ""
				user := ""
				if len(args) == 2 {
					name = args[1]
				}
				if len(args) >= 3 {
					name = args[1]
					user = args[2]
				}
				if len(name) == 0 {
					zlog.Warn("db name is empty, skip")
					return
				}
				if err := db.Drop(name, user); err != nil {
					zlog.Errorf("drop db err: %v", err)
					return
				}
				zlog.Infof("drop db %s done", name)
			case "show":
				dbs, err := db.Show()
				if err != nil {
					zlog.Errorf("show db err: %v", err)
					return
				}
				zlog.Infof("show db: %v", len(dbs))
				for _, ds := range dbs {
					zlog.Debugf("db: %s", ds.Name)
				}
			default:
				zlog.Debugf("unknown mysql command: %s", args[0])
			}
		},
	}
	mysqlcmd.PersistentFlags().IntVar(&opt.Port, "port", 3306, "mysql port $MYSQL_PORT")
	mysqlcmd.PersistentFlags().StringVar(&opt.User, "user", "root", "mysql user $MYSQL_USER")
	mysqlcmd.PersistentFlags().StringVar(&opt.Pass, "pass", "", "mysql pass $MYSQL_PASS")
	mysqlcmd.PersistentFlags().StringVar(&opt.Name, "name", "", "mysql name $MYSQL_NAME")
	mysqlcmd.PersistentFlags().StringVar(&opt.Dsn, "dsn", "", "mysql dsn, $MYSQL_DSN")
	mysqlcmd.PersistentFlags().StringVar(&opt.Host, "host", "127.0.0.1", "mysql host $MYSQL_HOST")
	return mysqlcmd
}
