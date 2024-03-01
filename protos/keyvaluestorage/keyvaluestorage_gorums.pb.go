// Code generated by protoc-gen-gorums. DO NOT EDIT.
// versions:
// 	protoc-gen-gorums v0.7.0-devel
// 	protoc            v4.25.2
// source: protos/keyvaluestorage/keyvaluestorage.proto

package keyvaluestorage

import (
	context "context"
	fmt "fmt"

	gorums "github.com/relab/gorums"
	encoding "google.golang.org/grpc/encoding"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = gorums.EnforceVersion(7 - gorums.MinVersion)
	// Verify that the gorums runtime is sufficiently up-to-date.
	_ = gorums.EnforceVersion(gorums.MaxVersion - 7)
)

// A Configuration represents a static set of nodes on which quorum remote
// procedure calls may be invoked.
type Configuration struct {
	gorums.RawConfiguration
	nodes []*Node
	qspec QuorumSpec
}

// ConfigurationFromRaw returns a new Configuration from the given raw configuration and QuorumSpec.
//
// This function may for example be used to "clone" a configuration but install a different QuorumSpec:
//
//	cfg1, err := mgr.NewConfiguration(qspec1, opts...)
//	cfg2 := ConfigurationFromRaw(cfg1.RawConfig, qspec2)
func ConfigurationFromRaw(rawCfg gorums.RawConfiguration, qspec QuorumSpec) *Configuration {
	// return an error if the QuorumSpec interface is not empty and no implementation was provided.
	var test interface{} = struct{}{}
	if _, empty := test.(QuorumSpec); !empty && qspec == nil {
		panic("QuorumSpec may not be nil")
	}
	return &Configuration{
		RawConfiguration: rawCfg,
		qspec:            qspec,
	}
}

// Nodes returns a slice of each available node. IDs are returned in the same
// order as they were provided in the creation of the Manager.
//
// NOTE: mutating the returned slice is not supported.
func (c *Configuration) Nodes() []*Node {
	if c.nodes == nil {
		c.nodes = make([]*Node, 0, c.Size())
		for _, n := range c.RawConfiguration {
			c.nodes = append(c.nodes, &Node{n})
		}
	}
	return c.nodes
}

// And returns a NodeListOption that can be used to create a new configuration combining c and d.
func (c Configuration) And(d *Configuration) gorums.NodeListOption {
	return c.RawConfiguration.And(d.RawConfiguration)
}

// Except returns a NodeListOption that can be used to create a new configuration
// from c without the nodes in rm.
func (c Configuration) Except(rm *Configuration) gorums.NodeListOption {
	return c.RawConfiguration.Except(rm.RawConfiguration)
}

func init() {
	if encoding.GetCodec(gorums.ContentSubtype) == nil {
		encoding.RegisterCodec(gorums.NewCodec())
	}
}

// Manager maintains a connection pool of nodes on
// which quorum calls can be performed.
type Manager struct {
	*gorums.RawManager
}

// NewManager returns a new Manager for managing connection to nodes added
// to the manager. This function accepts manager options used to configure
// various aspects of the manager.
func NewManager(opts ...gorums.ManagerOption) (mgr *Manager) {
	mgr = &Manager{}
	mgr.RawManager = gorums.NewRawManager(opts...)
	return mgr
}

// NewConfiguration returns a configuration based on the provided list of nodes (required)
// and an optional quorum specification. The QuorumSpec is necessary for call types that
// must process replies. For configurations only used for unicast or multicast call types,
// a QuorumSpec is not needed. The QuorumSpec interface is also a ConfigOption.
// Nodes can be supplied using WithNodeMap or WithNodeList, or WithNodeIDs.
// A new configuration can also be created from an existing configuration,
// using the And, WithNewNodes, Except, and WithoutNodes methods.
func (m *Manager) NewConfiguration(opts ...gorums.ConfigOption) (c *Configuration, err error) {
	if len(opts) < 1 || len(opts) > 2 {
		return nil, fmt.Errorf("wrong number of options: %d", len(opts))
	}
	c = &Configuration{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case gorums.NodeListOption:
			c.RawConfiguration, err = gorums.NewRawConfiguration(m.RawManager, v)
			if err != nil {
				return nil, err
			}
		case QuorumSpec:
			// Must be last since v may match QuorumSpec if it is interface{}
			c.qspec = v
		default:
			return nil, fmt.Errorf("unknown option type: %v", v)
		}
	}
	// return an error if the QuorumSpec interface is not empty and no implementation was provided.
	var test interface{} = struct{}{}
	if _, empty := test.(QuorumSpec); !empty && c.qspec == nil {
		return nil, fmt.Errorf("missing required QuorumSpec")
	}
	return c, nil
}

// Nodes returns a slice of available nodes on this manager.
// IDs are returned in the order they were added at creation of the manager.
func (m *Manager) Nodes() []*Node {
	gorumsNodes := m.RawManager.Nodes()
	nodes := make([]*Node, 0, len(gorumsNodes))
	for _, n := range gorumsNodes {
		nodes = append(nodes, &Node{n})
	}
	return nodes
}

// Node encapsulates the state of a node on which a remote procedure call
// can be performed.
type Node struct {
	*gorums.RawNode
}

