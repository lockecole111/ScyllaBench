package main

import (
    "github.com/gocql/gocql"
    "time"
    "log"
    "strings"
    "flag"
    "sync"
)

type Config struct {    
    ratio int
    ttl int 
    write_thread int
    read_thread int
    nodes []string
    id gocql.UUID
    interval int
}


func Conn(nodes ...string )  *gocql.Session  {

    cluster := gocql.NewCluster(nodes...)
    cluster.Keyspace = "test"
    cluster.Consistency = gocql.Quorum
    cluster.NumConns = 1
    cluster.ConnectTimeout = time.Second * 60
    session, _ := cluster.CreateSession()
    if session != nil {
        return session
    }
    return Conn()
}

func TestWrite(write chan bool,config Config){    
    session := Conn(config.nodes...)
    
    for {
        if <- write {
            if err := session.Query(`INSERT INTO example (id, title, content) VALUES (?, ?, ?) using TTL ?;`,
            config.id, RandString(16),RandString(32),config.ttl ).Exec(); err != nil {
            }
        }
    }            
}


func TestRead(read chan bool, config Config){   
    session := Conn(config.nodes...)
    for {
        if <- read {
            if err := session.Query(`select * from example where id=? limit 1`,
            config.id ).Exec(); err != nil {
            }
        }
    }   
}


func CountQPS(total_write, total_read *int, config Config){
    
    for {
        time.Sleep(time.Second * time.Duration(config.interval))
        log.Print("QPS Write:",*total_write/config.interval)
        log.Print("QPS Read:", *total_read/config.interval)
        *total_write = 0
        *total_read = 0
        }    
}



func ParseParam(config *Config) {
    var nodes_flag string
    flag.StringVar(&nodes_flag, "nodes", "localhost", "nodes' ip or domain.")
    flag.IntVar(&config.ttl, "ttl", 600, "TTL.")
    flag.IntVar(&config.interval, "interval", 5, "sampling interval second(s) of countqps.")
    flag.IntVar(&config.ratio, "ratio", 1, "ratio of write/read,and if set to 0 means read and write are independent without ratio limit.")
    flag.IntVar(&config.write_thread, "write_thread", 10, "number of write thread(s).")
    flag.IntVar(&config.read_thread, "read_thread", 10, "number of read thread(s).")
    flag.Parse()
    config.nodes = strings.Split(nodes_flag,",")
}    
func Read(read chan bool,total_read *int) {
    for {
        *total_read = *total_read + 1
        read <- true
        }
    }
func Write(write chan bool,total_write *int) {
    for {
        *total_write = *total_write + 1
        write <- true
        }
    }

func main() {
    
    write := make(chan bool)
    read := make(chan bool)
    var wait_group sync.WaitGroup
    var total_write int = 0
    var total_read int = 0

    var config Config
    ParseParam(&config)
    config.id = gocql.TimeUUID()

    for i := 0; i < config.write_thread; i++ {
        go TestWrite(write, config)
    }
    for i := 0; i < config.read_thread; i++ {
        go TestRead(read, config)
    }
    go CountQPS(&total_write, &total_read, config)

    if config.ratio==0 {
        go Read(read, &total_read)
        go Write(write, &total_write)
        wait_group.Wait()
        }
    for {
        
        total_write = total_write + int(config.ratio)
        for i := 0;i < config.ratio ; i++ {
            write <- true
        }
        total_read = total_read + 1
        read <- true
    }

}
