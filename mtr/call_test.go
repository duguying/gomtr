// Copyright 2018. All rights reserved.
// This file is part of gomtr project
// Created by duguying on 2018/4/10.

package mtr

import (
	"fmt"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	c := New("mtr")
	err := c.SimpleCall("duguying.net", 60, time.Minute)
	if err != nil {
		fmt.Println(err.Error())
	}
}
