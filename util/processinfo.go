// Copyright 2017 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"github.com/pingcap/parser/ast"
	"github.com/sirupsen/logrus"
	"reflect"
	"strings"
	"time"
)

// ProcessInfo is a struct used for show processlist statement.
type ProcessInfo struct {
	ID      uint64
	User    string
	Host    string
	DB      string
	Command string
	Time    time.Time
	State   uint16
	Info    string
	Mem     int64
}

// SessionManager is an interface for session manage. Show processlist and
// kill statement rely on this interface.
type SessionManager interface {
	// ShowProcessList returns map[connectionID]ProcessInfo
	ShowProcessList() map[uint64]ProcessInfo
	Kill(connectionID uint64, query bool)
}

type AstVisitor struct {
	Level *int
}

func NewAstVisitor() AstVisitor {
	var i int
	return AstVisitor{Level: &i}
}

func (a AstVisitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
	var p []string
	for i := 0; i < *a.Level; i++ {
		p = append(p, " -- ")
	}

	logrus.Infof("ast enter: %s%s", strings.Join(p, ""), reflect.TypeOf(in))
	switch tt := in.(type) {
	case *ast.BinaryOperationExpr:
		logrus.Infof("left: %s ,op:%s, right: %s", reflect.TypeOf(tt.L), tt.Op, reflect.TypeOf(tt.R))
	case *ast.VariableExpr:
		logrus.Infof("VariableExpr__ name: %s value: %s", tt.Name, reflect.TypeOf(tt.Value))
	case *ast.ColumnNameExpr:

	}

	*a.Level += 1
	return in, false
}

func (a AstVisitor) Leave(in ast.Node) (out ast.Node, ok bool) {
	*a.Level -= 1
	var p []string
	for i := 0; i < *a.Level; i++ {
		p = append(p, " -- ")
	}
	//logrus.Infof("ast leave: %s%s", strings.Join(p, ""), reflect.TypeOf(in))

	return in, true
}

func Accept(n ast.Node) {
	var l int
	v := AstVisitor{Level: &l}
	logrus.Info("------------------ show ast")
	n.Accept(v)
}
