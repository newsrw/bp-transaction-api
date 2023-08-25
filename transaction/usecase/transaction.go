package usecase

import (
	"bp-transaction-api/configs"
	"bp-transaction-api/domain"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	requestTimeoutLimit = 30 * time.Second
)

type transactionUsecase struct {
	contextTimeout time.Duration
	config         *configs.Config
	logger         *zap.Logger
}

// NewTransactionUsecase create new an transactionUsecase object representation of domain.TransactionUsecase interface
func NewTransactionUsecase(
	timeout time.Duration,
	config *configs.Config,
	logger *zap.Logger,
) domain.TransactionUsecase {
	return &transactionUsecase{
		contextTimeout: timeout,
		config:         config,
		logger:         logger,
	}
}

// BroadcastTransaction broadcast transaction
func (uc *transactionUsecase) BroadcastTransaction(c context.Context, req *domain.BroadcastTransactionRequest) (rsp domain.BroadcastTransactionResponse, err error) {
	_, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	requestURL := fmt.Sprintf("%s/broadcast", uc.config.Server.Address)
	var jsonData []byte
	if jsonData, err = json.Marshal(req); err != nil {
		uc.logger.Error("marshal json error", zap.Error(err))
		return
	}

	broadcastReq, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(jsonData))
	if err != nil {
		uc.logger.Error("could not create request", zap.Error(err))
		return
	}
	broadcastReq.Header.Add("Content-Type", "application/json")

	client := http.Client{Timeout: requestTimeoutLimit}
	res, err := client.Do(broadcastReq)
	if err != nil {
		uc.logger.Error("error making http request", zap.Error(err))
		return
	}
	defer res.Body.Close()

	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		uc.logger.Error("could not read response body", zap.Error(err))
		return
	}

	err = json.Unmarshal(resBody, &rsp)
	if err != nil {
		uc.logger.Error("unmarshal json error", zap.Error(err))
		return
	}

	return
}

// MonitoringStatus monitoring transaction status
func (uc *transactionUsecase) MonitoringStatus(c context.Context, req *domain.MonitoringTransactionStatusRequest) (rsp *domain.MonitoringTransactionStatusResponse, err error) {
	_, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	if req.TxHash == "" {
		err = domain.ErrInternalServerError
		uc.logger.Error("tx_hash is empty", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("==== Start Transaction Status Monitoring ====")
	for {
		rsp := &domain.MonitoringTransactionStatusResponse{}
		requestURL := fmt.Sprintf("%s/check/%s", uc.config.Server.Address, req.TxHash)
		checkStatusReq, err := http.NewRequest(http.MethodGet, requestURL, nil)
		if err != nil {
			uc.logger.Error("could not create request", zap.Error(err))
			return nil, err
		}

		client := http.Client{Timeout: requestTimeoutLimit}
		httpRsp, err := client.Do(checkStatusReq)
		if err != nil {
			uc.logger.Error("error making http request", zap.Error(err))
			return nil, err
		}
		defer httpRsp.Body.Close()

		resBody, err := io.ReadAll(httpRsp.Body)
		if err != nil {
			uc.logger.Error("could not read response body", zap.Error(err))
			return nil, err
		}

		err = json.Unmarshal(resBody, rsp)
		if err != nil {
			uc.logger.Error("unmarshal json error", zap.Error(err))
			return nil, err
		}

		currentStatus := rsp.TxStatus
		uc.logger.Info(
			"Monitoring Info",
			zap.String("Tx Hash", req.TxHash),
			zap.Any("Status", currentStatus),
		)

		if currentStatus == domain.TxStatusConfirmed || currentStatus == domain.TxStatusFailed || currentStatus == domain.TxStatusDNE {
			return rsp, nil
		}

		fiveSecondsDuration := time.Second * time.Duration(uc.config.Client.MonitorDelayDuration)
		time.Sleep(fiveSecondsDuration)
	}
}
