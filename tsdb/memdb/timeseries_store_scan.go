package memdb

import (
	"github.com/lindb/lindb/series"
)

// scan scans time series data, then finds field store by field id
func (ts *timeSeriesStore) scan(sCtx *series.ScanContext, version series.Version, seriesID uint32, fieldMetas map[uint16]*fieldMeta) {
	worker := sCtx.Worker
	for _, fieldID := range sCtx.FieldIDs {
		ts.sl.Lock()
		fStore, ok := ts.getFStore(fieldID)
		ts.sl.Unlock()
		if !ok {
			continue
		}
		fieldMeta := fieldMetas[fieldID]
		fStore.scan(sCtx, version, seriesID, fieldMeta, ts)
	}

	// send msg to notify current series scan completed
	worker.Complete(seriesID)
}
