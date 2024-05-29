package member

import (
	"errors"
	"sync"

	"tradeengine/server/web/rest/param"
	"tradeengine/service/db/model"
	dbTypes "tradeengine/service/db/types"
	"tradeengine/service/member/types"

	"gorm.io/gorm"

	"tradeengine/utils/logger"
	"tradeengine/utils/tool"
)

var (
	memberSrv *MemberService
	once      sync.Once
)

func NewService(db *dbTypes.DBService) *MemberService {
	once.Do(func() {
		memberSrv = &MemberService{
			db: db,
		}
	})
	return memberSrv
}

func GetService() *MemberService {
	return memberSrv
}

type MemberService struct {
	db *dbTypes.DBService
}

func (s *MemberService) Auth(param *param.MemberAuthParam) error {
	var errPreFix string = "failed to auth member"

	// check step
	if err := param.Check(); err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}

	model := model.Member{
		Account: param.Account,
	}
	matchMember, err := s.queryMemberByModel(model, true)
	if err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	if !tool.CheckPassword(param.Password, matchMember.Password) {
		return tool.PrefixError(errPreFix, errors.New("password is incorrect"))
	}
	logger.SERVER.Info("member %s auth successfully!\n", param.Account)
	return nil
}

func (s *MemberService) AuthAndMember(param *param.MemberAuthParam) (*types.Member, error) {
	var errPreFix string = "failed to auth member and get"

	// check step
	if err := param.Check(); err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}

	model := model.Member{
		Account: param.Account,
	}
	matchMember, err := s.queryMemberByModel(model, true)
	if err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}
	if !tool.CheckPassword(param.Password, matchMember.Password) {
		return nil, tool.PrefixError(errPreFix, errors.New("password is incorrect"))
	}
	logger.SERVER.Info("member %s auth successfully!\n", param.Account)
	return matchMember, nil
}

func (s *MemberService) Create(param param.MemberCreateParam) (*types.Member, error) {
	var errPreFix string = "failed to create member"

	// check step
	if err := param.Check(); err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return nil, err
	}

	// check user is existed
	findModel := model.Member{
		Account: param.Account,
	}
	if _, err := s.queryMemberByModel(findModel, false); err == nil {
		return nil, tool.PrefixError(errPreFix, errors.New("user is existed"))
	}

	// prepare member info for create
	pwd, err := tool.HashPassword(param.Password)
	if err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return nil, err
	}

	createModel := &model.Member{
		Account:  param.Account,
		Name:     param.Name,
		Password: pwd,
		Email:    param.Email,
		Phone:    param.Phone,
	}

	// create member
	if err = s.db.Create(createModel).Error; err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return nil, err
	}
	logger.SERVER.Info("member %s create successfully!\n", param.Account)
	return ModelToMember(*createModel, false), nil
}

func (s *MemberService) Edit(param param.MemberEditParam) error {
	var errPreFix string = "failed to member edit"

	// check step
	if err := param.Check(); err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}

	// update member
	findModel := model.Member{
		Account: param.Account,
	}
	matchMember, err := s.queryMemberByModel(findModel, false)
	if err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}

	queryModel := model.Member{
		Model: gorm.Model{
			ID: matchMember.ID,
		},
	}
	editModel := types.Member{
		Account: param.Account,
		Name:    param.Name,
		Email:   param.Email,
		Phone:   param.Phone,
	}
	if err := s.db.Where(queryModel).Take(&queryModel).Updates(&editModel).Error; err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("member %s edit successfully!\n", param.Account)
	return nil
}

func (s *MemberService) Delete(param param.MemberDeleteParam) error {
	var errPreFix string = "failed to member delete"

	// check step
	if err := param.Check(); err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}

	// delete member
	account := param.Account
	deleteMember := model.Member{
		Account: account,
	}
	if err := s.db.Where(deleteMember).Take(&deleteMember).Delete(&deleteMember).Error; err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	deleteUnscopedMember := model.Member{
		Account: account,
	}
	if err := s.db.Unscoped().Where(deleteUnscopedMember).Take(&deleteUnscopedMember).Delete(&deleteUnscopedMember).Error; err != nil {
		err = tool.PrefixError(errPreFix, err)
		logger.SERVER.Debug(err.Error())
		return err
	}
	logger.SERVER.Info("member %s delete successfully!\n", param.Account)
	return nil
}

func (s *MemberService) Member(param param.MemberInfoParam) (*types.Member, error) {
	var errPreFix string = "failed to get member"

	// check step
	if err := param.Check(); err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}

	// find member
	findModel := model.Member{
		Account: param.Account,
	}
	matchMember, err := s.queryMemberByModel(findModel, false)
	return matchMember, tool.PrefixError(errPreFix, err)
}

func (s *MemberService) Members() ([]types.Member, error) {
	var errPreFix string = "failed to get member list"

	// find member list
	var memberList []types.Member = make([]types.Member, 0)
	var queryMemberList []model.Member
	if err := s.db.Find(&queryMemberList).Error; err != nil {
		return nil, tool.PrefixError(errPreFix, err)
	}
	for _, memberModel := range queryMemberList {
		memberList = append(memberList, *ModelToMember(memberModel, false))
	}
	return memberList, nil
}

func (s *MemberService) queryMemberByModel(findModel model.Member, includePassword bool) (*types.Member, error) {
	var queryMemberList []model.Member
	if err := s.db.Where(findModel).Find(&queryMemberList).Error; err != nil {
		return nil, err
	}
	var matchMember *types.Member
	for _, queryMember := range queryMemberList {
		if queryMember.Account == findModel.Account {
			matchMember = ModelToMember(queryMember, includePassword)
			break
		}
	}
	if matchMember == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return matchMember, nil
}

func ModelToMember(m model.Member, includePassword bool) *types.Member {
	member := &types.Member{
		ID:      m.ID,
		Account: m.Account,
		Name:    m.Name,
		Email:   m.Email,
		Phone:   m.Phone,
	}
	if includePassword {
		member.Password = m.Password
	}
	return member
}
