package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
	"unsafe"
)

func main() {
	conn, err := net.Dial("tcp", "192.168.2.20:20080")
	//conn, err := net.Dial("tcp", "172.21.234.16:20080")
	fmt.Println("Connection established")
	defer conn.Close()
	if err != nil {
		log.Fatalln("Unable to connect to the server")
	}

	b := make([]byte, 6144)
	for i := 0; i < len(b)/4; i++ {
		for j := 0; j < 4; j++ {
			b[i*4+3-j] = *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&i)) + uintptr(j)))
		}
	}

	tstamp := make(chan time.Time)
	done := make(chan bool)

	go func() {
		run := 0
		num_err := 0
		f, _ := os.Create("wifi.log")
		defer f.Close()

		rb := make([]byte, 6144)
		for {
			time_sent, result := <-tstamp
			if result == false {
				fmt.Println("Done")
				return
			}
			_, err := conn.Read(rb[0:])
			time_received := time.Now()
			fmt.Println(time_received.Sub(time_sent))
			if err != nil {
				fmt.Println("Error reading data")
				continue
			}

			// Output received data to file in a readable format
			/*
				fname := fmt.Sprintf("run%d.txt", run)

				f, err := os.Create(fname)
				if err != nil {
					log.Fatal("Failed to create file\n")
				}
				defer f.Close()
				val := int32(0)
				for i := 0; i < len(b)/4; i++ {
					for j := 0; j < 4; j++ {
						*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(j))) = rb[i*4+3-j]
					}
					str := fmt.Sprintf("Expected %d, received %d\n", i, val)
					f.WriteString(str)
				}
			*/
			for i := 0; i < len(rb); i++ {
				if b[i] != rb[i] {
					//fmt.Printf("Data mismatch: expected 0x%x, received 0x%x\n", b[i], rb[i])
					num_err++
				}
			}

			str := fmt.Sprintf("Run %d completed with %d errors\n", run, num_err)
			f.WriteString(str)

			num_err = 0
			run++
			done <- true
		}
	}()

	for i := 0; i < 20; i++ {
		tstamp <- time.Now()
		conn.Write([]byte(b))
		_ = <-done
	}

	close(tstamp)
	close(done)
}
