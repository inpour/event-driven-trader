package trade

type Side string

const (
	Sell Side = "sell"
	Buy  Side = "buy"
)

type Position struct {
	Id                  int
	Enter, Sl, Tp, Exit float64
	Side                Side
}
