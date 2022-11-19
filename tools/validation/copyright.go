/*
Copyright 2021 Flant JSC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var EELicenseRe = regexp.MustCompile(`(?s)Copyright 202[1-9] Flant JSC.*Licensed under the Deckhouse Platform Enterprise Edition \(EE\) license.*See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE`)

var CELicenseRe = regexp.MustCompile(`(?s)[/#{!-]*(\s)*Copyright 202[1-9] Flant JSC[-!}\n#/]*
[/#{!-]*(\s)*Licensed under the Apache License, Version 2.0 \(the \"License\"\);[-!}\n]*
[/#{!-]*(\s)*you may not use this file except in compliance with the License.[-!}\n]*
[/#{!-]*(\s)*You may obtain a copy of the License at[-!}\n#/]*
[/#{!-]*(\s)*http://www.apache.org/licenses/LICENSE-2.0[-!}\n#/]*
[/#{!-]*(\s)*Unless required by applicable law or agreed to in writing, software[-!}\n]*
[/#{!-]*(\s)*distributed under the License is distributed on an \"AS IS\" BASIS,[-!}\n]*
[/#{!-]*(\s)*WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.[-!}\n]*
[/#{!-]*(\s)*See the License for the specific language governing permissions and[-!}\n]*
[/#{!-]*(\s)*limitations under the License.[-!}\n]*`)

var fileToCheckRe = regexp.MustCompile(`\.go$|/[^/.]+$|\.sh$|\.lua$|\.py$`)
var fileToSkipRe = regexp.MustCompile(`geohash.lua$|\.github/CODEOWNERS|Dockerfile$|Makefile$|/docs/documentation/|/docs/site/|bashrc$|inputrc$|modules_menu_skip$`)

func RunCopyrightValidation(info *DiffInfo) (exitCode int) {
	fmt.Printf("Run 'copyright' validation ...\n")

	if len(info.Files) == 0 {
		fmt.Printf("Nothing to validate, diff is empty\n")
		os.Exit(0)
	}

	exitCode = 0
	msgs := NewMessages()
	for _, fileInfo := range info.Files {
		if !fileInfo.HasContent() {
			continue
		}

		fileName := fileInfo.NewFileName

		if fileToCheckRe.MatchString(fileName) && !fileToSkipRe.MatchString(fileName) {
			msgs.Add(checkFileCopyright(fileName))
		} else {
			msgs.Add(NewSkip(fileName, ""))
		}
	}
	msgs.PrintReport()

	if msgs.CountErrors() > 0 {
		exitCode = 1
	}

	return exitCode
}

var copyrightOrAutogenRe = regexp.MustCompile(`Copyright The|autogenerated|DO NOT EDIT`)
var copyrightRe = regexp.MustCompile(`Copyright`)
var flantRe = regexp.MustCompile(`Flant|Deckhouse`)

// checkFileCopyright returns true if file is readable and has no copyright information in it.
func checkFileCopyright(fName string) Message {
	// Original script 'validate_copyright.sh' used 'head -n 10'.
	// Here we just read first 1024 bytes.
	headBuf, err := readFileHead(fName, 1024)
	if err != nil {
		return NewSkip(fName, err.Error())
	}

	// Skip autogenerated file or file already has other than Flant copyright
	if copyrightOrAutogenRe.Match(headBuf) {
		return NewSkip(fName, "generated code or other license")
	}

	// Check Flant license if file contains keywords.
	if flantRe.Match(headBuf) {
		return checkFlantLicense(fName, headBuf)
	}

	// Skip file with some other copyright
	if copyrightRe.Match(headBuf) {
		return NewSkip(fName, "contains other license")
	}

	return NewError(fName, "no copyright or license information", "")
}

func checkFlantLicense(fName string, buf []byte) Message {
	if strings.HasPrefix(fName, "/ee/") || strings.HasPrefix(fName, "ee/") {
		if !EELicenseRe.Match(buf) {
			return NewError(fName, "EE related file should contain EE license", "")
		}
	} else {
		if !CELicenseRe.Match(buf) {
			return NewError(fName, "should contain CE license", "")
		}
	}

	return NewOK(fName)
}

func readFileHead(fName string, size int) ([]byte, error) {
	file, err := os.Open(fName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return nil, fmt.Errorf("directory")
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("symlink")
	}

	headBuf := make([]byte, size)
	_, err = file.Read(headBuf)
	if err != nil {
		return nil, err
	}

	return headBuf, nil
}
