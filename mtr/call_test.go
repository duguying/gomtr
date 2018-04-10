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

	time.Sleep(time.Minute)
}

//var report = `Start: Tue Apr 10 14:02:03 2018
//HOST: Test-91                     Loss%   Snt   Last   Avg  Best  Wrst StDev
//  1.|-- gateway                    0.0%    10    0.3   0.3   0.2   0.4   0.0
//  2.|-- 218.17.157.193             0.0%    10    4.1  28.5   3.2  49.7  16.5
//  3.|-- 10.1.100.73                0.0%    10   12.1  26.2   2.1  56.8  19.9
//  4.|-- 218.17.52.42              10.0%    10   28.3  30.5   1.5  84.7  30.2
//  5.|-- 119.145.47.126            10.0%    10   48.7  41.5   1.8  99.5  34.5
//  6.|-- 219.133.30.234             0.0%    10   25.1  40.1   2.1 102.4  35.7
//  7.|-- ???                       100.0    10    0.0   0.0   0.0   0.0   0.0
//  8.|-- 183.56.65.2                0.0%    10    9.8  44.1   9.8  99.3  27.9
//  9.|-- 202.97.42.102             40.0%    10   33.3  57.1  33.3  87.6  20.6
// 10.|-- 150.138.128.130           90.0%    10   39.9  39.9  39.9  39.9   0.0
// 11.|-- ???                       100.0    10    0.0   0.0   0.0   0.0   0.0
// 12.|-- ???                       100.0    10    0.0   0.0   0.0   0.0   0.0
// 13.|-- 119.38.215.29              0.0%    10   56.1  54.6  40.6  79.2  13.6
// 14.|-- ???                       100.0    10    0.0   0.0   0.0   0.0   0.0
// 15.|-- 114.215.151.25             0.0%    10   42.1  60.3  41.0  96.2  17.5`
//
//func TestParse(t *testing.T) {
//	c:=New("mtr")
//	c.parseReport(report)
//}