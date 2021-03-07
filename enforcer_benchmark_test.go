/*
Copyright 2020 RS4
@Author: Weny Xu
@Date: 2021/03/08 0:53
*/

package casbin

import (
	"sync"
	"testing"
)

//
func (e *Enforcer) BatchEnforceBasic(requests [][]interface{}) ([]bool, error) {
	var results []bool
	for _, request := range requests {
		result, err := e.enforce("", nil, request...)
		if err != nil{
			return results, err
		}
		results = append(results, result)
	}
	return results, nil
}

func testBatchEnforceBenchmark(t *testing.B, e *Enforcer, requests [][]interface{}, results []bool) {
	t.Helper()
	myRes, _ := e.BatchEnforce(requests)
	if len(myRes) != len(results) {
		t.Errorf("%v supposed to be %v", myRes, results)
	}
	for i, v := range myRes {
		if v != results[i] {
			t.Errorf("%v supposed to be %v", myRes, results)
		}
	}
}

func testBatchEnforceBenchmarkBasic(t *testing.B, e *Enforcer, requests [][]interface{}, results []bool) {
	t.Helper()
	myRes, _ := e.BatchEnforceBasic(requests)
	if len(myRes) != len(results) {
		t.Errorf("%v supposed to be %v", myRes, results)
	}
	for i, v := range myRes {
		if v != results[i] {
			t.Errorf("%v supposed to be %v", myRes, results)
		}
	}
}

func BenchmarkEnforcer_BatchEnforceBasicSync(b *testing.B){
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		testBatchEnforceBenchmarkBasic(b, e, [][]interface{}{
			{"alice", "data1", "write"},
			{"alice", "data1", "read"},
			{"bob", "data2", "read"},
			{"bob", "data2", "write"},
			{"data1_deny_group", "data1", "read"},
			{"data1_deny_group", "data1", "write"},
			{"data2_allow_group", "data2", "read"},
			{"data2_allow_group", "data2", "write"},
			{"alice", "data1", "write"},
			{"alice", "data1", "read"},
			{"bob", "data2", "read"},
			{"bob", "data2", "write"},
			{"data1_deny_group", "data1", "read"},
			{"data1_deny_group", "data1", "write"},
			{"data2_allow_group", "data2", "read"},
			{"data2_allow_group", "data2", "write"},
		}, []bool{
			true, true, false, true, false, false, true, true,true, true, false, true, false, false, true, true,
		})
	}
}

func BenchmarkEnforcer_BatchEnforceSync(b *testing.B) {
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	b.ResetTimer()
	for i:=0;i<b.N;i++{
		testBatchEnforceBenchmark(b, e, [][]interface{}{
			{"alice", "data1", "write"},
			{"alice", "data1", "read"},
			{"bob", "data2", "read"},
			{"bob", "data2", "write"},
			{"data1_deny_group", "data1", "read"},
			{"data1_deny_group", "data1", "write"},
			{"data2_allow_group", "data2", "read"},
			{"data2_allow_group", "data2", "write"},
			{"alice", "data1", "write"},
			{"alice", "data1", "read"},
			{"bob", "data2", "read"},
			{"bob", "data2", "write"},
			{"data1_deny_group", "data1", "read"},
			{"data1_deny_group", "data1", "write"},
			{"data2_allow_group", "data2", "read"},
			{"data2_allow_group", "data2", "write"},
		}, []bool{
			true, true, false, true, false, false, true, true,true, true, false, true, false, false, true, true,
		})
	}
}

func BenchmarkEnforcer_BatchEnforce(b *testing.B) {
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	b.ResetTimer()
	var wg sync.WaitGroup
	for i:=0;i<b.N;i++{
		wg.Add(1)
		go func() {
			defer wg.Done()
			testBatchEnforceBenchmark(b, e, [][]interface{}{
				{"alice", "data1", "write"},
				{"alice", "data1", "read"},
				{"bob", "data2", "read"},
				{"bob", "data2", "write"},
				{"data1_deny_group", "data1", "read"},
				{"data1_deny_group", "data1", "write"},
				{"data2_allow_group", "data2", "read"},
				{"data2_allow_group", "data2", "write"},
				{"alice", "data1", "write"},
				{"alice", "data1", "read"},
				{"bob", "data2", "read"},
				{"bob", "data2", "write"},
				{"data1_deny_group", "data1", "read"},
				{"data1_deny_group", "data1", "write"},
				{"data2_allow_group", "data2", "read"},
				{"data2_allow_group", "data2", "write"},
			}, []bool{
				true, true, false, true, false, false, true, true,true, true, false, true, false, false, true, true,
			})
		}()
	}
	wg.Wait()
}
