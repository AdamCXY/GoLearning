package rxgo

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type filteringOperator struct {
	opFunc func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool)
}

func (ftop filteringOperator) op(ctx context.Context, o *Observable) {
	in := o.pred.outflow
	out := o.outflow
	//fmt.Println(o.name, "operator in/out chan ", in, out)
	var wg sync.WaitGroup
	go func() {
		end := false
		for x := range in {
			if end {
				break
			}
			// can not pass a interface as parameter (pointer) to gorountion for it may change its value outside!
			xv := reflect.ValueOf(x)
			// send an error to stream if the flip not accept error
			if e, ok := x.(error); ok && !o.flip_accept_error {
				o.sendToFlow(ctx, e, out)
				continue
			}
			// scheduler
			switch threading := o.threading; threading {
			case ThreadingDefault:
				if ftop.opFunc(ctx, o, xv, out) {
					end = true
				}
			case ThreadingIO:
				fallthrough
			case ThreadingComputing:
				wg.Add(1)
				go func() {
					defer wg.Done()
					if ftop.opFunc(ctx, o, xv, out) {
						end = true
					}
				}()
			default:
			}
		}
		if o.flip != nil {
			buffer := (reflect.ValueOf(o.flip))
			if buffer.Kind() != reflect.Slice {
				panic("flip is not buffer")
			}
			for i := 0; i < buffer.Len(); i++ {
				o.sendToFlow(ctx, buffer.Index(i).Interface(), out)
			}
		}
		wg.Wait() //waiting all go-routines completed
		o.closeFlow(out)
	}()
}

// Debounce : 间隔timespan时间输出item；
func (parent *Observable) Debounce(timespan time.Duration) (o *Observable) {

	o = parent.newTransformObservable("Debounce")

	o.flip_accept_error = false
	o.flip_sup_ctx = false
	o.flip = nil
	o.threading = ThreadingComputing
	count := 0
	o.operator = filteringOperator{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			count++
			var tempCount = count
			//fmt.Printf("Debunce x %d with tempcount %d\n", item.Interface().(int), tempCount)
			time.Sleep(timespan)
			time.Sleep(5 * time.Microsecond)
			//fmt.Printf("Debunce x %d with count %d\n", item.Interface().(int), count)
			if tempCount == count {
				end = o.sendToFlow(ctx, item.Interface(), out)
			}
			return
		}}
	return o
}

func (parent *Observable) Distinct() (o *Observable) {
	o = parent.newTransformObservable("distinct")
	o.flip_accept_error = true
	o.flip_sup_ctx = true
	m := map[string]bool{}
	o.operator = filteringOperator{
		opFunc: func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			itemStr := fmt.Sprintf("%v", item)
			if _, ok := m[itemStr]; !ok {
				m[itemStr] = true
				o.sendToFlow(ctx, item.Interface(), out)
			}
			return false
		},
	}
	return o
}

// ElementAt: 返回位于第id位的元素，从0开始
func (parent *Observable) ElementAt(id int) (o *Observable) {
	o = parent.newTransformObservable("ElementAt")
	o.flip_accept_error = false
	o.flip_sup_ctx = false
	o.flip = nil
	count := 0
	o.operator = filteringOperator{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			if count == id {
				end = o.sendToFlow(ctx, item.Interface(), out)
			}
			count++
			return
		}}
	return o
}

func (parent *Observable) First() (o *Observable) {
	o = parent.newTransformObservable("first")
	o.flip_accept_error = true
	o.flip_sup_ctx = true
	o.operator = filteringOperator{
		opFunc: func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			o.sendToFlow(ctx, item.Interface(), out)
			return true
		},
	}
	return o
}

// IgnoreElements: 忽略所有元素，只发送结束或是错误信息；
func (parent *Observable) IgnoreElements() (o *Observable) {
	o = parent.newTransformObservable("IgnoreElements")
	o.flip_accept_error = false
	o.flip_sup_ctx = false
	o.flip = nil
	o.operator = filteringOperator{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			return
		}}
	return o
}

func (parent *Observable) Last() (o *Observable) {
	o = parent.newTransformObservable("last")
	o.flip_accept_error = true
	o.flip_sup_ctx = true
	o.operator = filteringOperator{
		opFunc: func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			o.flip = append([]interface{}{}, item.Interface())
			return false
		},
	}
	return o
}

// Sample: 定期发射数据
func (parent *Observable) Sample(timespan time.Duration) (o *Observable) {
	o = parent.newTransformObservable("Sample")
	o.flip_accept_error = false
	o.flip_sup_ctx = false
	o.flip = nil
	o.threading = ThreadingComputing
	temp := make([]reflect.Value, 128)
	var hastimer = false
	var count int
	var wg sync.WaitGroup
	count = 0
	o.operator = filteringOperator{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			temp[count] = item
			count = count + 1
			//fmt.Printf("append: %d\n", item.Interface().(int))
			if !hastimer {
				hastimer = true
				wg.Add(1)
				go func() {
					defer wg.Done()
					for {
						time.Sleep(timespan + 50*time.Millisecond)
						select {
						case <-ctx.Done():
							return
						default:
							if count == 0 {
								return
							}
							//fmt.Printf("sample %d with beforeitem %d\n", temp[count-1].Interface(), count-1)
							if o.sendToFlow(ctx, temp[count-1].Interface(), out) {
								return
							}
							count = 0
						}
					}
				}()
				wg.Wait()
			}
			return false
		}}
	return o
}

func (parent *Observable) Skip(num int) (o *Observable) {
	o = parent.newTransformObservable("skip.n")
	o.flip_accept_error = true
	o.flip_sup_ctx = true
	count := 0
	o.operator = filteringOperator{
		opFunc: func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			count++
			if count > num {
				o.sendToFlow(ctx, item.Interface(), out)
			}
			return false
		},
	}

	return o
}

// Skiplast：跳过后n项后再发送
func (parent *Observable) Skiplast(n int) (o *Observable) {
	o = parent.newTransformObservable("Skiplast")
	o.flip_accept_error = false
	o.flip_sup_ctx = false
	o.flip = nil
	var temp []reflect.Value
	var tempcount = 0
	o.operator = filteringOperator{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			if tempcount == n {
				end = o.sendToFlow(ctx, temp[0].Interface(), out)
				temp = temp[1:]
			} else {
				tempcount++
			}
			temp = append(temp, item)
			return
		}}
	return o
}

func (parent *Observable) Take(num int) (o *Observable) {
	o = parent.newTransformObservable("take.n")
	o.flip_accept_error = true
	o.flip_sup_ctx = true
	count := 0
	o.operator = filteringOperator{
		opFunc: func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			count++
			if count > num {
				return true
			}
			o.sendToFlow(ctx, item.Interface(), out)
			return false
		},
	}

	return o
}

func (parent *Observable) TakeLast(num int) (o *Observable) {
	o = parent.newTransformObservable("takeLast.n")
	o.flip_accept_error = true
	o.flip_sup_ctx = true
	count := 0
	var lasts []interface{}
	o.operator = filteringOperator{
		opFunc: func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool) {
			count++
			if count <= num {
				lasts = append(lasts, item.Interface())
			} else {
				lasts = lasts[1:]
				lasts = append(lasts, item.Interface())
			}
			o.flip = lasts
			return false
		},
	}

	return o
}
