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
	count_recv:=0
	addr, err := net.ResolveTCPAddr("tcp4", "localhost:8899")
	handleErr(err)
	receiver, err := net.DialTCP("tcp4", nil, addr)
	handleErr(err)
	fmt.Println("Connect to server,start dispatching on port 8900!")
	buffer:=make([]byte,1024)
	receivers:=make([]*net.TCPConn,1024)
	var sendbuffer []byte
	ch1 :=make(chan struct{})
	go func() {
		for {
			cnt, err = receiver.Read(buffer)
			sendbuffer=make([]byte,cnt)
			copy(sendbuffer,buffer)
			for i:=0;i<count_recv;i++{
				receivers[i].Write(sendbuffer)
			}
		}
	}()
	go func() {
		addr2,err:=net.ResolveTCPAddr("tcp4",":8900")
		handleErr(err)
		listenerconn,err:=net.ListenTCP("tcp4",addr2)
		for {
			listener, err := listenerconn.AcceptTCP()
			receivers[count_recv]=listener
			handleErr(err)
			count_recv++
			fmt.Printf("Receiver %v connected,current receiver count = %d\n",listener.RemoteAddr(),count_recv)

		}
	}()
	<-ch1

}

