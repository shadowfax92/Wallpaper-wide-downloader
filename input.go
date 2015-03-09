package main

import (
	"bufio"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Data structures for storing input config and Urls
type inputConfig struct {
	request_rate          uint64
	duration              time.Duration
	max_cpu               uint64
	workers               uint64
	url_file              string
	data_collector_script string
}

type urlData struct {
	request_type string //GET or POST

	url string
}

var (
	errorFileRead = errors.New("Error while opening input file")
	errorYamlRead = errors.New("Error while reading yaml data")
)

//TODO add logging for each of the error
func ReadInputFromYaml(file_path string) (inputConfig, error) {
	var input_data inputConfig

	defer func() {
		err := recover()
		if err != nil {
			Error.Fatalln("Function paniked with error = ", err)
		}
	}()

	byte_string, err := ioutil.ReadFile(file_path)

	if err != nil {
		return input_data, errorFileRead
	}

	m := make(map[interface{}]interface{})

	err = yaml.Unmarshal(byte_string, &m)

	if err != nil {
		return input_data, errorYamlRead
	}

	for key, val := range m {
		if key == "request_rate" {
			input_data.request_rate, err = strconv.ParseUint(val.(string), 0, 64)
		} else if key == "duration" {
			input_data.duration, err = time.ParseDuration(val.(string))
		} else if key == "max_cpus" {
			input_data.max_cpu, err = strconv.ParseUint(val.(string), 0, 64) //strconv.Atoi(val.(string))
		} else if key == "workers" {
			input_data.workers, err = strconv.ParseUint(val.(string), 0, 64)
		} else if key == "urls_file" {
			input_data.url_file = val.(string)
		} else if key == "data_collector_command" {
			input_data.data_collector_script = val.(string)
		}
	}

	if err != nil {
		return input_data, errorYamlRead
	}

	return input_data, err
}

func ReadUrlFile(file_path string) ([]urlData, error) {
	url_data := make([]urlData, 100) // slices automatically re-size
	inputFile, err := os.Open(file_path)

	if err != nil {
		log.Fatal("Error opening input file:", errorFileRead)
		return url_data, errorFileRead
	}
	defer inputFile.Close()

	url_regex := regexp.MustCompile(`(http)\S+`)
	request_regex := regexp.MustCompile(`GET|POST`)
	scanner := bufio.NewScanner(inputFile)

	// scanner.Scan() advances to the next token returning false if an error was encountered
	for scanner.Scan() {
		line := scanner.Text()

		var tmp_url urlData
		for _, val := range strings.Fields(line) {
			if url_regex.MatchString(val) {
				tmp_url.url = val
			}
			if request_regex.MatchString(val) {
				tmp_url.request_type = val
			}

		}
		url_data = append(url_data, tmp_url)
	}

	for _, val := range url_data {
		fmt.Println(val.url)
	}

	// When finished scanning if any error other than io.EOF occured
	// it will be returned by scanner.Err().
	if err := scanner.Err(); err != nil {
		log.Fatal(scanner.Err())
		return url_data, err
	}

	return url_data, nil
}
