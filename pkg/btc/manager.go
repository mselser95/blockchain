package btc

// import (
// 	"bytes"
// 	"encoding/json"
// 	"errors"
// 	"math/big"
// 	"net/http"

// 	"github.com/mselser95/blockchain/pkg/signer"
// 	"github.com/mselser95/blockchain/pkg/utils"
// )

// // Manager is a manager for interacting with the Bitcoin blockchain.
// type Manager struct {
// 	rpcURL string
// 	signer signer.TransactionSigner
// }

// // Connect establishes a connection to the Bitcoin node.
// func (m *Manager) Connect(url string) error {
// 	m.rpcURL = url
// 	return nil
// }

// // GetBalance retrieves the balance of the specified address.
// func (m *Manager) GetBalance(address string, token *utils.Token) (*big.Int, error) {
// 	if token.Type != utils.Native {
// 		return nil, errors.New("unsupported token type for Bitcoin")
// 	}

// 	// Bitcoin does not directly support querying balance by address via JSON-RPC.
// 	// You may need to implement this using a third-party service or additional logic.
// 	return nil, errors.New("not implemented")
// }

// // ReadCall performs a read-only call on the Bitcoin blockchain.
// func (m *Manager) ReadCall(tx *utils.Transaction) (interface{}, error) {
// 	// Bitcoin does not typically support read-only contract calls.
// 	return nil, errors.New("unsupported method")
// }

// // SendTransaction signs and sends a transaction to the Bitcoin network.
// func (m *Manager) SendTransaction(tx *utils.Transaction) (string, error) {
// 	if m.signer == nil {
// 		return "", errors.New("transaction signer is not set")
// 	}

// 	// Sign the transaction using the signer
// 	signedTx, err := m.signer.SignTransaction(tx)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Send the signed transaction to the Bitcoin network via RPC
// 	txHex := signedTx.Payload["signedTx"].(string)
// 	reqBody, err := json.Marshal(map[string]interface{}{
// 		"jsonrpc": "1.0",
// 		"id":      "curltest",
// 		"method":  "sendrawtransaction",
// 		"params":  []interface{}{txHex},
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	resp, err := http.Post(m.rpcURL, "application/json", bytes.NewBuffer(reqBody))
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	var respData map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
// 		return "", err
// 	}

// 	if respData["error"] != nil {
// 		return "", errors.New(respData["error"].(map[string]interface{})["message"].(string))
// 	}

// 	return respData["result"].(string), nil
// }

// // GetTransactionDetails retrieves the details of a transaction by its ID.
// func (m *Manager) GetTransactionDetails(txID string) (*utils.TransactionDetails, error) {
// 	reqBody, err := json.Marshal(map[string]interface{}{
// 		"jsonrpc": "1.0",
// 		"id":      "curltest",
// 		"method":  "getrawtransaction",
// 		"params":  []interface{}{txID, 1},
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	resp, err := http.Post(m.rpcURL, "application/json", bytes.NewBuffer(reqBody))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var respData map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
// 		return nil, err
// 	}

// 	if respData["error"] != nil {
// 		return nil, errors.New(respData["error"].(map[string]interface{})["message"].(string))
// 	}

// 	txData := respData["result"].(map[string]interface{})

// 	// Process the transaction data and convert it to the TransactionDetails format
// 	// Here you would need to parse the data according to the BTC JSON-RPC format
// 	details := &utils.TransactionDetails{
// 		ID:          txID,
// 		Status:      "confirmed", // Bitcoin doesn't directly provide status, assumed confirmed
// 		BlockNumber: uint64(txData["blockheight"].(float64)),
// 		Timestamp:   txData["time"].(float64), // Convert to time.Time if needed
// 		From:        "",                       // Parsing inputs to determine the sender
// 		To:          "",                       // Parsing outputs to determine the recipient
// 		Amount:      big.NewInt(int64(txData["vout"].([]interface{})[0].(map[string]interface{})["value"].(float64) * 1e8)),
// 		Fee:         nil, // Bitcoin fees can be calculated from the difference between inputs and outputs
// 		Logs:        nil, // No logs in Bitcoin transactions
// 		Events:      nil, // No events in Bitcoin transactions
// 	}

// 	return details, nil
// }

// // Start starts the Manager. This could involve setting up necessary resources.
// func (m *Manager) Start() error {
// 	// Implement any startup logic if needed
// 	return nil
// }

// // Stop stops the Manager and cleans up resources.
// func (m *Manager) Stop() error {
// 	// Implement any shutdown logic if needed
// 	return nil
// }
