package server

import (
	codec "github.com/chobie/momonga/encoding/mqtt"
	"github.com/chobie/momonga/util"
	"github.com/chobie/momonga/configuration"
)

const KILOBYTE = 1024
const MEGABYTE = 1024 * KILOBYTE
const MAX_REQUEST_SIZE = MEGABYTE * 2

func NewTcpServer(conf *configuration.Config) *TcpServer{
	server := &TcpServer{
		forceSSLUsers: map[string]bool{},
		Connections: map[string]Connection{},
		Engine: &Pidgey{
			Topics: map[string]*Topic{},
			Queue: make(chan codec.Message, 8192),
			OutGoingTable: util.NewMessageTable(),
			Qlobber: util.NewQlobber(),
			Retain: map[string]*codec.PublishMessage{},
			Connections: map[string]*MmuxConnection{},
			SubscribeMap: map[string]string{},
			RetryMap: map[string][]*Retryable{},
			ErrorChannel: make(chan *Retryable, 8192),
		},
	}
	server.Engine.SetupCallback()
	server.listenAddress = conf.GetListenAddress()
	server.SSLlistenAddress = conf.GetSSLListenAddress()

// TODO: Unix Socket
	//server.listenSocket = "/tmp/hoge.sock"//config.TcpInputSocketString()

// TODO: SSL
//	if config.TcpInputUseSSL {
//		cert, err := tls.LoadX509KeyPair(config.TcpInputSSLCert(), config.TcpInputSSLKey())
//		if err != nil {
//			log.Error("tcp server: loadkeys failed. disable ssl feature: %s", err)
//		} else {
//			tslConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
//			tslConfig.Rand = rand.Reader
//
//			server.tlsConfig = tslConfig
//			for _, name := range config.TcpInputForceSSL() {
//				server.forceSSLUsers[name] = true
//			}
//			log.Debug("SSL Config loaded")
//		}
//	}

	server.shutdown = make(chan bool, 1)
	return server
}
