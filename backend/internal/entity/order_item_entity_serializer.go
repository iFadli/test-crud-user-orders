package entity

import "encoding/json"

func MarshalOrderItem(orderItem *OrderItem) ([]byte, error) {
	return json.Marshal(orderItem)
}

func UnmarshalOrderItem(data []byte, orderItem *OrderItem) error {
	err := json.Unmarshal(data, orderItem)
	if err != nil {
		return err
	}
	return nil
}

func MarshalOrderItems(orderItems []*OrderItem) ([]byte, error) {
	return json.Marshal(orderItems)
}

func UnmarshalOrderItems(data []byte, orderItems *[]*OrderItem) error {
	err := json.Unmarshal(data, orderItems)
	if err != nil {
		return err
	}
	return nil
}
