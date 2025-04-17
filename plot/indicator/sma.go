package indicator

import (
	"fmt"
	"github.com/markcheno/go-talib"
	"time"

	"github.com/hunternsk/ninjabot/model"
	"github.com/hunternsk/ninjabot/plot"
)

func SMA(period int, color string) plot.Indicator {
	return &sma{
		Period: period,
		Color:  color,
	}
}

type sma struct {
	Period int
	Color  string
	Values model.Series[float64]
	Time   []time.Time
}

func (s sma) Warmup() int {
	return s.Period
}

func (s sma) Name() string {
	return fmt.Sprintf("SMA(%d)", s.Period)
}

func (s sma) Overlay() bool {
	return true
}

func (s *sma) Load(dataframe *model.Dataframe) {
	if len(dataframe.Time) < s.Period {
		return
	}

	s.Values = talib.Sma(dataframe.Close, s.Period)[s.Period:]
	s.Time = dataframe.Time[s.Period:]
}

func (s sma) Metrics() []plot.IndicatorMetric {
	return []plot.IndicatorMetric{
		{
			Style:  "line",
			Color:  s.Color,
			Values: s.Values,
			Time:   s.Time,
		},
	}
}
