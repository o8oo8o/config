package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	fileName string
	comment  []string
}

type Section map[string]string

func (s Section) GetInt(key string) (int, error) {
	data := 0
	if val, ok := s[key]; ok {
		data, err := strconv.Atoi(val)
		if err == nil {
			return data, nil
		}
	}
	return data, errors.New("GetInt Error")
}

func (s Section) GetFloat(key string) (float64, error) {
	data := 0.0
	if val, ok := s[key]; ok {
		data, err := strconv.ParseFloat(val, 64)
		if err == nil {
			return data, nil
		}
	}
	return data, errors.New("GetFloat Error")
}

func (s Section) GetString(key string) (string, error) {
	if val, ok := s[key]; ok {
		return val, nil
	}
	return "", errors.New("GetString Error")
}

func (s Section) GetBool(key string) (bool, error) {
	if val, ok := s[key]; ok {
		fmt.Println(val)
		data, err := strconv.ParseBool(val)
		if err == nil {
			return data, nil
		}
	}
	return false, errors.New("GetBool Error")
}

// 读取配置文件的每一行
func (c Config) ReadLines() (lines []string, err error) {
	fd, err := os.Open(c.fileName)
	if err != nil {
		return
	}
	defer fd.Close()
	lines = make([]string, 0)
	reader := bufio.NewReader(fd)
	prefix := ""
	var isLongLine bool
	for {
		byteLine, isPrefix, er := reader.ReadLine()
		if er != nil && er != io.EOF {
			return nil, er
		}
		if er == io.EOF {
			break
		}
		line := string(byteLine)
		if isPrefix {
			prefix += line
			continue
		} else {
			isLongLine = true
		}

		line = prefix + line
		if isLongLine {
			prefix = ""
		}
		line = strings.TrimSpace(line)
		// 跳过空白行
		if len(line) == 0 {
			continue
		}
		// 跳过注释行
		var breakLine = false
		for _, v := range c.comment {
			if strings.HasPrefix(line, v) {
				breakLine = true
				break
			}
		}
		if breakLine {
			continue
		}

		lines = append(lines, line)
	}
	return lines, nil
}

// 获取所有配置
func (c Config) GetAllConfig() map[string]map[string]string {
	allConfig := make(map[string]map[string]string)
	lines, err := c.ReadLines()
	if err != nil {
		log.Fatalln(err)
	}
	var section = make(map[string]string, 1)

	for _, line := range lines {
		if line[0] == '[' && line[len(line)-1] == ']' {
			sectionName := line[1 : len(line)-1]
			section = make(map[string]string, 1)
			allConfig[sectionName] = section
		} else {
			configKeyVal := strings.Split(line, "=")
			key := strings.TrimSpace(configKeyVal[0])
			val := strings.TrimSpace(strings.Join(configKeyVal[1:], "="))
			section[key] = val
		}
	}
	return allConfig
}

// 获取某一段配置
func (c Config) GetSection(section string) (Section, error) {
	if data, ok := c.GetAllConfig()[section]; ok {
		return data, nil
	}
	return map[string]string{}, nil
}

func LoadConfigFile(filename string, comment []string) (Config, error) {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("file not exist:", err)
			return Config{}, err
		}
	}
	return Config{
		fileName: filename,
		comment:  comment,
	}, nil
}
