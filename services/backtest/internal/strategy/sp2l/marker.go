package sp2l

import (
	"fmt"
	"strconv"

	"github.com/inpour/event-driven-trader/shared/backtest"
)

type markerType int

const (
	gapMarker markerType = iota
	buyMarker
	sellMarker
	tpMarker
	slMarker
	MaxEntranceMarker
)

var markerPosition = map[markerType]string{
	gapMarker:         "belowBar",
	buyMarker:         "belowBar",
	sellMarker:        "aboveBar",
	tpMarker:          "aboveBar",
	slMarker:          "aboveBar",
	MaxEntranceMarker: "aboveBar",
}

var markerShape = map[markerType]string{
	gapMarker:         "circle",
	buyMarker:         "arrowUp",
	sellMarker:        "arrowDown",
	tpMarker:          "square",
	slMarker:          "square",
	MaxEntranceMarker: "circle",
}

var markerColor = map[markerType]string{
	gapMarker:         "#e2c100",
	buyMarker:         "#00c853",
	sellMarker:        "#ff1744",
	tpMarker:          "#00c853",
	slMarker:          "#ff1744",
	MaxEntranceMarker: "#ffffff",
}

var markerSize = map[markerType]float32{
	gapMarker:         0.1,
	buyMarker:         0.3,
	sellMarker:        0.3,
	tpMarker:          0.3,
	slMarker:          0.3,
	MaxEntranceMarker: 0.1,
}

func marker(markerType_ markerType, id int, time int64, price float64) *backtest.Marker {
	text := "ID:" + strconv.Itoa(id)
	if markerType_ == MaxEntranceMarker {
		text += " Max Distance"
	} else if markerType_ != gapMarker {
		text += " @ " + fmt.Sprintf("%.2f", price)
	}

	return &backtest.Marker{
		Time:     time,
		Price:    price,
		Position: markerPosition[markerType_],
		Shape:    markerShape[markerType_],
		Color:    markerColor[markerType_],
		Size:     markerSize[markerType_],
		Text:     text,
	}
}
