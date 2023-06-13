package kraken

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// UnmarshalJSON - unmarshal update
func (u *DataUpdate) UnmarshalJSON(data []byte) error {
	var raw []interface{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	if len(raw) < 3 {
		return fmt.Errorf("invalid data length: %#v", raw)
	}

	if len(raw) == 3 {
		var ok bool
		u.Data = raw[0]

		if u.ChannelName, ok = raw[1].(string); !ok {
			return fmt.Errorf("expected message to have channel name as 2nd element but got %#v instead", raw[1])
		}

		var sequenceMap map[string]interface{}
		if sequenceMap, ok = raw[2].(map[string]interface{}); !ok {
			return fmt.Errorf("expected message to have JSON object as 3rd element but got %#v instead", raw[2])
		}

		var sequenceRaw interface{}
		if sequenceRaw, ok = sequenceMap["sequence"]; !ok {
			return fmt.Errorf("expected message to have sequence in JSON object as 3rd element but got %#v instead", raw[2])
		}

		var seq float64
		if seq, ok = sequenceRaw.(float64); !ok {
			return fmt.Errorf("expected message to have sequence integer in JSON object as 3rd element but got %#v instead", raw[2])
		}

		u.Sequence = int64(seq)
		return nil
	}

	chID, ok := raw[0].(float64)
	if !ok {
		return fmt.Errorf("expected message to start with a channel id but got %#v instead", raw[0])
	}

	u.ChannelID = int64(chID)
	u.ChannelName, ok = raw[len(raw)-2].(string)
	if !ok {
		return fmt.Errorf("expected message with (n - 2) element channel name but got %#v instead", raw[len(raw)-2])
	}
	u.Pair, ok = raw[len(raw)-1].(string)
	if !ok {
		return fmt.Errorf("expected message  with (n - 2) element pair but got %#v instead", raw[len(raw)-1])
	}
	u.Data = raw[1 : len(raw)-2][0]

	return nil
}

// ParseData parse a sub data from the ws update msg
func ParseData(data interface{}, pair string, dataParser *DataParser) (interface{}, error) {
	dataParser.Result.Pair = pair

	dataParser.body, dataParser.ok = data.(map[string]interface{})
	if !dataParser.ok {
		return dataParser.Result, fmt.Errorf("Can't parse data %#v", data)
	}
	dataParser.Result.Asks = dataParser.Result.Asks[:0]
	dataParser.Result.Bids = dataParser.Result.Bids[:0]

	for dataParser.k, dataParser.v = range dataParser.body {
		switch dataParser.k {
		case "c":
			dataParser.checkSum, dataParser.ok = dataParser.v.(string)
			if !dataParser.ok {
				return nil, fmt.Errorf("[bookFactory] Invalid checkSum type: %v %T", dataParser.v, dataParser.v)
			}
			dataParser.Result.CheckSum = dataParser.checkSum
		default:
			dataParser.itemsParser.items = dataParser.itemsParser.items[:0]
			dataParser.err = parseItems(dataParser.v, &dataParser.itemsParser)
			if dataParser.err != nil {
				return nil, dataParser.err
			}
			dataParser.Result.IsSnapshot = len(dataParser.k) == 2 && strings.HasSuffix(dataParser.k, "s")
			if strings.HasPrefix(dataParser.k, "a") {
				dataParser.Result.Asks = dataParser.itemsParser.items
			} else {
				dataParser.Result.Bids = dataParser.itemsParser.items
			}
		}
	}
	return dataParser.Result, nil
}

func parseItems(value interface{}, itemsParser *ItemsParser) error {

	itemsParser.updates, itemsParser.ok = value.([]interface{})
	if !itemsParser.ok {
		return fmt.Errorf("[bookFactory] Invalid items type: %v %T", value, value)
	}
	for _, itemsParser.item = range itemsParser.updates {
		itemsParser.entity = itemsParser.item.([]interface{})

		itemsParser.orderBookItem.Price = itemsParser.entity[0].(string)
		itemsParser.orderBookItem.Volume = itemsParser.entity[1].(string)
		// itemsParser.orderBookItem.Time = valToFloat64(itemsParser.entity[2])
		// itemsParser.orderBookItem.Republish = (len(itemsParser.entity) == 4 && itemsParser.entity[3] == "r")

		itemsParser.items = append(itemsParser.items, itemsParser.orderBookItem)
	}
	return nil
}

func valToFloat64(value interface{}) float64 {
	if v, ok := value.(string); ok {
		result, err := strconv.ParseFloat(v, 64)
		if err != nil {
			log.Printf("Can't parse float %#v", value)
			return .0
		}
		return result
	}
	return .0
}
