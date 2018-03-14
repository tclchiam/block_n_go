package main

import (
	"github.com/sirupsen/logrus"

	"github.com/tclchiam/oxidize-go/account"
	"github.com/tclchiam/oxidize-go/cmd/interrupt"
	"github.com/tclchiam/oxidize-go/node"
	"github.com/tclchiam/oxidize-go/oxylogger"
	"github.com/tclchiam/oxidize-go/p2p"
	"github.com/tclchiam/oxidize-go/server/httpserver"
	"github.com/tclchiam/oxidize-go/server/rpc"
	"github.com/tclchiam/oxidize-go/storage"
)

func init() {
	accountLogger := oxylogger.NewLogger(logrus.InfoLevel)
	httpLogger := oxylogger.NewLogger(logrus.InfoLevel)
	interruptLogger := oxylogger.NewLogger(logrus.InfoLevel)
	nodeLogger := oxylogger.NewLogger(logrus.WarnLevel)
	p2pLogger := oxylogger.NewLogger(logrus.WarnLevel)
	rpcLogger := oxylogger.NewLogger(logrus.WarnLevel)
	storageLogger := oxylogger.NewLogger(logrus.WarnLevel)

	account.UseLogger(accountLogger)
	httpserver.UseLogger(httpLogger)
	interrupt.UseLogger(interruptLogger)
	node.UseLogger(nodeLogger)
	p2p.UseLogger(p2pLogger)
	rpc.UseLogger(rpcLogger)
	storage.UseLogger(storageLogger)
}
