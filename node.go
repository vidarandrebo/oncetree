package oncetree

import (
	"github.com/relab/gorums"
	"github.com/vidarandrebo/oncetree/protos"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"time"
)

type Node struct {
	rpcAddr         string
	keyValueStorage *KeyValueStorage
	gorumsServer    *gorums.Server
	id              string
}

func NewNode(id string, rpcAddr string) *Node {
	return &Node{
		rpcAddr:         rpcAddr,
		keyValueStorage: NewKeyValueStorage(),
		id:              id,
	}
}

// Run starts the main loop of the node
func (n *Node) Run() {
	n.startGorumsServer(n.rpcAddr)

	for {
		select {
		case <-time.After(time.Second * 1):
			log.Println("send heartbeat here...")
		}
	}
}

func (n *Node) startGorumsServer(addr string) {
	n.gorumsServer = gorums.NewServer()
	protos.RegisterKeyValueServiceServer(n.gorumsServer, n)
	listener, listenErr := net.Listen("tcp", addr)
	if listenErr != nil {
		log.Panicf("could not listen to address %v", addr)
	}
	go func() {
		serveErr := n.gorumsServer.Serve(listener)
		if serveErr != nil {
			log.Panicf("gorums server could not serve key value server")
		}
	}()
}

func (n *Node) Write(ctx gorums.ServerCtx, request *protos.WriteRequest) (response *emptypb.Empty, err error) {
	n.keyValueStorage.WriteValue(request.GetID(), request.GetKey(), request.GetValue())
	return &emptypb.Empty{}, nil
}

func (n *Node) Read(ctx gorums.ServerCtx, request *protos.ReadRequest) (response *protos.ReadResponse, err error) {
	value, err := n.keyValueStorage.ReadValue(request.Key)
	if err != nil {
		return &protos.ReadResponse{Value: 0}, err
	}
	return &protos.ReadResponse{Value: value}, nil
}

func (n *Node) ReadAll(ctx gorums.ServerCtx, request *protos.ReadRequest) (response *protos.ReadAllResponse, err error) {
	//TODO implement me
	panic("implement me")
}
