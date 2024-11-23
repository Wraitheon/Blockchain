package blockchain

import "fmt"

type Transaction struct {
	Sender    string
	Recipient string
	Amount    float64
}

func CreateTransaction(sender string, recipient string, amount float64) Transaction {
	return Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
}

func SerializeTransaction(transactions []Transaction) []string {
	var result []string
	for _, tx := range transactions {
		record := tx.Sender + "->" + tx.Recipient + ": " + fmt.Sprintf("%.2f", tx.Amount)
		result = append(result, record)
	}

	return result
}
