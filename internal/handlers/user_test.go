package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
	"github.com/xuri/excelize/v2"

	mocks "forum/internal/repo/mocks"
)

var Log = logrus.New()

func InitLogger() {
	Log.SetOutput(os.Stdout)
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Log.SetLevel(logrus.InfoLevel)
}

func TestMain(m *testing.M) {
	InitLogger()
	logrus.Info("=== Starting Test Suite ===")
	exitCode := m.Run()
	logrus.Info("=== Test Suite Completed ===")
	os.Exit(exitCode)
}

type SignupTestCase struct {
	Name          string
	Username      string
	Email         string
	Password      string
	PasswordAgain string
	WantCode      int
}

func loadSignupTestData(fileName, sheetName string) ([]SignupTestCase, error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", fileName, err)
	}
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows from sheet %s: %v", sheetName, err)
	}

	var tests []SignupTestCase
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 6 {
			continue
		}
		wantCode, err := strconv.Atoi(row[5])
		if err != nil {
			return nil, fmt.Errorf("invalid WantCode in row %d: %w", i, err)
		}
		testCase := SignupTestCase{
			Name:          row[0],
			Username:      row[1],
			Email:         row[2],
			Password:      row[3],
			PasswordAgain: row[4],
			WantCode:      wantCode,
		}
		tests = append(tests, testCase)
	}
	return tests, nil
}

type LoginTestCase struct {
	Name     string
	Email    string
	Password string
	WantCode int
}

func loadLoginTestData(fileName, sheetName string) ([]LoginTestCase, error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", fileName, err)
	}
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows from sheet %s: %v", sheetName, err)
	}

	var tests []LoginTestCase
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 4 {
			continue
		}
		wantCode, err := strconv.Atoi(row[3])
		if err != nil {
			return nil, fmt.Errorf("invalid WantCode in row %d: %v", i, err)
		}
		testCase := LoginTestCase{
			Name:     row[0],
			Email:    row[1],
			Password: row[2],
			WantCode: wantCode,
		}
		tests = append(tests, testCase)
	}
	return tests, nil
}

func TestSignUp(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	logrus.Info("TestSignUp: Starting Excel-driven tests for /signup")

	signupTests, err := loadSignupTestData("testdata_signup.xlsx", "Sheet1")
	if err != nil {
		t.Fatalf("Error loading signup test data: %v", err)
	}

	for _, tt := range signupTests {
		t.Run(tt.Name, func(t *testing.T) {
			logrus.Infof("Running signup test case: %q", tt.Name)

			form := url.Values{}
			form.Add("name", tt.Username)
			form.Add("email", tt.Email)
			form.Add("password", tt.Password)
			form.Add("password", tt.PasswordAgain)

			code, _, _ := ts.postForm(t, "/signup", form)

			if code != tt.WantCode {
				logrus.Errorf("Signup test FAILED for %q: got code %d, want %d", tt.Name, code, tt.WantCode)
			} else {
				logrus.Infof("Signup test PASSED for %q: got code %d (as expected)", tt.Name, code)
			}
			mocks.Equal(t, code, tt.WantCode)
		})
	}
	logrus.Info("TestSignUp: Completed Excel-driven tests for /signup")
}

func TestUserLoginPost(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.Close()

	logrus.Info("TestUserLoginPost: Starting Excel-driven tests for /login")

	loginTests, err := loadLoginTestData("testdata_login.xlsx", "Sheet1")
	if err != nil {
		t.Fatalf("Error loading login test data: %v", err)
	}

	for _, tt := range loginTests {
		t.Run(tt.Name, func(t *testing.T) {
			logrus.Infof("Running login test case: %q", tt.Name)

			form := url.Values{}
			form.Add("email", tt.Email)
			form.Add("password", tt.Password)
			fmt.Println(form)
			code, _, _ := ts.postForm(t, "/login", form)

			if code != tt.WantCode {
				logrus.Errorf("Login test FAILED for %q: got %d, want %d", tt.Name, code, tt.WantCode)
			} else {
				logrus.Infof("Login test PASSED for %q: got %d (as expected)", tt.Name, code)
			}
			mocks.Equal(t, code, tt.WantCode)
		})
	}
	logrus.Info("TestUserLoginPost: Completed Excel-driven tests for /login")
}

func waitForElement(wd selenium.WebDriver, by, value string, timeout time.Duration) error {
	end := time.Now().Add(timeout)
	for {
		if time.Now().After(end) {
			return fmt.Errorf("timeout waiting for element %s=%s", by, value)
		}
		_, err := wd.FindElement(by, value)
		if err == nil {
			return nil
		}
		time.Sleep(time.Second)
	}
}

func waitForErrorElement(wd selenium.WebDriver, timeout time.Duration) error {
	end := time.Now().Add(timeout)
	for {
		if time.Now().After(end) {
			return fmt.Errorf("timeout waiting for any error message to appear")
		}
		elements, err := wd.FindElements(selenium.ByName, "email")
		if err == nil && len(elements) > 0 {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
}

func TestUserLoginBrowserStack(t *testing.T) {
	logrus.Info("TestUserLoginBrowserStack: Starting BrowserStack E2E tests for /login")

	loginTests, err := loadLoginTestData("testdata_login.xlsx", "Sheet1")
	if err != nil {
		t.Fatalf("Error loading login test data: %v", err)
	}

	bsUser := "cowbuno_7Tam42"
	bsKey := "QJsbG7ySCnDoqzB2tFt9"

	caps := selenium.Capabilities{
		"browserName":     "Chrome",
		"browser_version": "latest",
		"os":              "Windows",
		"os_version":      "10",
	}
	caps["browserstack.user"] = bsUser
	caps["browserstack.key"] = bsKey

	bsHubURL := "http://hub-cloud.browserstack.com/wd/hub"
	wd, err := selenium.NewRemote(caps, bsHubURL)
	if err != nil {
		t.Fatalf("Failed to create remote WebDriver: %v", err)
	}
	defer wd.Quit()

	forumURL := "http://188.227.35.5:8080/login"

	for _, tc := range loginTests {
		t.Run(tc.Name, func(t *testing.T) {
			if err := wd.Get(forumURL); err != nil {
				t.Fatalf("Failed to navigate to login page: %v", err)
			}

			time.Sleep(3 * time.Second)

			emailElem, err := wd.FindElement(selenium.ByName, "email")
			if err != nil {
				t.Fatalf("Failed to find email input: %v", err)
			}
			passwordElem, err := wd.FindElement(selenium.ByName, "password")
			if err != nil {
				t.Fatalf("Failed to find password input: %v", err)
			}
			emailElem.Clear()
			emailElem.SendKeys(tc.Email)
			passwordElem.Clear()
			passwordElem.SendKeys(tc.Password)

			loginButton, err := wd.FindElement(selenium.ByXPATH, "//input[@type='submit' and @value='Login']")
			if err != nil {
				t.Fatalf("Failed to find login button: %v", err)
			}
			if err := loginButton.Click(); err != nil {
				t.Fatalf("Failed to click login button: %v", err)
			}

			if tc.WantCode == http.StatusSeeOther {
				err = waitForElement(wd, selenium.ByID, "user-home", 10*time.Second)
				if err != nil {
					t.Errorf("Expected successful login, but user-home element did not appear: %v", err)
				}
			} else {
				err = waitForErrorElement(wd, 10*time.Second)
				if err != nil {
					t.Errorf("Expected an error message to appear, but it did not: %v", err)
				}
			}
		})
	}

	logrus.Info("TestUserLoginBrowserStack: Completed BrowserStack E2E tests for /login")
}
