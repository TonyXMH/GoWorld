package dispatcher_client

import (
	"../../../common"
	"../../../proto"
	"errors"
	"github.com/TonyXMH/GoWorld/config"
	"github.com/TonyXMH/GoWorld/gwlog"
	"github.com/TonyXMH/GoWorld/netutil"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	LOOP_DELAY_ON_DISPATCHER_CLIENT_ERROR = time.Second
)

var (
	_dispatcherClient         *DistpatcherClient
	dispatcherClientDelegate  IDispatcherClientDelegate
	errDispatcherNotConnected = errors.New("dispatcher not connected")
)

func getDispatcherClient() *DistpatcherClient {
	addr := (*uintptr)(unsafe.Pointer(&_dispatcherClient))
	return (*DistpatcherClient)(unsafe.Pointer(atomic.LoadUintptr(addr)))
}

func setDispatcherClient(dc *DistpatcherClient) {
	addr := (*uintptr)(unsafe.Pointer(&_dispatcherClient))
	atomic.StoreUintptr(addr, uintptr(unsafe.Pointer(dc)))
}

func assureConnectedDispatcherClient() *DistpatcherClient {
	dispatcherClient := getDispatcherClient()
	for dispatcherClient == nil {
		dispatcherClient, err := connectDispatchClient()
		if err != nil {
			gwlog.Error("Connect to dispatcher failed:%s", err.Error())
			time.Sleep(LOOP_DELAY_ON_DISPATCHER_CLIENT_ERROR)
			continue
		}
		setDispatcherClient(dispatcherClient)
		dispatcherClientDelegate.OnDispatcherClientConnect()
		gwlog.Info("dispatcher_client: connected to dispatcher:%s", dispatcherClient)
	}
	return dispatcherClient
}

func connectDispatchClient() (*DistpatcherClient, error) {
	dispatcherConfig := config.GetDispatcher()
	conn, err := netutil.ConnectTCP(dispatcherConfig.Ip, dispatcherConfig.Port)
	if err != nil {
		return nil, err
	}
	return newDispatcherClient(conn), nil
}

type IDispatcherClientDelegate interface {
	OnDispatcherClientConnect()
	HandleDeclareService(entityID common.EntityID, serviceName string)
}

func Initialize(delegate IDispatcherClientDelegate) {
	dispatcherClientDelegate = delegate
	assureConnectedDispatcherClient()
	go netutil.ServeForever(serveDispatcherClient) //start the recv routine
}

func GetDispatcherClientForSend() *DistpatcherClient {
	return getDispatcherClient()
}

func serveDispatcherClient() {
	gwlog.Debug("serveDispatcherClient: start serving dispatcher client...")
	for {
		dispatcherClient := assureConnectedDispatcherClient()
		var msgtype proto.MsgType_t
		pkt, err := dispatcherClient.Recv(&msgtype)
		if err != nil {
			gwlog.Error("serveDispatcherClient: RecvMsgPacket error: %s", err.Error())
			dispatcherClient.Close()
			setDispatcherClient(nil)
			time.Sleep(LOOP_DELAY_ON_DISPATCHER_CLIENT_ERROR)
			continue
		}
		gwlog.Info("%s.RecvPacket: msgtype=%v,data=%v", dispatcherClient, msgtype, pkt.Payload())
		if msgtype == proto.MT_DECLEARE_SERVICE {
			eid := pkt.ReadEntityID()
			serviceName := pkt.ReadVarStr()
			dispatcherClientDelegate.HandleDeclareService(eid, serviceName)
		} else {
			gwlog.TraceError("unknown msgtype:%v", msgtype)
		}
		pkt.Release()
	}
}
