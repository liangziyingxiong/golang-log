package conf

import (
	"os"
	"fmt"
	"bufio"
	"unicode"
	"strings"
	"logs"
)

type Conf struct {
	file     string
	loger    *logs.Log

	k_v      map[string]string
}

func New(cfile string, loger *logs.Log) *Conf {
	f, err := os.Open(cfile)
	if err != nil {
		loger.Error("%v\n", err)
		os.Exit(-1)
	}
	defer f.Close()

	_conf := &Conf{file: cfile, loger: loger}
	_conf.mapKeyValue(f)
	return _conf
}

func (c *Conf) mapKeyValue(f *os.File) {
	r := bufio.NewReader(f)

	c.k_v = make(map[string]string)
	for {
		var key_end, value_start, value_end, quote_start, quote_end, value_char_start int
		var key, value string
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimLeftFunc(line, unicode.IsSpace)
		if strings.HasPrefix(line, "#") {
			continue
		}
		if y := strings.Index(line, "#"); y != -1 {
			line = line[0:y]
		}
		for k, v := range line {
			switch v {
				case ' ':
					if key_end == 0 {
						key = line[:k]
						key = strings.TrimSpace(key)
						key_end = k
					} else if value_end == 0 && value_char_start != 0 {
						value_end = k
						value = line[value_start:value_end]
						value = strings.TrimSpace(value)
					}
				case '=':
					if value_start == 0 {
						value_start = k + 1
					}
				case '"':
					if quote_start == 0 {
						quote_start = k
						value_start = quote_start + 1
					} else if quote_end == 0 {
						quote_end = k
						value_end = quote_end
						value = line[value_start:value_end]
					}
				case '#', '\n', '\r':
					if value_start == 0 {
						fmt.Fprintf(os.Stderr, "%s is invalid format\n", line)
						os.Exit(-1)
					}
					if quote_start != 0 && quote_end == 0 {
						fmt.Fprintf(os.Stderr, "%s need end quote\n", line)
						os.Exit(-1)
					}
					if value_end == 0 {
						value_end = k
						value = line[value_start:value_end]
						value = strings.TrimSpace(value)
					}
					break
				default:
					if value_start != 0 && value_char_start == 0{
						value_start = k
						value_char_start = k
					}
			}
		}
		c.k_v[key] = value
	}
	c.loger.Printf("c.k_v : %q\n", c.k_v)
}

func (c *Conf) GetValueByKey(key string) string {
	if val, ok := c.k_v[key]; ok == true {
		return val
	}
	return ""
}
