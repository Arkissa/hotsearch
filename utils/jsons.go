package utils

import (
	"encoding/json"
	"hotsearch/log"
)

type JsonDate struct {
	Date   []byte
    Decode any
}

func (j *JsonDate) Decoder() {

    if err := json.Unmarshal(j.Date, j.Decode); err != nil {
        log.LogOutErr("json unmarshal err", err)
    }
}
