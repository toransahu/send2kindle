package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	user2 "os/user"
	"path"
	"strconv"
	"strings"

	"github.com/toransahu/send2kindle/util"
)

type config struct {
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	StorePath string `json:"storepath"`
	Password  string `json:"password"`
	Server    string `json:"server"`
	Port      int    `json:"port"`
}

const DefaultTimeout = 120
const XdgConfigHome = "XDG_CONFIG_HOME"
const ConfigFolderName = "send2kindle"

var instance *config

func isGmail(mail string) bool {
	//use regex here
	return strings.Contains(mail, "@gmail.com")
}
func DefaultConfigPath() (string, error) {
	user, err := user2.Current()
	if err != nil {
		util.Red.Println("Couldn't get current user ", err)
		os.Exit(1)
	}
	xdgConfigHome := os.Getenv(XdgConfigHome)
	var configFolder string
	if len(xdgConfigHome) == 0 {
		configFolder = path.Join(user.HomeDir, ".config")
		configFolder = path.Join(configFolder, ConfigFolderName)
		//util.Cyan.Println("Config home not set, will look for config at ", configFolder)
	} else {
		configFolder = path.Join(xdgConfigHome, ConfigFolderName)
	}
	_ = os.Mkdir(configFolder, os.ModePerm)

	return path.Join(configFolder, "KindleConfig.json"), nil
}
func exists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		util.Red.Println(err)
		return false
	}
	return true
}
func NewConfig() *config {
	config := config{}
	config.Server = "smtp.gmail.com"
	config.Port = 465
	return &config
}
func CreateConfig() *config {

	util.CyanBold.Println("CONFIGURE KINDLE-SEND")

	configuration := NewConfig()
	util.Cyan.Printf("Email of your device and press enter (eg. purple_terminal@kindle.com) : ")
	configuration.Receiver = util.ScanlineTrim()
	util.Cyan.Printf("Email that'll be used to send documents to device (eg. yourname@gmail.com) : ")
	configuration.Sender = util.ScanlineTrim()

	if !isGmail(configuration.Sender) {
		util.Cyan.Println("Sender email is different then Gmail, " +
			"can you help with SMTP server address and SMTP port for your email provider\n" +
			"Just search SMTP settings for <your email domain>.com on internet")

		util.Cyan.Printf("Enter SMTP Server Address (eg. smtp.gmail.com) : ")
		configuration.Server = util.ScanlineTrim()
		for {
			util.Cyan.Printf("Enter SMTP port (usually 587 or 465) : ")
			portStr := util.ScanlineTrim()
			portInt, err := strconv.Atoi(portStr)
			if err != nil {
				util.Red.Println("Entered port number is either invalid or not an integer, please try again")
				continue
			}
			configuration.Port = portInt
			break
		}
	}

	util.Cyan.Printf("Enter password for Sender %s (password remains encrypted in your machine) : ", configuration.Sender)
	configuration.Password = util.ScanlineTrim()

	util.Cyan.Printf("File path to store all the documents on your computer (empty is ok) :")
	configuration.StorePath = util.ScanlineTrim()
	encryptedPass, err := Encrypt(configuration.Sender, configuration.Password)
	if err != nil {
		util.Red.Println("Error encrypting password ", err)
		os.Exit(1)
	}
	configuration.Password = encryptedPass

	return configuration
}

func handleCreation(filename string) error {
	util.Red.Println("Configuration file doesn't exist\n Answer next few questions to create config file")
	configuration := CreateConfig()
	err := Save(*configuration, filename)
	if err != nil {
		util.Red.Println("Error while writing config to ", filename, err)
		return err
	}
	util.Green.Printf("Config created successfully and stored at %s, you can directly edit it later on \n", filename)
	return nil
}
func Load(filename string) (config, error) {
	if !exists(filename) {
		err := handleCreation(filename)
		if err != nil {
			return config{}, err
		}
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		util.Red.Println("Error reading config ", err)
		return config{}, err
	}
	var c config
	err = json.Unmarshal(data, &c)
	if err != nil {
		util.Red.Println("Error converting config to json ", err)
		return config{}, err
	}
	decryptedPass, err := Decrypt(c.Sender, c.Password)
	if err != nil {
		util.Red.Println("Error decrypting password : ", err)
		os.Exit(1)
	}
	c.Password = decryptedPass
	InitializeConfig(&c)
	return c, nil
}

func Save(c config, filename string) error {
	data, err := json.MarshalIndent(c, "", "	")
	if err != nil {
		util.Red.Println("Error parsing configuration for writing")
		return err
	}
	return ioutil.WriteFile(filename, data, 0644)
}

func InitializeConfig(c *config) {
	if instance == nil {
		instance = c
		util.Green.Println("Loaded configuration")
	}
}

func GetInstance() *config {
	return instance
}
