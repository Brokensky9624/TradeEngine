package internal

import (
	"reflect"
	"testing"
	"time"
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetWriteTimeout(tt.args.writeTimeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mySQLConnConfig.SetWriteTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mySQLConnConfig_SetMaxIdleConns(t *testing.T) {
	type args struct {
		maxIdleConns int
	}
	tests := []struct {
		name   string
		config *mySQLConnConfig
		args   args
		want   *mySQLConnConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetMaxIdleConns(tt.args.maxIdleConns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mySQLConnConfig.SetMaxIdleConns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mySQLConnConfig_SetMaxOpenConns(t *testing.T) {
	type args struct {
		maxOpenConns int
	}
	tests := []struct {
		name   string
		config *mySQLConnConfig
		args   args
		want   *mySQLConnConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetMaxOpenConns(tt.args.maxOpenConns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mySQLConnConfig.SetMaxOpenConns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mySQLConnConfig_SetConnMaxLifetime(t *testing.T) {
	type args struct {
		connMaxLifetime time.Duration
	}
	tests := []struct {
		name   string
		config *mySQLConnConfig
		args   args
		want   *mySQLConnConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.config.SetConnMaxLifetime(tt.args.connMaxLifetime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mySQLConnConfig.SetConnMaxLifetime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mySQLConnConfig_getBaseDsn(t *testing.T) {
	tests := []struct {
		name string
		b    *mySQLConnConfig
		want string
	}{}
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
		// TODO: Add test cases.
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
			name: "full",
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
				SetWriteTimeout("30s").
				SetMaxIdleConns(10).
				SetMaxOpenConns(500).
				SetConnMaxLifetime(time.Minute * 60),
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
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseToSimpleDsn(tt.args.field, tt.args.fieldValue); got != tt.want {
				t.Errorf("parseToSimpleDsn() = %v, want %v", got, tt.want)
			}
		})
	}
}
