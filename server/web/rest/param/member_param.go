package param

import "tradeengine/utils/tool"

type MemberAuthParam struct {
	ID       uint   `json:"id"`
	Account  string `json:"account" required:"true"`
	Password string `json:"password" required:"true"`
}

func (param MemberAuthParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type MemberCreateParam struct {
	Account  string `json:"account" required:"true"`
	Name     string `json:"name" required:"true"`
	Password string `json:"password" required:"true"`
	Email    string `json:"email" required:"true"`
	Phone    string `json:"phone"`
}

func (param MemberCreateParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type MemberEditParam struct {
	Account string `json:"account" required:"true"`
	Name    string `json:"name" required:"true"`
	Email   string `json:"email" required:"true"`
	Phone   string `json:"phone" required:"true"`
}

func (param MemberEditParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type MemberDeleteParam struct {
	Account string `json:"account" required:"true"`
}

func (param MemberDeleteParam) Check() error {
	return tool.CheckRequiredFields(param)
}

type MemberInfoParam struct {
	Account string `json:"account" required:"true"`
}

func (param MemberInfoParam) Check() error {
	return tool.CheckRequiredFields(param)
}
