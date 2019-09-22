package memdb

import (
	"fmt"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source ./segment_store.go -destination=./segment_store_mock_test.go -package memdb

const (
	emptySimpleFieldStoreSize = 8 + // familyTime
		8 + // aggFunc
		8 // block pointer
)

// sStoreINTF represents segment-store,
// which abstracts a store for storing field data based on family start time
type sStoreINTF interface {
	GetFamilyTime() int64

	AggType() field.AggType

	SlotRange() (
		startSlot,
		endSlot int,
		err error)

	Bytes(
		needSlotRange bool,
	) (
		data []byte,
		startSlot,
		endSlot int,
		err error)

	// WriteInt writes a int value, and returns the written length
	WriteInt(
		value int64,
		writeCtx writeContext,
	) int

	// WriteFloat writes a float64 value, and returns the written length
	WriteFloat(value float64,
		writeCtx writeContext,
	) int

	MemSize() int

	// scan scans segment store data based on query time range
	scan(agg aggregation.SeriesAggregator, memScanCtx *memScanContext)
}

// singleFieldStore stores single field
type simpleFieldStore struct {
	familyTime int64
	aggFunc    field.AggFunc
	block      block
}

// newSingleFieldStore returns a new segment store for simple field store
func newSimpleFieldStore(familyTime int64, aggFunc field.AggFunc) sStoreINTF {
	return &simpleFieldStore{
		familyTime: familyTime,
		aggFunc:    aggFunc,
	}
}

func (fs *simpleFieldStore) GetFamilyTime() int64 {
	return fs.familyTime
}

func (fs *simpleFieldStore) AggType() field.AggType {
	return fs.aggFunc.AggType()
}

func (fs *simpleFieldStore) WriteFloat(value float64, writeCtx writeContext) int {
	oldSize := fs.MemSize()
	pos, hasValue := fs.calcTimeWindow(writeCtx.blockStore, writeCtx.slotIndex, field.Float)
	currentBlock := fs.block
	if hasValue {
		// do rollup using agg func
		currentBlock.setFloatValue(pos, fs.aggFunc.AggregateFloat(currentBlock.getFloatValue(pos), value))
	} else {
		currentBlock.setFloatValue(pos, value)
	}
	return fs.MemSize() - oldSize
}

func (fs *simpleFieldStore) WriteInt(value int64, writeCtx writeContext) int {
	oldSize := fs.MemSize()
	pos, hasValue := fs.calcTimeWindow(writeCtx.blockStore, writeCtx.slotIndex, field.Integer)
	currentBlock := fs.block
	if hasValue {
		// do rollup using agg func
		currentBlock.setIntValue(pos, fs.aggFunc.AggregateInt(currentBlock.getIntValue(pos), value))
	} else {
		currentBlock.setIntValue(pos, value)
	}
	return fs.MemSize() - oldSize
}

// calcTimeWindow calculates time window's block for storing field data based on slot time and value type.
// return int=>pos(slot in time window), bool=>needRollup(if rollup with old value)
// 1) block is nil, create new block, return 0, false
// 2) slot time out of current time window, need compress time window then create new one, return 0, false
// 3) in current time window, if has old value return pos, true, else return pos, false
func (fs *simpleFieldStore) calcTimeWindow(blockStore *blockStore, slotTime int,
	valueType field.ValueType) (int, bool) {
	currentBlock := fs.block

	// block is nil
	if currentBlock == nil {
		currentBlock = blockStore.allocBlock(valueType)
		currentBlock.setStartTime(slotTime)
		fs.block = currentBlock
		return 0, false
	}

	startTime := currentBlock.getStartTime()

	// if current slot time out of current time window, need compress block data, start new time window
	if slotTime < startTime || slotTime >= startTime+blockStore.timeWindow {
		_, _, err := currentBlock.compact(fs.aggFunc)
		if err != nil {
			memDBLogger.Error("compress block data error, data will lost", logger.Error(err))
		} else {
			// reset start time using slot time
			currentBlock.setStartTime(slotTime)
		}
		return 0, false
	}

	// in current time window, do rollup value
	pos := slotTime - startTime
	needRollup := false
	if currentBlock.hasValue(pos) {
		// has old value, need do rollup
		needRollup = true
	}
	return pos, needRollup
}

func (fs *simpleFieldStore) Bytes(needSlotRange bool) (data []byte, startSlot, endSlot int, err error) {
	if fs.block == nil {
		err = fmt.Errorf("block is empty")
		return
	}
	if startSlot, endSlot, err = fs.block.compact(fs.aggFunc); err != nil {
		err = fmt.Errorf("compact block data in simple field store error:%s", err)
		return
	}
	data = fs.block.bytes()
	return
}

func (fs *simpleFieldStore) SlotRange() (startSlot, endSlot int, err error) {
	if fs.block == nil {
		err = fmt.Errorf("block is empty")
		return
	}
	startSlot, endSlot = encoding.DecodeTSDTime(fs.block.bytes())
	return
}

func (fs *simpleFieldStore) MemSize() int {
	if fs.block == nil {
		return emptySimpleFieldStoreSize
	}
	return emptySimpleFieldStoreSize + fs.block.memsize()
}
