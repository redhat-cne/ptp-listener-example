package main

import (
	"os"
	"time"

	ptpEvent "github.com/redhat-cne/sdk-go/pkg/event/ptp"
	"github.com/sirupsen/logrus"
	exports "github.com/test-network-function/ptp-listener-exports"
	lib "github.com/test-network-function/ptp-listener-lib"
)

func main() {
	lib.InitPubSub()
	listener := func(name string, ch <-chan exports.StoredEvent) {
		for i := range ch {
			logrus.Infof("[%s] got %v\n", name, i)
		}
		logrus.Infof("[%s] done\n", name)
	}

	ch1 := lib.Ps.Subscribe(string(ptpEvent.OsClockSyncStateChange))

	err := lib.StartListening(
		exports.Port9085,
		exports.Port9085,
		"linuxptp-daemon-nxsds",
		"master2",
		"openshift-ptp",
		"/home/deliedit/.kube/config.bos2.cluster-2",
		"https://api.cluster-2.cnfcertlab.org:6443",
	)

	if err != nil {
		logrus.Errorf("could not start listening for events, err=%s", err)
		os.Exit(1)
	}

	go listener("1", ch1)
	const sleepTimeout = 30
	time.Sleep(time.Minute * sleepTimeout)
	lib.Ps.Close()
}
