package uservalidator

import (
	"fmt"
	"gameapp/dto"
	"gameapp/pkg/errmsg"
	"gameapp/pkg/richerror"
	"github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
}

type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{
		repo: repo,
	}
}

func (v Validator) ValidateRegisterRequest(req dto.RegisterRequest) (error, map[string]string) {
	const op = "uservalidator.ValidateRegisterRequest"
	// TODO: we should verify phone number by verification code

	if err := validation.ValidateStruct(&req,
		validation.Field(
			&req.Name,
			validation.Required,
			validation.Length(3, 50),
		),

		validation.Field(
			&req.Password,
			validation.Required,
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9!@#%^&*]{8,}$`)),
		),

		validation.Field(
			&req.PhoneNumber,
			validation.Required,
			validation.Match(regexp.MustCompile("^09[0-9]{9}$")).Error("phone number is not valid"),
			validation.By(v.checkPhoneNumberUniqueness),
		),
	); err != nil {
		fieldErrors := make(map[string]string)
		errV, ok := err.(validation.Errors)
		if ok {
			for key, value := range errV {
				if value != nil {
					fieldErrors[key] = value.Error()
				}
			}
		}

		return richerror.New(op).
			WithMessage(errmsg.ErrorMsgInvalidInput).WithErr(err).
			WithKind(richerror.KindInvalid).
			WithMeta(map[string]interface{}{"req": req}), fieldErrors
	}

	return nil, nil
}

func (v Validator) checkPhoneNumberUniqueness(value interface{}) error {
	phoneNumber := value.(string)

	if isUnique, err := v.repo.IsPhoneNumberUnique(phoneNumber); err != nil || !isUnique {
		if err != nil {
			return err
		}
		if !isUnique {
			return fmt.Errorf(errmsg.ErrorMsgPhoneNumberIsNotUnique)
		}
	}
	return nil
}
