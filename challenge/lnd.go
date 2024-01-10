package challenge

import (
	"context"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lntypes"
	"google.golang.org/grpc"
)

const (
	PaymentErr = "Failed to pay invoice!"
)

// LndClient represents a client for interacting with an LND (Lightning Network Daemon) node.
type LndClient struct {
	client lnrpc.LightningClient // The Lightning gRPC client for making API calls.
	conn   *grpc.ClientConn      // The gRPC client connection to the LND node.
}

func (lnd *LndClient) Client() *lnrpc.LightningClient {
	return &lnd.client
}

func (lnd *LndClient) Conn() *grpc.ClientConn {
	return lnd.conn
}

func NewLndClient(conn *grpc.ClientConn) LightningNode {
	client := lnrpc.NewLightningClient(conn)
	return &LndClient{
		client,
		conn,
	}
}

func (lnd *LndClient) SendPayment(cx context.Context, invoice PaymentRequest) (lntypes.Preimage, error) {
	sendRequest := lnrpc.SendRequest{PaymentRequest: invoice.GetPaymentRequest()}

	response, err := lnd.client.SendPaymentSync(cx, &sendRequest)

	if err != nil {
		return lntypes.Preimage{}, err
	}

	pre_image, err := lntypes.MakePreimage(response.PaymentPreimage)

	if err != nil {
		return lntypes.Preimage{}, err
	}

	return pre_image, nil
}

func (lnd *LndClient) CreateInvoice(cx context.Context, invoice InvoiceBuilder) (PaymentRequest, error) {
	paymentRequest, err := lnd.client.AddInvoice(cx, &invoice)
	return *paymentRequest, err
}
