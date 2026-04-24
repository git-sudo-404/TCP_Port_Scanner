package main;

import (
  "fmt"
  "net"
  "sync"
  "time"
)

func scanPort(host string,port int,wg *sync.WaitGroup){
    
  defer wg.Done()

  address:=fmt.Sprintf("%s:%d",host,port)
  conn,err := net.DialTimeout("tcp",address,1 * time.Second)

  if err!=nil{
    return
  }

  conn.Close()
  fmt.Printf("[OPEN] Port %d\n",port)

}

func main(){

  var wg sync.WaitGroup
  host:="localhost"

  for port:=1;port<=65535;port++{
    wg.Add(1)
    go scanPort(host,port,&wg)
  }

  wg.Wait()


}
