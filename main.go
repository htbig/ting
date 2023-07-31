package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"
)

func task(exitChanal chan int) {
	var startTick = time.Now().UnixNano() / 1e6
	fmt.Println("start task time:", startTick)
	for {
		if (time.Now().UnixNano()/1e6)-startTick < 200 {
			continue
		}
		time.Sleep(800 * time.Millisecond)
		startTick = time.Now().UnixNano() / 1e6
		select {
		case val, ok := <-exitChanal:
			if !ok {
				log.Printf("Channel closed")
				return
			}
			log.Printf("Revice dataChan %d\n", val)
			return
		default:
			break
		}
	}
}

func getCurrentMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

var memSlice []byte

func doWork(exitChanal chan int) {
	file, err := os.OpenFile("test.txt", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	const memSize = 100 * 1024 * 1024
	memSlice = make([]byte, memSize)

	for i := 0; i < memSize; i++ {
		memSlice[i] = 1
	}

	fmt.Printf("use the memory: %v\n", getCurrentMemoryUsage()/(1024*1024))

	count, err := file.WriteString("Hello Golang")
	if err != nil {
	}
	fmt.Println("hello world", count, err)
	go task(exitChanal)
}
func main() {
	go func() {
		http.ListenAndServe("localhost:8080", nil)
	}()
	runtime.GOMAXPROCS(3)
	fmt.Println("time hour now:", time.Now().Hour())
	var exitChanal = make(chan int)
	for {
		if time.Now().Hour() >= 0 && time.Now().Hour() < 8 {
			doWork(exitChanal)
			time.Sleep(time.Duration((8*60 - (time.Now().Hour()*60 + time.Now().Minute()))) * 60000000000)
			exitChanal <- 1
			memSlice = nil
		} else if time.Now().Hour() >= 8 && time.Now().Hour() < 16 {
			doWork(exitChanal)
			doWork(exitChanal)
			time.Sleep(time.Duration((16*60 - (time.Now().Hour()*60 + time.Now().Minute()))) * 60000000000)
			exitChanal <- 1
			exitChanal <- 2
			memSlice = nil
		} else {
			doWork(exitChanal)
			doWork(exitChanal)
			doWork(exitChanal)
			fmt.Println("after do work:", time.Now().Minute(), 24*60-(time.Now().Hour()*60+time.Now().Minute()), 60000000000*time.Duration((24*60-(time.Now().Hour()*60+time.Now().Minute()))))
			time.Sleep(time.Duration((24*60 - (time.Now().Hour()*60 + time.Now().Minute()))) * 60000000000)
			exitChanal <- 1
			exitChanal <- 2
			exitChanal <- 3
			memSlice = nil
		}

	}

	select {}
}
