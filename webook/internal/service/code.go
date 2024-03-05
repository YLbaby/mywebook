package service

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/service/sms"
	"context"
	"fmt"
	"math/rand"
)

const codeTplId = "1877556"

type CodeService interface {
}

type codeService struct {
	repo   repository.CachedCodeRepository
	smsSvc sms.Service
}

func (svc *codeService) Send(ctx context.Context, biz, phone string) error {
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	// 到这儿验证码已经保存到redis里面了，现在需要发送出去给用户
	err = svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
	return err
}

func (svc *codeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%6d", num)
}
