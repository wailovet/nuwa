package nuwa

import (
	"crypto/rsa"
	"crypto/tls"
	"time"
	"unsafe"

	"github.com/go-sql-driver/mysql"
	"xorm.io/core"
	"xorm.io/xorm"
)

type MysqledConfig struct {
	User                    string            `json:"user,omitempty"`                      // Username
	Passwd                  string            `json:"passwd,omitempty"`                    // Password (requires User)
	Net                     string            `json:"net,omitempty"`                       // Network type
	Addr                    string            `json:"addr,omitempty"`                      // Network address (requires Net)
	DBName                  string            `json:"db_name,omitempty"`                   // Database name
	Params                  map[string]string `json:"params,omitempty"`                    // Connection parameters
	Collation               string            `json:"collation,omitempty"`                 // Connection collation
	Loc                     *time.Location    `json:"loc,omitempty"`                       // Location for time.Time values
	MaxAllowedPacket        int               `json:"max_allowed_packet,omitempty"`        // Max packet size allowed
	ServerPubKey            string            `json:"server_pub_key,omitempty"`            // Server public key name
	pubKey                  *rsa.PublicKey    `json:"-"`                                   // Server public key
	TLSConfig               string            `json:"tls_config,omitempty"`                // TLS configuration name
	tls                     *tls.Config       `json:"-"`                                   // TLS configuration
	Timeout                 time.Duration     `json:"timeout,omitempty"`                   // Dial timeout
	ReadTimeout             time.Duration     `json:"read_timeout,omitempty"`              // I/O read timeout
	WriteTimeout            time.Duration     `json:"write_timeout,omitempty"`             // I/O write timeout
	AllowAllFiles           bool              `json:"allow_all_files,omitempty"`           // Allow all files to be used with LOAD DATA LOCAL INFILE
	AllowCleartextPasswords bool              `json:"allow_cleartext_passwords,omitempty"` // Allows the cleartext client side plugin
	AllowNativePasswords    bool              `json:"allow_native_passwords,omitempty"`    // Allows the native password authentication method
	AllowOldPasswords       bool              `json:"allow_old_passwords,omitempty"`       // Allows the old insecure password method
	ClientFoundRows         bool              `json:"client_found_rows,omitempty"`         // Return number of matching rows instead of rows changed
	ColumnsWithAlias        bool              `json:"columns_with_alias,omitempty"`        // Prepend table alias to column names
	InterpolateParams       bool              `json:"interpolate_params,omitempty"`        // Interpolate placeholders into query string
	MultiStatements         bool              `json:"multi_statements,omitempty"`          // Allow multiple statements in one query
	ParseTime               bool              `json:"parse_time,omitempty"`                // Parse time values to time.Time
	RejectReadOnly          bool              `json:"reject_read_only,omitempty"`          // Reject read-only connections
}

type Mysqled struct {
	engine *xorm.Engine
	config *MysqledConfig
	dsn    string
	prefix string
	isLog  bool
}

func (s *Mysqled) ConfigDSN(dsn string, tablePrefix ...string) (err error) {
	s.dsn = dsn
	mconfig, err := mysql.ParseDSN(dsn)
	if err != nil {
		return err
	}

	s.config = (*MysqledConfig)(unsafe.Pointer(mconfig))

	s.engine, err = xorm.NewEngine("mysql", dsn)
	if err != nil {
		return err
	}

	prefix := ""
	if len(tablePrefix) > 0 {
		prefix = tablePrefix[0]
	}
	s.engine.SetTableMapper(core.NewPrefixMapper(core.SnakeMapper{}, prefix))
	return
}

func (s *Mysqled) Config(config *MysqledConfig, tablePrefix ...string) (err error) {
	s.config = config
	mconfig := (*mysql.Config)(unsafe.Pointer(config))
	s.dsn = mconfig.FormatDSN()
	s.engine, err = xorm.NewEngine("mysql", s.dsn)
	if err != nil {
		return err
	}

	prefix := ""
	if len(tablePrefix) > 0 {
		prefix = tablePrefix[0]
	}
	s.engine.SetTableMapper(core.NewPrefixMapper(core.SnakeMapper{}, prefix))
	return
}

func (s *Mysqled) SimpleConfig(addr, username, password, dbName string, tablePrefix ...string) (err error) {
	mysqlConfig := mysql.NewConfig()
	mysqlConfig.User = username
	mysqlConfig.DBName = dbName
	mysqlConfig.Passwd = password
	mysqlConfig.Net = "tcp"
	mysqlConfig.Addr = addr
	mysqlConfig.Params = map[string]string{
		"charset":   "utf8mb4",
		"parseTime": "true",
	}

	return s.ConfigDSN(mysqlConfig.FormatDSN(), tablePrefix...)
}

func (s *Mysqled) SetPerformanceConfig(maxIdle int, maxActive int, maxLifetime time.Duration) {
	s.engine.SetMaxIdleConns(maxIdle)
	s.engine.SetMaxOpenConns(maxActive)
	s.engine.SetConnMaxLifetime(maxLifetime)
}

func (s *Mysqled) Xorm() *xorm.Engine {
	return s.engine
}
