package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const DIRECTORY = "logs"

type Config struct {
	directory string
	filetype  string
}

func Init(filetype string) (Config, error) {
	var config Config
	var filename string

	if _, err := os.Stat(DIRECTORY); os.IsNotExist(err) {
		err = os.MkdirAll(DIRECTORY, 0755)

		if err != nil {
			return config, err
		}
	}

	switch filetype {
	case "log":
		filename = "error.log"
	case "json":
		filename = "error.json"
	default:
		return config, errors.New("adapter not defined")
	}

	path := fmt.Sprintf("%s/%s", DIRECTORY, filename)

	var _, err = os.Stat(path)

	if os.IsNotExist(err) {
		file, err := os.Create(path)

		if err != nil {
			fmt.Println(err)
		}

		defer file.Close()
	}

	return Config{
		directory: path,
		filetype:  filetype,
	}, nil
}

func (c *Config) Log(err error) {
	switch c.filetype {
	case "log":
		_ = c.logFile(err)
	case "json":
		_ = c.logJson(err)
	default:
		return
	}
}

func (c *Config) logFile(e error) error {
	path := c.directory

	var file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	newError := fmt.Sprintf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), e.Error())

	_, err = fmt.Fprintln(file, newError)
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) logJson(e error) error {
	path := c.directory

	var file, err = ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var jsonData []map[string]interface{}

	_ = json.Unmarshal(file, &jsonData)
	newError := map[string]interface{}{
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
		"error":     e.Error(),
	}
	jsonData = append(jsonData, newError)
	jsonString, err := json.Marshal(jsonData)

	if err != nil {
		return err
	}

	_ = ioutil.WriteFile(path, jsonString, os.ModePerm)

	return nil
}
