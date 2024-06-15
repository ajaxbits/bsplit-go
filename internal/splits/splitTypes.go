package splits

type SplitType int

const (
	Even = iota
	Percent
	Adjustment
	Exact
)

var splitTypeName = map[SplitType]string{
	Even:       "even",
	Percent:    "percent",
	Adjustment: "adjustment",
	Exact:      "exact",
}
