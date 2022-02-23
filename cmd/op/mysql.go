package op

import (
	"fmt"

	"github.com/ergoapi/log"
	"github.com/ergoapi/util/environ"
	"github.com/spf13/cobra"
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
		opt.Dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", opt.User, opt.Pass, opt.User, opt.Pass, opt.Name)
	}
}

func MysqlCmd() *cobra.Command {
	opt := &mysqlOption{}
	mysqlcmd := &cobra.Command{
		Use:     "mysql",
		Short:   "ping",
		Version: "3.1.0",
		Run: func(cobraCmd *cobra.Command, args []string) {
			zlog := log.GetInstance()
			opt.Check()
			zlog.Debugf("mysql dsn: %s", opt.Dsn)
			zlog.Debugf("args: %v", len(args))
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
