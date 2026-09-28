package logging

import (
	"fmt"
	"time"

	kafkago "github.com/segmentio/kafka-go"

	"github.com/inpour/event-driven-trader/shared/events"
)

// LogCandleReceived writes a readable record of a received CandleClosed event.
func LogCandleReceived(
	event events.Envelope[events.CandleClosed],
	message kafkago.Message,
) {
	candle := event.Payload.Candle

	fmt.Printf(
		"received candle-close event: "+
			"time=%s "+
			"run_id=%s "+
			"event_id=%s "+
			"symbol=%s "+
			"timeframe=%s "+
			"close_time=%d "+
			"close=%f "+
			"partition=%d "+
			"offset=%d\n",
		time.Now().UTC(),
		event.RunID,
		event.EventID,
		event.Payload.Symbol,
		event.Payload.Timeframe,
		candle.Time,
		candle.Close,
		message.Partition,
		message.Offset,
	)
}
