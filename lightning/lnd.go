package lightning

import (
	"github.com/lightningnetwork/lnd/lnrpc"
	"google.golang.org/grpc"
)

type LndClient struct {
	client lnrpc.LightningClient
	conn   grpc.ClientConn
	opt    grpc.CallOption
}
