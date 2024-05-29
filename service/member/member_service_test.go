package member

import (
	"context"
	"fmt"
	"testing"
	"tradeengine/server/web/rest/param"
	"tradeengine/service/db"
	"tradeengine/service/stock"
	"tradeengine/service/stockinfo"
	"tradeengine/service/wallet"
)

var (
	walletSrv    *wallet.WalletService
	stockSrv     *stock.StockService
	stockInfoSrv *stockinfo.StockInfoService
)

func init() {
	rootContext := context.WithoutCancel(context.Background())
	dbMngr := db.NewDBManager(rootContext)
	dbMngr.Run()
	NewService(dbMngr.DefaultDBService())
	walletSrv = wallet.NewService(dbMngr.DefaultDBService())
	stockSrv = stock.NewService(dbMngr.DefaultDBService())
	stockInfoSrv = stockinfo.NewService(dbMngr.DefaultDBService())
}

func TestMemberService_Create(t *testing.T) {
	type args struct {
		username string
		param    []param.MemberCreateParam
	}
	tests := []struct {
		name    string
		s       *MemberService
		args    args
		wantErr bool
	}{
		{
			name: "create 5 members Jason1 ~ Jason5",
			s:    GetService(),
			args: args{
				username: "Jason",
				param: []param.MemberCreateParam{
					{
						Password: "12345678",
						Email:    "test@gmail.com",
					},
					{
						Password: "12345678",
						Email:    "test@gmail.com",
					},
					{
						Password: "12345678",
						Email:    "test@gmail.com",
					},
					{
						Password: "12345678",
						Email:    "test@gmail.com",
					},
					{
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
			for i, p := range tt.args.param {
				p.Account = fmt.Sprintf("%s%d", tt.args.username, (i + 1))
				p.Name = p.Account
				if _, err := tt.s.Create(p); (err != nil) != tt.wantErr {
					t.Errorf("MemberService.Create() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestMemberService_Create_With1000G_50StockID1(t *testing.T) {
	type args struct {
		username      string
		param         []param.MemberCreateParam
		money         uint
		stockQuantity uint
		stockInfoID   uint
	}
	tests := []struct {
		name      string
		memberSrv *MemberService
		walletSrv *wallet.WalletService
		stockSrv  *stock.StockService
		args      args
		wantErr   bool
	}{
		{
			name:      "create 5 members Jason1 ~ Jason5",
			memberSrv: GetService(),
			walletSrv: walletSrv,
			stockSrv:  stockSrv,
			args: args{
				username: "Jason",
				param: []param.MemberCreateParam{
					{
						Password: "12345678",
						Email:    "test@gmail.com",
					},
					{
						Password: "12345678",
						Email:    "test@gmail.com",
					},
					{
						Password: "12345678",
						Email:    "test@gmail.com",
					},
					{
						Password: "12345678",
						Email:    "test@gmail.com",
					},
					{
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
			for i, p := range tt.args.param {
				p.Account = fmt.Sprintf("%s%d", tt.args.username, (i + 1))
				p.Name = p.Account
				oneMember, err := tt.memberSrv.Create(p)
				if (err != nil) != tt.wantErr {
					// t.Errorf("MemberService.Create() error = %v, wantErr %v", err, tt.wantErr)
				}
				walletParam := param.WalletCreateParam{
					OwnerID:        oneMember.ID,
					AvailableMoney: 1000,
					PendingMoney:   0,
				}
				if err := tt.walletSrv.Create(walletParam); (err != nil) != tt.wantErr {
					// t.Errorf("WalletService.Create() error = %v, wantErr %v", err, tt.wantErr)
				}
				stockParam := param.OneStockCreateParam{
					OwnerID:           oneMember.ID,
					StockInfoID:       1,
					AvailableQuantity: 50,
					PendingQuantity:   0,
				}
				if err := tt.stockSrv.Create(stockParam); (err != nil) != tt.wantErr {
					// t.Errorf("StockService.Create() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestMemberService_Delete(t *testing.T) {
	type args struct {
		username string
		param    []param.MemberDeleteParam
	}
	tests := []struct {
		name    string
		s       *MemberService
		args    args
		wantErr bool
	}{
		{
			name: "delete 5 members Jason1 ~ Jason5",
			s:    GetService(),
			args: args{
				username: "Jason",
				param: []param.MemberDeleteParam{
					{},
					{},
					{},
					{},
					{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, p := range tt.args.param {
				p.Account = fmt.Sprintf("%s%d", tt.args.username, (i + 1))
				if err := tt.s.Delete(p); (err != nil) != tt.wantErr {
					t.Errorf("MemberService.Create() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestMemberService_Members(t *testing.T) {
	tests := []struct {
		name    string
		s       *MemberService
		wantErr bool
	}{
		{
			name:    "list all members",
			s:       GetService(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Members()
			if (err != nil) != tt.wantErr {
				t.Errorf("MemberService.Members() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(got)
		})
	}
}
