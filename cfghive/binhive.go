package cfghive

import (
	"bufio"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"github.com/hashicorp/go-msgpack/codec"
	"io"
	"reflect"
)

var (
	mh codec.MsgpackHandle
)

func configureCodec() *codec.MsgpackHandle {
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))
	return &mh
}

// BinHive is a hive built on top of MemHive that can be serialized and loaded to/from a writer.
type BinHive struct {
	hive       *MemHive
	hasChange  bool
	isCommited bool
	Stream     *bufio.ReadWriter
	comp       bool
	compLevel  uint8
}

func NewBinHive(compression bool, level uint8) *BinHive {
	h, _ := NewMemHive()
	return &BinHive{
		hive:       h,
		hasChange:  false,
		isCommited: false,
		comp:       compression,
		compLevel:  level,
	}
}

func (h *BinHive) Characteristics() HiveCharacteristics {
	return HiveCharacteristics{false, true, false}
}

func (h *BinHive) Set(key string, value interface{}) error {
	h.hasChange = true
	return h.hive.Set(key, value)
}

func (h *BinHive) SetBool(key string, value bool) error {
	h.hasChange = true
	return h.hive.SetBool(key, value)
}

func (h *BinHive) SetInt(key string, value int) error {
	h.hasChange = true
	return h.hive.SetInt(key, value)
}

func (h *BinHive) SetFloat(key string, value float64) error {
	h.hasChange = true
	return h.hive.SetFloat(key, value)
}

func (h *BinHive) SetString(key string, value string) error {
	h.hasChange = true
	return h.hive.SetString(key, value)
}

func (h *BinHive) Get(key string) (interface{}, error) {
	return h.hive.Get(key)
}

func (h *BinHive) GetBool(key string) (bool, error) {
	return h.hive.GetBool(key)
}

func (h *BinHive) GetInt(key string) (int, error) {
	return h.hive.GetInt(key)
}

func (h *BinHive) GetFloat(key string) (float64, error) {
	return h.hive.GetFloat(key)
}

func (h *BinHive) GetString(key string) (*string, error) {
	return h.hive.GetString(key)
}

func (h *BinHive) NewSub(key string) {
	h.hasChange = true
	h.hive.NewSub(key)
}

func (h *BinHive) Delete(key string) {
	h.hasChange = true
	h.hive.Delete(key)
}

func (h *BinHive) Commit() error {
	if !h.hasChange {
		return nil
	}
	if h.isCommited {
		return nil
	}
	h.isCommited = true
	h.hasChange = false
	return h.saveToWriter(*h.Stream)
}

func (h *BinHive) Load() error {
	return h.loadFromReader(*h.Stream)
}

func (h *BinHive) Save() error {
	return h.saveToWriter(*h.Stream)
}

func (h *BinHive) GetData() *map[string]HiveValue {
	return h.hive.GetData()
}

func (h *BinHive) saveToWriter(w io.Writer) error {
	len_ := uint64(HiveSize(h.hive.data))
	dataRaw := HiveMapToGeneric(h.hive.data)
	fmt.Printf("hive size: %d\n", len_)
	if h.comp {
		header := []byte{0xC1, h.compLevel}
		header = binary.BigEndian.AppendUint64(header, len_)
		_, err := w.Write(header)
		if err != nil {
			return err
		}
		cw, err := gzip.NewWriterLevel(w, int(h.compLevel))
		if err != nil {
			return err
		}
		defer cw.Close()
		ch := configureCodec()
		c := codec.NewEncoder(cw, ch)
		err = c.Encode(dataRaw)
		err = cw.Flush()
		if err != nil {
			return err
		}
		return err
	} else {
		header := []byte{0xC0, 0x00}
		header = binary.BigEndian.AppendUint64(header, len_)
		_, err := w.Write(header)
		if err != nil {
			return err
		}
		ch := configureCodec()
		c := codec.NewEncoder(w, ch)
		err = c.Encode(dataRaw)
		return err
	}
}

func (h *BinHive) loadFromReader(r io.Reader) error {
	header := make([]byte, 10)
	i, err := r.Read(header)
	if err != nil {
		return err
	}
	if i != 10 {
		return fmt.Errorf("invalid header length: %d", i)
	}
	if header[0] != 0xC0 && header[0] != 0xC1 {
		return fmt.Errorf("invalid header byte: %x", header[0])
	}
	if header[0] == 0xC1 {
		h.comp = true
		h.compLevel = header[1]
	} else {
		h.comp = false
		h.compLevel = 0
	}
	rawData := make(map[string]interface{})
	if h.comp {
		cr, err := gzip.NewReader(r)
		if err != nil {
			return err
		}
		defer cr.Close()
		ch := configureCodec()
		c := codec.NewDecoder(cr, ch)
		err = c.Decode(&rawData)
		if err != nil {
			return err
		}

	} else {
		ch := configureCodec()
		c := codec.NewDecoder(r, ch)
		err = c.Decode(&rawData)
		if err != nil {
			return err
		}
	}
	data, err := GenericMapToSubMap(rawData)
	if err != nil {
		return err
	}
	h.hive.data = data
	return nil
}
