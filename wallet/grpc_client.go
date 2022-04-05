package wallet

import (
	"context"
	"encoding/hex"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	client zarb.ZarbClient
}

func MewGRPCClient(rpcEndpoint string) (*GrpcClient, error) {
	conn, err := grpc.Dial(rpcEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &GrpcClient{
		client: zarb.NewZarbClient(conn),
	}, nil
}

func (c *GrpcClient) GetStamp() (hash.Stamp, error) {
	info, err := c.client.GetBlockchainInfo(context.Background(), &zarb.BlockchainInfoRequest{})
	if err != nil {
		return hash.Stamp{}, err
	}
	h, _ := hash.FromBytes(info.LastBlockHash)
	return h.Stamp(), nil
}

func (c *GrpcClient) GetAccountBalance(addr crypto.Address) (int64, error) {
	acc, err := c.client.GetAccount(context.Background(), &zarb.AccountRequest{Address: addr.Bytes()})
	if err != nil {
		return 0, err
	}

	return acc.Account.Balance, nil
}

func (c *GrpcClient) GetAccountSequence(addr crypto.Address) (int32, error) {
	acc, err := c.client.GetAccount(context.Background(), &zarb.AccountRequest{Address: addr.Bytes()})
	if err != nil {
		return 0, err
	}

	return acc.Account.Sequence + 1, nil
}

func (c *GrpcClient) GetValidatorSequence(addr crypto.Address) (int32, error) {
	val, err := c.client.GetValidator(context.Background(), &zarb.ValidatorRequest{Address: addr.Bytes()})
	if err != nil {
		return 0, err
	}

	return val.Validator.Sequence + 1, nil
}

func (c *GrpcClient) GetValidatorStake(addr crypto.Address) (int64, error) {
	val, err := c.client.GetValidator(context.Background(), &zarb.ValidatorRequest{Address: addr.Bytes()})
	if err != nil {
		return 0, err
	}

	return val.Validator.Stake, nil
}

func (c *GrpcClient) SendTx(payload []byte) (string, error) {
	res, err := c.client.SendRawTransaction(context.Background(), &zarb.SendRawTransactionRequest{
		Data: hex.EncodeToString(payload),
	})

	if err != nil {
		return "", err
	}

	return res.Id, nil
}
