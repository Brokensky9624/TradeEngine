package member

import (
	"context"
	"testing"
	"tradeengine/server/web/rest/param"
	"tradeengine/service/db"
	"tradeengine/service/wallet"
)

func init() {
	rootContext := context.WithoutCancel(context.Background())
	dbMngr := db.NewDBManager(rootContext)
	dbMngr.Run()
	walletSrv := wallet.NewService(dbMngr.DefaultDBService())
	NewService(dbMngr.DefaultDBService(), walletSrv)
}

func TestMemberService_Create(t *testing.T) {
	type args struct {
		param []param.MemberCreateParam
	}
	tests := []struct {
		name    string
		s       *MemberService
		args    args
		wantErr bool
	}{
		{
			name: "create 5 members",
			s:    GetService(),
			args: args{
				param: []param.MemberCreateParam{
					{
						Account:  "Jason1",
						Name:     "Jason1",
						Password: "12345678",
						Email:    "test@gmail.com",
					},
					{
						Account:  "Jason2",
						Name:     "Jason2",
						Password: "12345678",
						Email:    "test@gmail.com",
					},
					{
						Account:  "Jason3",
						Name:     "Jason3",
						Password: "12345678",
						Email:    "test@gmail.com",
					},
					{
						Account:  "Jason4",
						Name:     "Jason4",
						Password: "12345678",
						Email:    "test@gmail.com",
					},
					{
						Account:  "Jason5",
						Name:     "Jason5",
						Password: "12345678",
						Email:    "test@gmail.com",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, p := range tt.args.param {
				if err := tt.s.Create(p); (err != nil) != tt.wantErr {
					t.Errorf("MemberService.Create() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
