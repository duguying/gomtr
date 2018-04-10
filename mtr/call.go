// Copyright 2018. All rights reserved.
// This file is part of gomtr project
// Created by duguying on 2018/4/10.

package mtr

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type ToolMtr struct {
	ToolPath string `json:"tool_path"`
}

func New(path string) *ToolMtr {
	return &ToolMtr{
		ToolPath: path,
	}
}

func (tm *ToolMtr) SimpleCall(host string, size int, timeout time.Duration) (err error) {
	content, err := tm.call(60)
	if err != nil {
		return err
	}
	fmt.Println("=======>",content)
	tm.parseReport(content)
	return nil
}

func (tm *ToolMtr) call(size int) (content string, err error) {
	c := exec.Command(tm.ToolPath, "-s", fmt.Sprintf("%d", size), "-n", "-r")
	data, err := c.CombinedOutput()
	if err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}

func (tm *ToolMtr) parseReport(content string) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.Replace(line, "\r", "", -1)
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "Start:") {
			// 时间行
		} else if strings.HasPrefix(line, "HOST:") {
			// 表头
		} else {
			line = strings.Replace(line, "\t", " ", -1)
			exp, _ := regexp.Compile(`[ ]+`)
			line = exp.ReplaceAllString(line, " ")
			fmt.Println(line)
		}
	}
}
