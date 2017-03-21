package dbconf

import "fmt"

import "github.com/magiconair/properties"

// DbConf wraps connection details to the database
type DbConf struct {
	user     string
	password string
	dbName   string
	host     string
	port     int
}

// New returns a new DbConf reading from rhn.conf
func New() *DbConf {
	p := properties.MustLoadFile("/etc/rhn/rhn.conf", properties.UTF8)

	return &DbConf{p.MustGetString("db_user"),
		p.MustGetString("db_password"),
		p.MustGetString("db_name"),
		p.MustGetString("db_host"),
		p.MustGetInt("db_port")}
}

// ConnectionString returns a libpq connection string
func (d *DbConf) ConnectionString() (result string) {
	return fmt.Sprintf("user='%s' password='%s' dbname='%s' host='%s' port='%d' sslmode=disable", d.user, d.password, d.dbName, d.host, d.port)
}
