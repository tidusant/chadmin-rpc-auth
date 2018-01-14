package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"strconv"
	"strings"

	"github.com/tidusant/c3m-common/c3mcommon"
	"github.com/tidusant/c3m-common/log"
	rpch "github.com/tidusant/chadmin-repo/cuahang"
	"github.com/tidusant/chadmin-repo/models"
)

const (
	defaultcampaigncode string = "XVsdAZGVmY"
)

type Arith int

func (t *Arith) Run(data string, result *string) error {
	log.Debugf("Call RPCAuth args:" + data)
	*result = ""
	//parse args
	args := strings.Split(data, "|")

	var usex models.UserSession
	usex.Session = args[0]
	usex.Action = args[2]
	userIP := args[1]
	usex.Params = ""

	if len(args) > 3 {
		usex.Params = args[3]
	}

	if usex.Action == "l" {
		*result = login(usex, userIP)
	} else if usex.Action == "test" {
		*result = test(usex, userIP)
	} else if usex.Action == "aut" {
		*result = rpch.GetLogin(usex.Session, userIP)
	} else if usex.Action == "submitorder" {
		//*result = submitorder(siteid, mongoSession, data2)
	} else { //default
		*result = ""
	}

	return nil
}

func test(usex models.UserSession, userIP string) string {
	if rpch.GetLogin(usex.Session, userIP) != "" {
		return c3mcommon.ReturnJsonMessage("1", "", "user logged in", "")
	}
	return c3mcommon.ReturnJsonMessage("0", "user not logged in", "", "")
}

func login(usex models.UserSession, userIP string) string {
	args := strings.Split(usex.Params, ",")
	if len(args) < 2 {
		return c3mcommon.ReturnJsonMessage("0", "empty username or pass", "", "")
	}
	user := args[0]
	pass := args[1]
	userid := rpch.Login(user, pass, usex.Session, userIP)
	if userid != "" {
		return c3mcommon.ReturnJsonMessage("1", "", "login success", "")
	}
	return c3mcommon.ReturnJsonMessage("0", "login fail", "", "")

}
func main() {
	var port int
	var debug bool
	flag.IntVar(&port, "port", 9877, "help message for flagname")
	flag.BoolVar(&debug, "debug", false, "Indicates  if debug messages should be printed in log files")
	flag.Parse()

	logLevel := log.DebugLevel
	if !debug {
		logLevel = log.InfoLevel

	}

	log.SetOutputFile(fmt.Sprintf("adminAuth-"+strconv.Itoa(port)), logLevel)
	defer log.CloseOutputFile()
	log.RedirectStdOut()

	//init db
	//

	arith := new(Arith)
	rpc.Register(arith)
	log.Infof("running with port:" + strconv.Itoa(port))

	//			rpc.HandleHTTP()
	//			l, e := net.Listen("tcp", ":"+strconv.Itoa(port))
	//			if e != nil {
	//				log.Debug("listen error:", e)
	//			}
	//			http.Serve(l, nil)

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port))
	c3mcommon.CheckError("rpc dail:", err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	c3mcommon.CheckError("rpc init listen", err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(conn)
	}

}
