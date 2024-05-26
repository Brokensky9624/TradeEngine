package internal

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type baseConnConfig struct {
	ip           string
	port         string
	userName     string
	password     string
	databaseName string
}

var tagNameForDsn = "dsnTag"

type mySQLConnConfig struct {
	baseConnConfig
	Charset         string `dsnTag:"charset"`
	ParseTime       string `dsnTag:"parseTime"`
	Loc             string `dsnTag:"loc"`
	Timeout         string `dsnTag:"timeout"`
	ReadTimeout     string `dsnTag:"readTimeout"`
	WriteTimeout    string `dsnTag:"writeTimeout"`
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

func NewMySQLConnConfig(ip, port, username, password, databaseName string) *mySQLConnConfig {
	return &mySQLConnConfig{
		baseConnConfig: baseConnConfig{
			ip:           ip,
			port:         port,
			userName:     username,
			password:     password,
			databaseName: databaseName,
		},
	}
}

func (config *mySQLConnConfig) SetCharset(charset string) *mySQLConnConfig {
	config.Charset = charset
	return config

}

func (config *mySQLConnConfig) SetParseTime(parseTime string) *mySQLConnConfig {
	config.ParseTime = parseTime
	return config
}

func (config *mySQLConnConfig) SetLoc(loc string) *mySQLConnConfig {
	config.Loc = loc
	return config
}

func (config *mySQLConnConfig) SetTimeout(timeout string) *mySQLConnConfig {
	config.Timeout = timeout
	return config
}

func (config *mySQLConnConfig) SetReadTimeout(readTimeout string) *mySQLConnConfig {
	config.ReadTimeout = readTimeout
	return config
}

func (config *mySQLConnConfig) SetWriteTimeout(writeTimeout string) *mySQLConnConfig {
	config.WriteTimeout = writeTimeout
	return config
}

func (config *mySQLConnConfig) SetMaxIdleConns(maxIdleConns int) *mySQLConnConfig {
	config.MaxIdleConns = maxIdleConns
	return config
}
func (config *mySQLConnConfig) SetMaxOpenConns(maxOpenConns int) *mySQLConnConfig {
	config.MaxOpenConns = maxOpenConns
	return config
}
func (config *mySQLConnConfig) SetConnMaxLifetime(connMaxLifetime time.Duration) *mySQLConnConfig {
	config.ConnMaxLifetime = connMaxLifetime
	return config
}

func (config *mySQLConnConfig) getBaseDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.userName, config.password, config.ip, config.port, config.databaseName)
}

func (config *mySQLConnConfig) getExtraDsn() string {
	tmpDsnList := []string{}
	cfgType := reflect.TypeOf(*config)
	cfgValue := reflect.ValueOf(*config)
	baseType := reflect.TypeOf(config.baseConnConfig)

	for i := 0; i < cfgType.NumField(); i++ {
		field := cfgType.Field(i)
		fieldValue := cfgValue.Field(i)
		if field.Type == baseType {
			continue
		}
		// Check if the field is set to its zero value
		if reflect.Zero(fieldValue.Type()).Interface() != fieldValue.Interface() {
			if tmpDsn := parseToSimpleDsn(field, fieldValue); tmpDsn != "" {
				tmpDsnList = append(tmpDsnList, tmpDsn)
			}
		}
	}
	return strings.Join(tmpDsnList, "&")
}

func (config *mySQLConnConfig) getDSN() string {
	dsn := config.getBaseDsn()
	extraDsn := config.getExtraDsn()
	if len(extraDsn) > 0 {
		dsn = strings.Join([]string{dsn, extraDsn}, "?")
	}
	return dsn
}

func parseToSimpleDsn(field reflect.StructField, fieldValue reflect.Value) string {
	tagName := field.Tag.Get(tagNameForDsn)
	if tagName != "" && fieldValue.String() != "" {
		return fmt.Sprintf("%s=%s", tagName, fieldValue.String())
	}
	return ""
}
