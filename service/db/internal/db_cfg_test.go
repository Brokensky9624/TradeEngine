package internal

import (
	"reflect"
	"testing"
)

func TestNewMySQLConnConfig(t *testing.T) {
	type args struct {
		ip           string
		port         string
		username     string
		password     string
		databaseName string
	}
	tests := []struct {
		name string
		args args
		want *mySQLConnConfig
	}{
		{
			name: "check correct case of TestNewMySQLConnConfig",
			args: args{
				"127.0.0.1",
				"3306",
				"initial_user",
				"password",
				"testDB",
			},
			want: &mySQLConnConfig{
				baseConnConfig: baseConnConfig{
					ip:           "127.0.0.1",
					port:         "3306",
					userName:     "initial_user",
					password:     "password",
					databaseName: "testDB",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMySQLConnConfig(tt.args.ip, tt.args.port, tt.args.username, tt.args.password, tt.args.databaseName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMySQLConnConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mySQLConnConfig_SetCharset(t *testing.T) {
	type args struct {
		charset string
	}
	tests := []struct {
		name   string
		config *mySQLConnConfig
		args   args
		want   *mySQLConnConfig
	}{
		{
			name: "check correct case of Test_mySQLConnConfig_SetCharset",
			config: NewMySQLConnConfig(
				"127.0.0.1",
				"3306",
				"initial_user",
				"password",
				"testDB",
			),
			args: args{
				charset: "utf8mb4",
			},
			want: &mySQLConnConfig{
				baseConnConfig: baseConnConfig{
					ip:           "127.0.0.1",
					port:         "3306",
					userName:     "initial_user",
					password:     "password",
					databaseName: "testDB",
				},
				Charset: "utf8mb4",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetCharset(tt.args.charset); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mySQLConnConfig.SetCharset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mySQLConnConfig_SetParseTime(t *testing.T) {
	type args struct {
		parseTime string
	}
	tests := []struct {
		name   string
		config *mySQLConnConfig
		args   args
		want   *mySQLConnConfig
	}{
		{
			name: "check correct case of Test_mySQLConnConfig_SetParseTime",
			config: NewMySQLConnConfig(
				"127.0.0.1",
				"3306",
				"initial_user",
				"password",
				"testDB",
			),
			args: args{
				parseTime: "True",
			},
			want: &mySQLConnConfig{
				baseConnConfig: baseConnConfig{
					ip:           "127.0.0.1",
					port:         "3306",
					userName:     "initial_user",
					password:     "password",
					databaseName: "testDB",
				},
				ParseTime: "True",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetParseTime(tt.args.parseTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mySQLConnConfig.SetParseTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mySQLConnConfig_SetLoc(t *testing.T) {
	type args struct {
		loc string
	}
	tests := []struct {
		name   string
		config *mySQLConnConfig
		args   args
		want   *mySQLConnConfig
	}{
		{
			name: "check correct case of Test_mySQLConnConfig_SetLoc",
			config: NewMySQLConnConfig(
				"127.0.0.1",
				"3306",
				"initial_user",
				"password",
				"testDB",
			),
			args: args{
				loc: "Local",
			},
			want: &mySQLConnConfig{
				baseConnConfig: baseConnConfig{
					ip:           "127.0.0.1",
					port:         "3306",
					userName:     "initial_user",
					password:     "password",
					databaseName: "testDB",
				},
				Loc: "Local",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetLoc(tt.args.loc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mySQLConnConfig.SetLoc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mySQLConnConfig_SetTimeout(t *testing.T) {
	type args struct {
		timeout string
	}
	tests := []struct {
		name   string
		config *mySQLConnConfig
		args   args
		want   *mySQLConnConfig
	}{
		{
			name: "check correct case of Test_mySQLConnConfig_SetTimeout",
			config: NewMySQLConnConfig(
				"127.0.0.1",
				"3306",
				"initial_user",
				"password",
				"testDB",
			),
			args: args{
				timeout: "30s",
			},
			want: &mySQLConnConfig{
				baseConnConfig: baseConnConfig{
					ip:           "127.0.0.1",
					port:         "3306",
					userName:     "initial_user",
					password:     "password",
					databaseName: "testDB",
				},
				Timeout: "30s",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetTimeout(tt.args.timeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mySQLConnConfig.SetTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mySQLConnConfig_SetReadTimeout(t *testing.T) {
	type args struct {
		readTimeout string
	}
	tests := []struct {
		name   string
		config *mySQLConnConfig
		args   args
		want   *mySQLConnConfig
	}{
		{
			name: "check correct case of Test_mySQLConnConfig_SetReadTimeout",
			config: NewMySQLConnConfig(
				"127.0.0.1",
				"3306",
				"initial_user",
				"password",
				"testDB",
			),
			args: args{
				readTimeout: "30s",
			},
			want: &mySQLConnConfig{
				baseConnConfig: baseConnConfig{
					ip:           "127.0.0.1",
					port:         "3306",
					userName:     "initial_user",
					password:     "password",
					databaseName: "testDB",
				},
				ReadTimeout: "30s",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetReadTimeout(tt.args.readTimeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mySQLConnConfig.SetReadTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mySQLConnConfig_SetWriteTimeout(t *testing.T) {
	type args struct {
		writeTimeout string
	}
	tests := []struct {
		name   string
		config *mySQLConnConfig
		args   args
		want   *mySQLConnConfig
	}{
		{
			name: "check correct case of Test_mySQLConnConfig_SetWriteTimeout",
			config: NewMySQLConnConfig(
				"127.0.0.1",
				"3306",
				"initial_user",
				"password",
				"testDB",
			),
			args: args{
				writeTimeout: "30s",
			},
			want: &mySQLConnConfig{
				baseConnConfig: baseConnConfig{
					ip:           "127.0.0.1",
					port:         "3306",
					userName:     "initial_user",
					password:     "password",
					databaseName: "testDB",
				},
				WriteTimeout: "30s",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetWriteTimeout(tt.args.writeTimeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mySQLConnConfig.SetWriteTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mySQLConnConfig_getBaseDsn(t *testing.T) {
	tests := []struct {
		name string
		b    *mySQLConnConfig
		want string
	}{
		{
			name: "check correct case of Test_mySQLConnConfig_getBaseDsn",
			b: NewMySQLConnConfig(
				"127.0.0.1",
				"3306",
				"initial_user",
				"password",
				"testDB",
			),
			want: "initial_user:password@tcp(127.0.0.1:3306)/testDB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.getBaseDsn(); got != tt.want {
				t.Errorf("mySQLConnConfig.getBaseDsn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mySQLConnConfig_getExtraDsn(t *testing.T) {
	tests := []struct {
		name string
		b    *mySQLConnConfig
		want string
	}{
		{
			name: "check correct case of Test_mySQLConnConfig_getExtraDsn",
			b: NewMySQLConnConfig(
				"127.0.0.1",
				"3306",
				"initial_user",
				"password",
				"testDB",
			).
				SetCharset("utf8mb4").
				SetParseTime("True").
				SetLoc("Local").
				SetTimeout("10s").
				SetReadTimeout("30s").
				SetWriteTimeout("30s"),
			want: "charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.getExtraDsn(); got != tt.want {
				t.Errorf("mySQLConnConfig.getExtraDsn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mySQLConnConfig_getDSN(t *testing.T) {
	tests := []struct {
		name string
		b    *mySQLConnConfig
		want string
	}{
		{
			name: "check correct case of Test_mySQLConnConfig_getDSN",
			b: NewMySQLConnConfig(
				"127.0.0.1",
				"3306",
				"initial_user",
				"password",
				"testDB",
			).
				SetCharset("utf8mb4").
				SetParseTime("True").
				SetLoc("Local").
				SetTimeout("10s").
				SetReadTimeout("30s").
				SetWriteTimeout("30s"),
			want: "initial_user:password@tcp(127.0.0.1:3306)/testDB?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.getDSN(); got != tt.want {
				t.Errorf("mySQLConnConfig.getDSN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseToSimpleDsn(t *testing.T) {
	type args struct {
		field      reflect.StructField
		fieldValue reflect.Value
	}
	type tmpConfig struct {
		abc string `dsnTag:"abc"`
	}
	getFieldAndFieldValue := func(config *tmpConfig) (reflect.StructField, reflect.Value) {
		var field reflect.StructField
		var fieldValue reflect.Value
		cfgType := reflect.TypeOf(*config)
		cfgValue := reflect.ValueOf(*config)
		for i := 0; i < cfgType.NumField(); i++ {
			field = cfgType.Field(i)
			fieldValue = cfgValue.Field(i)
			break
		}
		return field, fieldValue
	}
	tmpCfg := &tmpConfig{abc: "123"}
	field1, fieldValue1 := getFieldAndFieldValue(tmpCfg)
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "check correct case of Test_parseToSimpleDsn",
			args: args{
				field:      field1,
				fieldValue: fieldValue1,
			},
			want: "abc=123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseToSimpleDsn(tt.args.field, tt.args.fieldValue); got != tt.want {
				t.Errorf("parseToSimpleDsn() = %v, want %v", got, tt.want)
			}
		})
	}
}
