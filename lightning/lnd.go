package lightning

import (
	"context"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lntypes"
	"google.golang.org/grpc"
)

const (
	PaymentErr = "Failed to pay invoice!"
)

type LndClient struct {
	client lnrpc.LightningClient
	conn   *grpc.ClientConn
}

func (lnd *LndClient) Client() *lnrpc.LightningClient {
	return &lnd.client
}

func (lnd *LndClient) Conn() *grpc.ClientConn {
	return lnd.conn
}

func NewLndClient(conn *grpc.ClientConn) LndClient {
	client := lnrpc.NewLightningClient(conn)
	return LndClient{
		client,
		conn,
	}
}

func (lnd *LndClient) Pay(cx context.Context, invoice Invoice) (lntypes.Preimage, error) {
	encoded_pr := lnrpc.SendRequest{PaymentRequest: invoice.GetPaymentRequest()}

	response, err := lnd.client.SendPaymentSync(cx, &encoded_pr)

	if err != nil {
		return lntypes.Preimage{}, err
	}

	pre_image, _ := lntypes.MakePreimage(response.PaymentPreimage)

	return pre_image, nil
}

func (lnd *LndClient) CreateInvoice(cx context.Context, invoice InvoiceBuilder) (Invoice, error) {
	invoice_resp, err := lnd.client.AddInvoice(cx, &invoice, grpc.EmptyCallOption{})
	return *invoice_resp, err
}
