/*
Copyright 2020 RS4
@Author: Weny Xu
@Date: 2021/03/08 0:53
*/

package casbin

import (
	"fmt"
	"sync"
	"testing"
)

var (
	InputData = [][]interface{}{
		{"alice", "data1", "write"},
		{"alice", "data1", "read"},
		{"bob", "data2", "read"},
		{"bob", "data2", "write"},
		{"data1_deny_group", "data1", "read"},
		{"data1_deny_group", "data1", "write"},
		{"data2_allow_group", "data2", "read"},
		{"data2_allow_group", "data2", "write"},
	}
	Result = []bool{
		true, true, false, true, false, false, true, true,
	}
)

func testSetGen(scale int) (input [][]interface{}, result []bool) {
	for i := 0; i != scale; i++ {
		input = append(input, InputData...)
		result = append(result, Result...)
	}
	return
}

// BatchEnforce enforce in batches
func (e *Enforcer) BatchEnforceV2(requests [][]interface{}) ([]bool, error) {
	var results = make([]bool, len(requests))
	var mu1 sync.Mutex
	var errArr []error
	var wg sync.WaitGroup

	for index, request := range requests {
		wg.Add(1)
		go func(index int, request []interface{}) {
			defer wg.Done()
			result, err := e.enforce("", nil, request...)
			if err != nil {
				mu1.Lock()
				errArr = append(errArr, err)
				mu1.Unlock()
				//return results, err
			}
			results[index] = result
			//results = append(results, result)
		}(index, request)
	}
	wg.Wait()

	if len(errArr) != 0 {
		return nil, fmt.Errorf("%s \n", errArr)
	}
	return results, nil
}

//
func (e *Enforcer) BatchEnforceBasic(requests [][]interface{}) ([]bool, error) {
	var results []bool
	for _, request := range requests {
		result, err := e.enforce("", nil, request...)
		if err != nil {
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

func testBatchEnforceBenchmarkV2(t *testing.B, e *Enforcer, requests [][]interface{}, results []bool) {
	t.Helper()
	myRes, _ := e.BatchEnforceV2(requests)
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

//func BenchmarkEnforcer_BatchEnforceBasicSyncScale1(b *testing.B) {
//	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
//	input, _ := testSetGen(1)
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		_, _ = e.BatchEnforceBasic(input)
//	}
//}
//
//func BenchmarkEnforcer_BatchEnforceBasicSyncScale2(b *testing.B) {
//	scale := 2
//	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
//	input, _ := testSetGen(scale)
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		_, _ = e.BatchEnforceBasic(input)
//	}
//}
//
//func BenchmarkEnforcer_BatchEnforceBasicSyncScale5(b *testing.B) {
//	scale := 5
//	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
//	input, _ := testSetGen(scale)
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		_, _ = e.BatchEnforceBasic(input)
//	}
//}

func BenchmarkEnforcer_BatchEnforceBasicScale1(b *testing.B) {
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	input, _ := testSetGen(1)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforceBasic(input)
		}
	})
}

func BenchmarkEnforcer_BatchEnforceBasicScale2(b *testing.B) {
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	scale := 2
	input, _ := testSetGen(scale)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforceBasic(input)
		}
	})
}

func BenchmarkEnforcer_BatchEnforceBasicScale5(b *testing.B) {
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	scale := 5
	input, _ := testSetGen(scale)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforceBasic(input)
		}
	})
}

func BenchmarkEnforcer_BatchEnforceBasicScale10(b *testing.B) {
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	scale := 10
	input, _ := testSetGen(scale)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforceBasic(input)
		}
	})
}

func BenchmarkEnforcer_BatchEnforceBasicScale20(b *testing.B) {
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	scale := 20
	input, _ := testSetGen(scale)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforceBasic(input)
		}
	})
}

func BenchmarkEnforcer_BatchEnforceBasicScale50(b *testing.B) {
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	scale := 50
	input, _ := testSetGen(scale)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforceBasic(input)
		}
	})
}

func BenchmarkEnforcer_BatchEnforceBasicScale256(b *testing.B) {
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	scale := 256
	input, _ := testSetGen(scale)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforceBasic(input)
		}
	})
}

