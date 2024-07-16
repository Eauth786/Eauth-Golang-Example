package eauth

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Required configuration
const accountKey string = ""            /* Your account key goes here */
const applicationKey string = ""        /* Your application key goes here */
const applicationID string = ""         /* Your application ID goes here */
const applicationVersion string = "1.0" /* Your application version goes here */

// Advanced configuration
const invalidAccountKeyMessage string = "Invalid account key!"
const invalidApplicationKeyMessage string = "Invalid application key!"
const invalidRequestMessage string = "Invalid request!"
const outdatedVersionMessage string = "Outdated version, please upgrade!"
const busySessionsMessage string = "Please try again later!"
const unavaiableSessionMessage string = "Invalid session. Please re-launch the app!"
const usedSessionMessage string = "Why did the computer go to therapy? Because it had a case of 'Request Repeatitis' and couldn't stop asking for the same thing over and over again!"
const overcrowdedSessionMessage string = "Session limit exceeded. Please re-launch the app!"
const expiredSessionMessage string = "Your session has timed out. please re-launch the app!"
const invalidUserMessage string = "Incorrect login credentials!"
const bannedUserMessage string = "Access denied!"
const incorrectHwidMessage string = "Hardware ID mismatch. Please try again with the correct device!"
const expiredUserMessage string = "Your subscription has ended. Please renew to continue using our service!"
const usedNameMessage string = "Username already taken. Please choose a different username!"
const invalidKeyMessage string = "Invalid key. Please enter a valid key!"
const upgradeYourEauthMessage string = "Upgrade your Eauth plan to exceed the limits!"
const invalidEmailMessage string = "The email you entered is either already in use or unavailable or invalid!"

// Dynamic configuration
var initStatus bool = false
var login bool = false
var register bool = false

var sessionID string = ""
var AppName string = ""
var LoggedMessage string = ""
var RegisteredMessage string = ""

var UserRank string = ""
var RegisterDate string = ""
var ExpireDate string = ""
var HWID string = ""

var ErrorMessage string = ""

func raiseError(message string) {
	ErrorMessage = message
}

func ClearConsole() {
	cmd := exec.Command("cmd", "/C", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func computeSHA512(inputString string) string {
	hash := sha512.New()
	hash.Write([]byte(inputString))
	return hex.EncodeToString(hash.Sum(nil))
}

func generateAuthToken(message, appID string) string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)[:4]
	authToken := timestamp + message + appID
	return computeSHA512(authToken)
}

func RunRequest(data string) string {
	url := "https://eauth.us.to/api/1.1/"

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		fmt.Println("Error creating request:", err)
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}

	req.Header.Set("User-Agent", "e_a_u_t_h")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}

	var requestData map[string]interface{}
	err = json.Unmarshal([]byte(string(responseBody)), &requestData)
	if err != nil {
		fmt.Println("Error:", err)
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}
	if requestData["message"].(string) == "init_success" || requestData["message"].(string) == "login_success" {
		if resp.Header.Get("Key") != generateAuthToken(string(responseBody), applicationID) {
			os.Exit(1)
		}
	}

	return string(responseBody)
}

func userHWID() string {
	out, err := exec.Command("cmd", "/C", "wmic useraccount where name='%username%' get sid").Output()
	if err != nil {
		fmt.Println("Error:", err)
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}

	return string(strings.TrimSpace(strings.TrimPrefix(string(out), "SID"))) + "wffwz"
}

