package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/robertkrimen/otto"

	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"path"

	"crypto/sha256"
	"io/ioutil"
	"net/http"
)

const FILTER_BASE string = "http://192.241.174.50/filters/"
const REQUIREJS_URL string = "http://requirejs.org/docs/release/2.1.22/r.js"

type filterTuple struct {
	fromObjType string
	toObjType   string
}

func listScripts() ([]string, error) {
	var results []string

	path, err := getCachePath("")
	if err != nil {
		return nil, err
	}
	fi, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range fi {
		results = append(results, file.Name())
	}
	return results, nil
}

func exportScript(pipeline string, name string) {
	/*
			// Make directory...

			var cfgf = io.FileWriter //>>iopipe-cfg.json:
			var scriptf = io.Filewriter //>index.js
			var packagef = io.Filewriter //>package.json
			cfgf.write(`{"auth": {}}`)

			scriptf.write(`var iopipe = require("iopipe")
		        var config = require("./iopipe-cfg.json")
			iopipe.load_config(config)
			exports.run = function() {
		          iopipe.exec()
		        }`)

			packagef.write(`
			{
			  "name": "iopipe",
			  "private": true,
			  "version": "0.0.1",
			  "description": "iopipe sdk",
			  "author": "Eric Windisch",
			  "dependencies": {
			    "read-stream": "",
			    "request": ""
			  },
			  "main": "./iopipe.js"
			}
			`)
	*/
}

func makeFilter(script string) (func(input string) (string, error), error) {
	var (
		res  *http.Response
		body []byte
		err  error
	)

	res, err = http.Get(REQUIREJS_URL)
	if err != nil {
		log.Fatal(err)
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	res.Body.Close()
	rjs := string(body[:])

	return func(input string) (string, error) {
		vm := otto.New()

		logrus.Debug("Adding RequireJS")
		vm.Run(rjs)

		logrus.Debug("Executing script: " + script)
		vm.Set("msg", input)
		_, err := vm.Run(`
                        var module = { "exports": function() { } }
                `)
		if err != nil {
			return "", err
		}
		_, err = vm.Run(script)
		if err != nil {
			return "", err
		}
		val, err := vm.Run("module.exports(msg)")
		if err != nil {
			return "", err
		}
		return val.ToString()
	}, nil
}

func fetchFilter(filterPath string) ([]byte, error) {
	var (
		res  *http.Response
		body []byte
		err  error
	)

	path := path.Join(FILTER_BASE, filterPath)
	res, err = http.Get(path)
	if err != nil {
		return nil, err
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	/* Verify digest */
	chksum := sha256.Sum256(body[:])
	if filterPath != string(chksum[:]) {
		return nil, errors.New("Checksum failure")
	}

	return body, nil
}

func importScript(file string) (string, error) {
	var err error
	var fH io.Reader

	if file == "-" {
		fH = os.Stdin
	} else {
		fH, err = os.Open(file)
		if err != nil {
			return "", err
		}
	}
	reader := bufio.NewReader(fH)
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	id, err := writeFilterCache(body[:])
	if err != nil {
		return id, err
	}
	return id, nil
}

func getFilter(filterPath string) (func(input string) (string, error), error) {
	var script []byte
	var err error

	/* Do we have this cached? */
	if script, err := readFilterCache(filterPath); err == nil {
		return makeFilter(string(script[:]))
	} else {
		return nil, err
	}

	/* If not, fetch */
	if script, err = fetchFilter(filterPath); err != nil {
		return nil, err
	}
	if _, err = writeFilterCache(script); err != nil {
		return nil, err
	}

	return makeFilter(string(script[:]))
}

func getPipeline(filterPath string) (func(input string) (string, error), error) {
	var script []byte
	var err error

	diskPath, err := getCachePath("")
	if err != nil {
		return nil, err
	}

	/* Do we have this cached? */
	if _, err := os.Stat(diskPath); err == nil {
		script, err = ioutil.ReadFile(diskPath)
		return makeFilter(string(script[:]))
	}

	/* If not, fetch */
	if script, err = fetchFilter(filterPath); err != nil {
		return nil, err
	}
	/* Write cache */
	if _, err = writeFilterCache(script); err != nil {
		return nil, err
	}

	return makeFilter(string(script[:]))
}

func findFilters(fromObjType string, toObjType string) (string, error) {
	var (
		res  *http.Response
		body []byte
		err  error
	)

	path := path.Join(FILTER_BASE, fromObjType, toObjType)
	res, err = http.Get(path)
	if err != nil {
		return "", err
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	res.Body.Close()
	response := string(body[:])
	return response, nil
}

func findPipelines(fromObjType string, toObjType string) (string, error) {
	var (
		res  *http.Response
		body []byte
		err  error
	)

	path := path.Join(FILTER_BASE, fromObjType, toObjType)
	res, err = http.Get(path)
	if err != nil {
		return "", err
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	res.Body.Close()
	response := string(body[:])
	return response, nil
}

func publishPipeline(pipelineName string) {
}

func subscribePipeline(pipelineName string) {
}

func createPipeline(pipeparts []string) (string, error) {
	return "", nil
}

func removeFilter(filterid string) error {
	path, err := getCachePath(filterid)
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func tagPipeline(pipeline string) error {
	return nil
}

func tagFilter(filterid string, name string) error {
	path, err := getCachePath(filterid)
	if err != nil {
		return err
	}
	destpath, err := getCachePath(name)
	if err != nil {
		return err
	}
	err = os.Symlink(path, destpath)
	return err
}
