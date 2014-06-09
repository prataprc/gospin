// contains common test fixtures.

package failsafe

import(
    "github.com/goraft/raft"
    "path/filepath"
    "io/ioutil"
    "log"
    "net/http"
    "net"
    "os"
)

var testdir = "testdata"
var testRaftdir = filepath.Join(testdir, "server")
var host = "localhost"
var port = "4000"
var servAddr = "http://" +host + ":" + port
var lisAddr  = host + ":" + port
var smallJsonFile = filepath.Join(testdir, "small.json")
var smallJson, _ = ioutil.ReadFile(smallJsonFile)
var dummyFile = filepath.Join(testdir, "_dummy.json")

// servdir -> {*Server, net.Listener, *http.Server}
var activeServers map[string][]interface{}

// register commands for testing and start a server.
func init() {
    raft.RegisterCommand(&SetCommand{})
    raft.RegisterCommand(&DeleteCommand{})
    activeServers = make(map[string][]interface{})
}

func startTestServer(servdir string) (*Server, net.Listener, *http.Server) {
    if v, ok := activeServers[servdir]; ok {
        return v[0].(*Server), v[1].(net.Listener), v[2].(*http.Server)
    }

    os.RemoveAll(servdir)
    log.SetOutput(ioutil.Discard)
    mux := http.NewServeMux()
    srv, err := NewServer(servdir, host, port, mux)
    if err != nil {
        log.Fatal(err)
    }
    srv.Install("")

    daemon := &http.Server{Addr: lisAddr, Handler: mux}

    lis, err := net.Listen("tcp", lisAddr)
    if err != nil {
        log.Fatal(err)
    }
    go daemon.Serve(lis)
    activeServers[servdir] = []interface{}{srv, lis, daemon}
    return srv, lis, daemon
}