package main

import (
	"net"
	"log"
	"fmt"
)

func handleErr(err error){
	if err!=nil {
		log.Fatal(err)
	}
}

func main() {
	var cnt int
	count:=0
	count_recv:=0
	addr, err := net.ResolveTCPAddr("tcp4", "localhost:8899")
	handleErr(err)
	receiver, err := net.DialTCP("tcp4", nil, addr)
	handleErr(err)
	fmt.Println("Connect to server,start dispatching on port 8900!")
	buffer:=make([]byte,1024)
	ch:=make(chan struct{})
	ch2:=make(chan struct{})
	ch3:=make(chan struct{})
	go func() {
		for {
			<-ch
			count=0
			cnt, err = receiver.Read(buffer)
			for i:=0;i<count_recv;i++{
				ch3<- struct{}{}
			}
			handleErr(err)
		}
	}()
	go func() {
		addr2,err:=net.ResolveTCPAddr("tcp4","localhost:8900")
		handleErr(err)
		listenerconn,err:=net.ListenTCP("tcp4",addr2)
		for {
			listener, err := listenerconn.AcceptTCP()
			count_recv++
			fmt.Printf("Receiver %v connected,current receiver count = %d\n",listener.RemoteAddr(),count_recv)
			handleErr(err)
			go func() {
				for {
					<-ch3
					b := make([]byte, cnt)
					copy(b, buffer)
					listener.Write(b)
					count++
					if count >= count_recv {
						ch <- struct{}{}
					}
				}
			}()
		}
	}()
	ch<- struct{}{}
	<-ch2
}

