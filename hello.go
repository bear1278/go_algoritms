package main

import (
	"fmt"
	//"math/rand"
	"sync"
	"time"
)

const N = 5

func main() {

   fn := func(x int) int {
      //time.Sleep(time.Duration(rand.Int31n(N)) * time.Second)
      return x *2
   }
   in1 := make(chan int, N)
   in2 := make(chan int, N)
   out := make(chan int, N)

   start := time.Now()
   merge2Channels(fn, in1, in2, out, N+1)
   for i := 0; i < N+1; i++ {
      in1 <- i
      in2 <- i
   }

   orderFail := false
   EvenFail := false
   for i, prev := 0, 0; i < N; i++ {
      c := <-out
      if c%2 != 0 {
         EvenFail = true
      }
      if prev >= c && i != 0 {
         orderFail = true
      }
      prev = c
      fmt.Println(c)
   }
   if orderFail {
      fmt.Println("порядок нарушен")
   }
   if EvenFail {
      fmt.Println("Есть не четные")
   }
   duration := time.Since(start)
   if duration.Seconds() > N {
      fmt.Println("Время превышено")
   }
   fmt.Println("Время выполнения: ", duration)
}

// Merge2Channels below

func merge2Channels(fn func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
    waiter:=new(sync.WaitGroup)
    mu := new(sync.Mutex)
    res1:=make([]int,n)
    res2:=make([]int,n)
    val1:=0
    val2:=0
    waiter.Add(2*n)
    for i:=0;i<n;i++{
        go func(wg *sync.WaitGroup, m *sync.Mutex,ind int){
            defer wg.Done()
            m.Lock()
            val1=<-in1
            
            res1[ind]=fn(val1)
            m.Unlock()
        }(waiter,mu,i)
        go func(wg *sync.WaitGroup,m *sync.Mutex, ind int){
            defer wg.Done()
            m.Lock()
            val2=<-in2
            
            res2[ind]=fn(val2)
            m.Unlock()
        }(waiter,mu,i)
    }
   
   
    go func(){
        waiter.Wait()
        
        for i:=0;i<n;i++{
            
            out<-(res1[i]+res2[i])
            
        }
       
       
    }()
 }