//func BenchmarkEnforcer_BatchEnforceSyncScale1(b *testing.B) {
//	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
//	input, _ := testSetGen(1)
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		_, _ = e.BatchEnforce(input)
//	}
//}
//
//func BenchmarkEnforcer_BatchEnforceSyncScale2(b *testing.B) {
//	scale := 2
//	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
//	input, _ := testSetGen(scale)
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		_, _ = e.BatchEnforce(input)
//	}
//}
//
//func BenchmarkEnforcer_BatchEnforceSyncScale5(b *testing.B) {
//	scale := 5
//	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
//	input, _ := testSetGen(scale)
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		_, _ = e.BatchEnforce(input)
//	}
//}

func BenchmarkEnforcer_BatchEnforceScale1(b *testing.B) {
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	input, _ := testSetGen(1)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforce(input)
		}
	})
}

func BenchmarkEnforcer_BatchEnforceScale2(b *testing.B) {
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	scale := 2
	input, _ := testSetGen(scale)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforce(input)
		}
	})
}

func BenchmarkEnforcer_BatchEnforceScale5(b *testing.B) {
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	scale := 5
	input, _ := testSetGen(scale)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforce(input)
		}
	})
}

func BenchmarkEnforcer_BatchEnforceScale10(b *testing.B) {
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	scale := 10
	input, _ := testSetGen(scale)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforce(input)
		}
	})
}

func BenchmarkEnforcer_BatchEnforceScale20(b *testing.B) {
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	scale := 20
	input, _ := testSetGen(scale)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforce(input)
		}
	})
}

func BenchmarkEnforcer_BatchEnforceScale50(b *testing.B) {
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	scale := 50
	input, _ := testSetGen(scale)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforce(input)
		}
	})
}

func BenchmarkEnforcer_BatchEnforceScale256(b *testing.B) {
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	scale := 256
	input, _ := testSetGen(scale)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforce(input)
		}
	})
}

//func BenchmarkEnforcer_BatchEnforceV2SyncScale1(b *testing.B) {
//	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
//	input, _ := testSetGen(1)
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		_, _ = e.BatchEnforceV2(input)
//	}
//}
//
//func BenchmarkEnforcer_BatchEnforceV2SyncScale2(b *testing.B) {
//	scale := 2
//	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
//	input, _ := testSetGen(scale)
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		_, _ = e.BatchEnforceV2(input)
//	}
//}
//
//func BenchmarkEnforcer_BatchEnforceV2SyncScale5(b *testing.B) {
//	scale := 5
//	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
//	input, _ := testSetGen(scale)
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		_, _ = e.BatchEnforceV2(input)
//	}
//}

func BenchmarkEnforcer_BatchEnforceV2Scale1(b *testing.B) {
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	input, _ := testSetGen(1)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforceV2(input)
		}
	})
}

func BenchmarkEnforcer_BatchEnforceV2Scale2(b *testing.B) {
	scale := 2
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	input, _ := testSetGen(scale)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforceV2(input)
		}
	})
}

func BenchmarkEnforcer_BatchEnforceV2Scale5(b *testing.B) {
	scale := 5
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	input, _ := testSetGen(scale)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforceV2(input)
		}
	})
}

func BenchmarkEnforcer_BatchEnforceV2Scale10(b *testing.B) {
	scale := 10
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	input, _ := testSetGen(scale)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforceV2(input)
		}
	})
}

func BenchmarkEnforcer_BatchEnforceV2Scale20(b *testing.B) {
	scale := 20
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	input, _ := testSetGen(scale)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforceV2(input)
		}
	})
}

func BenchmarkEnforcer_BatchEnforceV2Scale50(b *testing.B) {
	scale := 50
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	input, _ := testSetGen(scale)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforceV2(input)
		}
	})
}

func BenchmarkEnforcer_BatchEnforceV2Scale256(b *testing.B) {
	scale := 256
	e, _ := NewEnforcer("examples/priority_model_explicit.conf", "examples/priority_policy_explicit.csv")
	input, _ := testSetGen(scale)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = e.BatchEnforceV2(input)
		}
	})

}
