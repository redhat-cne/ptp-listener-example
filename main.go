package main

import (
	"fmt"
	"os"
	"time"

	ptpEvent "github.com/redhat-cne/sdk-go/pkg/event/ptp"
	"github.com/sirupsen/logrus"
	exports "github.com/test-network-function/ptp-listener-exports"
	lib "github.com/test-network-function/ptp-listener-lib"
)

const (
	podName             = `linuxptp-daemon-5kzpd`
	nodeName            = `master1`
	kubeconfig          = `/home/deliedit/.kube/config.bos2.cluster-2`
	k8sAPI              = `https://api.cluster-2.cnfcertlab.org:6443`
	namespace           = `openshift-ptp`
	eventLocalhostPort  = 9085
	eventAPIRemotePort  = exports.Port9085
	localHTTPServerPort = 8989
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
		eventLocalhostPort,
		eventAPIRemotePort,
		localHTTPServerPort,
		podName,
		nodeName,
		namespace,
		kubeconfig,
		k8sAPI,
	)

	if err != nil {
		logrus.Errorf("could not start listening for events, err=%s", err)
		os.Exit(1)
	}

	go listener("1", ch1)
	const sleepTimeout = 1
	time.Sleep(time.Minute * sleepTimeout)

	lib.UnsubscribeAllEvents(fmt.Sprintf("localhost:%d", eventLocalhostPort), nodeName)
	lib.Ps.Close()
}
