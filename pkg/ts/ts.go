package ts

import (
	"time"
	"wfmon/pkg/utils/vec"
)

type Sample struct {
	Value     float64
	Timestamp time.Time
}

// TimeSeries represents samples and labels for a single time series.
type TimeSeries struct {
	MaxLen  uint
	Samples []Sample
}

func Empty() TimeSeries {
	return New(0)
}

func New(maxLen uint) TimeSeries {
	return TimeSeries{
		MaxLen:  maxLen,
		Samples: make([]Sample, 0, maxLen),
	}
}

func (ts TimeSeries) Add(val float64, timestamp time.Time) TimeSeries {
	sample := Sample{
		Value:     val,
		Timestamp: time.Now(),
	}
	ts.Samples = append(ts.Samples, sample)

	ts.Shrink()

	return ts
}

func (ts TimeSeries) Shrink() TimeSeries {
	start := len(ts.Samples) - int(ts.MaxLen)
	if start < 0 {
		return ts
	}

	ts.Samples = ts.Samples[start:]

	return ts
}

func (ts TimeSeries) Copy() TimeSeries {
	samples := make([]Sample, len(ts.Samples))
	copy(samples, ts.Samples)

	return TimeSeries{
		MaxLen:  ts.MaxLen,
		Samples: samples,
	}
}

func (ts TimeSeries) end(cnt int) []Sample {
	start := len(ts.Samples) - cnt
	if start < 0 {
		start = 0
	}
	return ts.Samples[start:]
}

type Vector []float64

func (ts TimeSeries) Last() (float64, bool) {
	vals := ts.Range(1)
	if len(vals) > 0 {
		return vals[0], true
	}

	return 0, false
}

func (ts TimeSeries) Range(cnt int) Vector {
	return ts.RangeAndMap(cnt, func(val float64) float64 { return val })
}

func (ts TimeSeries) RangeAndMap(cnt int, modifier func(val float64) float64) Vector {
	samples := ts.end(cnt)

	vec := make(Vector, len(samples))
	for i := 0; i < len(samples); i++ {
		vec[i] = modifier(samples[i].Value)
	}

	return vec
}

func (v Vector) Reverse() Vector {
	return vec.Reverse(v)
}