// QuorumSpec is the interface of quorum functions for KeyValueStorage.
type QuorumSpec interface {
	gorums.ConfigOption

	// ReadAllQF is the quorum function for the ReadAll
	// quorum call method. The in parameter is the request object
	// supplied to the ReadAll method at call time, and may or may not
	// be used by the quorum function. If the in parameter is not needed
	// you should implement your quorum function with '_ *ReadRequest'.
	ReadAllQF(in *ReadRequest, replies map[uint32]*ReadAllResponse) (*ReadAllResponse, bool)
}

// ReadAll is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (c *Configuration) ReadAll(ctx context.Context, in *ReadRequest) (resp *ReadAllResponse, err error) {
	cd := gorums.QuorumCallData{
		Message: in,
		Method:  "keyvaluestorage.KeyValueStorage.ReadAll",
	}
	cd.QuorumFunction = func(req protoreflect.ProtoMessage, replies map[uint32]protoreflect.ProtoMessage) (protoreflect.ProtoMessage, bool) {
		r := make(map[uint32]*ReadAllResponse, len(replies))
		for k, v := range replies {
			r[k] = v.(*ReadAllResponse)
		}
		return c.qspec.ReadAllQF(req.(*ReadRequest), r)
	}

	res, err := c.RawConfiguration.QuorumCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*ReadAllResponse), err
}

// Write is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (n *Node) Write(ctx context.Context, in *WriteRequest) (resp *emptypb.Empty, err error) {
	cd := gorums.CallData{
		Message: in,
		Method:  "keyvaluestorage.KeyValueStorage.Write",
	}

	res, err := n.RawNode.RPCCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*emptypb.Empty), err
}

// Read is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (n *Node) Read(ctx context.Context, in *ReadRequest) (resp *ReadResponse, err error) {
	cd := gorums.CallData{
		Message: in,
		Method:  "keyvaluestorage.KeyValueStorage.Read",
	}

	res, err := n.RawNode.RPCCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*ReadResponse), err
}

// Gossip is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (n *Node) Gossip(ctx context.Context, in *GossipMessage) (resp *emptypb.Empty, err error) {
	cd := gorums.CallData{
		Message: in,
		Method:  "keyvaluestorage.KeyValueStorage.Gossip",
	}

	res, err := n.RawNode.RPCCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*emptypb.Empty), err
}

// PrintState is a quorum call invoked on all nodes in configuration c,
// with the same argument in, and returns a combined result.
func (n *Node) PrintState(ctx context.Context, in *emptypb.Empty) (resp *emptypb.Empty, err error) {
	cd := gorums.CallData{
		Message: in,
		Method:  "keyvaluestorage.KeyValueStorage.PrintState",
	}

	res, err := n.RawNode.RPCCall(ctx, cd)
	if err != nil {
		return nil, err
	}
	return res.(*emptypb.Empty), err
}

// KeyValueStorage is the server-side API for the KeyValueStorage Service
type KeyValueStorage interface {
	Write(ctx gorums.ServerCtx, request *WriteRequest) (response *emptypb.Empty, err error)
	Read(ctx gorums.ServerCtx, request *ReadRequest) (response *ReadResponse, err error)
	ReadAll(ctx gorums.ServerCtx, request *ReadRequest) (response *ReadAllResponse, err error)
	Gossip(ctx gorums.ServerCtx, request *GossipMessage) (response *emptypb.Empty, err error)
	PrintState(ctx gorums.ServerCtx, request *emptypb.Empty) (response *emptypb.Empty, err error)
}

func RegisterKeyValueStorageServer(srv *gorums.Server, impl KeyValueStorage) {
	srv.RegisterHandler("keyvaluestorage.KeyValueStorage.Write", func(ctx gorums.ServerCtx, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*WriteRequest)
		defer ctx.Release()
		resp, err := impl.Write(ctx, req)
		gorums.SendMessage(ctx, finished, gorums.WrapMessage(in.Metadata, resp, err))
	})
	srv.RegisterHandler("keyvaluestorage.KeyValueStorage.Read", func(ctx gorums.ServerCtx, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*ReadRequest)
		defer ctx.Release()
		resp, err := impl.Read(ctx, req)
		gorums.SendMessage(ctx, finished, gorums.WrapMessage(in.Metadata, resp, err))
	})
	srv.RegisterHandler("keyvaluestorage.KeyValueStorage.ReadAll", func(ctx gorums.ServerCtx, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*ReadRequest)
		defer ctx.Release()
		resp, err := impl.ReadAll(ctx, req)
		gorums.SendMessage(ctx, finished, gorums.WrapMessage(in.Metadata, resp, err))
	})
	srv.RegisterHandler("keyvaluestorage.KeyValueStorage.Gossip", func(ctx gorums.ServerCtx, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*GossipMessage)
		defer ctx.Release()
		resp, err := impl.Gossip(ctx, req)
		gorums.SendMessage(ctx, finished, gorums.WrapMessage(in.Metadata, resp, err))
	})
	srv.RegisterHandler("keyvaluestorage.KeyValueStorage.PrintState", func(ctx gorums.ServerCtx, in *gorums.Message, finished chan<- *gorums.Message) {
		req := in.Message.(*emptypb.Empty)
		defer ctx.Release()
		resp, err := impl.PrintState(ctx, req)
		gorums.SendMessage(ctx, finished, gorums.WrapMessage(in.Metadata, resp, err))
	})
}

type internalReadAllResponse struct {
	nid   uint32
	reply *ReadAllResponse
	err   error
}