func Init() bool {
	if initStatus {
		return initStatus
	}

	var data map[string]interface{}
	err := json.Unmarshal([]byte(RunRequest("sort=init&appkey="+applicationKey+"&acckey="+accountKey+"&version="+applicationVersion+"&hwid="+userHWID())), &data)
	if err != nil {
		fmt.Println("Error:", err)
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}
	message := data["message"].(string)

	if message == "init_success" {
		initStatus = true
		sessionID = data["session_id"].(string)
		AppName = data["app_name"].(string)
		exec.Command("cmd", "/C", "title", AppName).Run()
		LoggedMessage = data["logged_message"].(string)
		RegisteredMessage = data["registered_message"].(string)
	} else if message == "invalid_account_key" {
		raiseError(invalidAccountKeyMessage)
	} else if message == "invalid_application_key" {
		raiseError(invalidApplicationKeyMessage)
	} else if message == "invalid_request" {
		raiseError(invalidRequestMessage)
	} else if message == "version_outdated" {
		downloadLink := data["download_link"].(string)
		if downloadLink != "" {
			exec.Command("cmd", "/C", "start", downloadLink).Start()
		}
		raiseError(outdatedVersionMessage)
	} else if message == "maximum_sessions_reached" {
		raiseError(busySessionsMessage)
	} else if message == "user_is_banned" {
		raiseError(bannedUserMessage)
	} else if message == "init_paused" {
		raiseError(data["paused_message"].(string))
	}

	return initStatus
}

func Login(username string, password string) bool {
	if login {
		return login
	}

	var data map[string]interface{}
	err := json.Unmarshal([]byte(RunRequest("sort=login&sessionid="+sessionID+"&username="+username+"&password="+password+"&hwid="+userHWID())), &data)
	if err != nil {
		fmt.Println("Error:", err)
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}
	message := data["message"].(string)

	if message == "login_success" {
		// Login success
		login = true
		UserRank = data["rank"].(string)
		RegisterDate = data["register_date"].(string)
		ExpireDate = data["expire_date"].(string)
		HWID = data["hwid"].(string)
	} else if message == "invalid_account_key" {
		raiseError(invalidAccountKeyMessage)
	} else if message == "session_unavailable" {
		raiseError(unavaiableSessionMessage)
	} else if message == "invalid_request" {
		raiseError(invalidRequestMessage) // This is usually not the case
	} else if message == "session_already_used" {
		raiseError(usedSessionMessage)
	} else if message == "session_overcrowded" {
		raiseError(overcrowdedSessionMessage)
	} else if message == "session_expired" {
		raiseError(expiredSessionMessage)
	} else if message == "account_unavailable" {
		raiseError(invalidUserMessage)
	} else if message == "user_is_banned" {
		raiseError(bannedUserMessage)
	} else if message == "hwid_incorrect" {
		raiseError(incorrectHwidMessage + "\n" + "HWID reset is available " + data["estimated_reset_time"].(string))
	} else if message == "subscription_expired" {
		raiseError(expiredUserMessage)
	}

	return login
}

func Register(username string, email string, password string, key string) bool {
	if register {
		return register
	}

	var data map[string]interface{}
	err := json.Unmarshal([]byte(RunRequest("sort=register&sessionid="+sessionID+"&username="+username+"&email="+email+"&password="+password+"&key="+key+"&hwid="+userHWID())), &data)
	if err != nil {
		fmt.Println("Error:", err)
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}
	message := data["message"].(string)

	switch message {
	case "register_success":
		// Register success
		register = true
	case "session_unavailable":
		raiseError(unavaiableSessionMessage)
	case "invalid_email":
		raiseError(invalidEmailMessage)
	case "session_already_used":
		raiseError(usedSessionMessage)
	case "invalid_request":
		raiseError(invalidRequestMessage) // This is usually not the case
	case "invalid_account_key":
		raiseError(invalidAccountKeyMessage)
	case "session_overcrowded":
		raiseError(overcrowdedSessionMessage)
	case "session_expired":
		raiseError(expiredSessionMessage)
	case "name_already_used":
		raiseError(usedNameMessage)
	case "key_unavailable":
		raiseError(invalidKeyMessage)
	case "maximum_users_reached":
		raiseError(upgradeYourEauthMessage)
	case "user_is_banned":
		raiseError(bannedUserMessage)
	}

	return register
}
