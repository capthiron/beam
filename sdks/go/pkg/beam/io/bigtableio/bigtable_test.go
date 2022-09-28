// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bigtableio

import (
	"errors"
	"reflect"
	"testing"
)

func TestHashStringToInt(t *testing.T) {
	equalVal := "equal"
	
	firstHash := hashStringToInt(equalVal)
	secondHash := hashStringToInt(equalVal)

	if firstHash != secondHash {
		t.Errorf("hashStringToInt(\"%s\") should equal hashStringToInt(\"%s\")", equalVal, equalVal)
	}

	if hashStringToInt("helloWorld") == hashStringToInt("helloworld") {
		t.Error("hashStringToInt(\"helloWorld\") should not equal hashStringToInt(\"helloworld\")")
	}

	if hashStringToInt("saturnsmoon1") == hashStringToInt("saturnsmoon2") {
		t.Error("hashStringToInt(\"saturnsmoon1\") should not equal hashStringToInt(\"saturnsmoon2\")")
	}
}

func TestAddGroupKeyFnGroupKeyGiven(t *testing.T) {
	bmWithGroupKey := NewMutation("rowKey").WithGroupKey("1")
	groupKey, _ := addGroupKeyFn(*bmWithGroupKey)
	if groupKey == 1 {
		t.Error("addGroupKeyFn should hash groupKey values properly, but projected \"1\" -> 1")
	}
}

func TestAddGroupKeyFnNoGroupKeyGiven(t *testing.T) {
	bmNoGroupKey := NewMutation("rowKey")
	groupKey, _ := addGroupKeyFn(*bmNoGroupKey)
	if groupKey != 1 {
		t.Errorf("addGroupKeyFn should assign 1 as hash if no groupKey is given, but projected nil -> %d", groupKey)
	}
}

func TestMustBeBigtableMutation(t *testing.T) {

	mutation := NewMutation("key")
	mutation.Set("family", "column", 0, []byte{})

	mutationWithGroupKey := NewMutation("key").WithGroupKey("groupKey")
	mutationWithGroupKey.Set("family", "column", 0, []byte{})

	passValues := []BigtableioMutation {
		{},
		{RowKey: "key"},
		{Mutations: []Mutation {{}}},
		*mutation,
		*mutationWithGroupKey,
	}

	for _, passValue := range passValues {
		passType := reflect.TypeOf(passValue)
		err := mustBeBigtableioMutation(passType)
		if err != nil {
			t.Errorf("input type %v should be considered a bigtableio.BigtableMutation", passType)
		}
	}
}

func TestMustNotBeBigtableMutation(t *testing.T) {
	failValues := []interface{} {
		1,
		1.0,
		"strings must fail",
		errors.New("errors must fail"),
	}

	for _, failValue := range failValues {
		failType := reflect.TypeOf(reflect.ValueOf(failValue))
		err := mustBeBigtableioMutation(failType)
		if err == nil {
			t.Errorf("input type %v should not be considered a bigtableio.BigtableMutation", failType)
		}
	}
}

func TestTryApplyBulk(t *testing.T) {
	err := tryApplyBulk(nil, nil)
	if err != nil {
		t.Error("tryApplyBulk should not return an error for inputs <nil, nil> but returned:\n", err)
	}

	err = tryApplyBulk(nil, errors.New("error"))
	if err == nil {
		t.Error("tryApplyBulk should return an error for inputs <nil, error>")
	}
	
	err = tryApplyBulk([]error { errors.New("error") }, nil)
	if err == nil {
		t.Error("tryApplyBulk should return an error for inputs <[]error, nil>")
	}
}

