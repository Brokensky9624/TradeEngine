package types

import "tradeengine/server/web/rest/param"

type IMemberSrv interface {
	Auth(param *param.MemberAuthParam) error
	AuthAndMember(param *param.MemberAuthParam) (*Member, error)
	Create(param param.MemberCreateParam) error
	Edit(param param.MemberEditParam) error
	Delete(param param.MemberDeleteParam) error
	Member(param param.MemberInfoParam) (*Member, error)
	Members() ([]Member, error)
}

type Member struct {
	ID       uint   `json:"id"`
	Account  string `json:"account"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}
