package twilio

import (
	"fmt"

	"github.com/akmal4410/gestapo/internal/config"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

type TwilioService interface {
	SendOTP(to string) error
	VerfiyOTP(to, code string) (bool, error)
}

type OTPService struct {
	twilio *config.Twilio
}

func NewOTPService(twilio *config.Twilio) TwilioService {
	return &OTPService{
		twilio: twilio,
	}
}

func (service *OTPService) SendOTP(to string) error {

	var client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: service.twilio.AccountSid,
		Password: service.twilio.AuthToken,
	})

	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")
	// params.SetCustomMessage("Your [Gestapo] verification code is:\n")

	resp, err := client.VerifyV2.CreateVerification(service.twilio.ServiceSid, params)
	if err != nil {
		return err
	}

	fmt.Printf("Verification has been send, Id :'%s'\n", *resp.AccountSid)
	return nil
}

func (service OTPService) VerfiyOTP(to, code string) (bool, error) {

	var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: service.twilio.AccountSid,
		Password: service.twilio.AuthToken,
	})
	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(service.twilio.ServiceSid, params)

	if err != nil {
		fmt.Println("error :", err.Error())
		return false, err
	}
	if *resp.Status == "approved" {
		return true, nil
	}
	return false, nil
}
