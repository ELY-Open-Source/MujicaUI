package finance

import (
	"fmt"
	"log/slog"

	"github.com/egoist/mygo/ui"

	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/core"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/data"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/input"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/navigation"
	"git.zacharyzhang.com/ZacharyZhang-NY/MujicaUI/theme"
)

// OrderType is limit, market, stop or stop limit.
type OrderType int

const (
	OrderLimit OrderType = iota
	OrderMarket
	OrderStop
	OrderStopLimit
)

// TimeInForce is how long an order rests.
type TimeInForce int

const (
	TIFDay TimeInForce = iota
	TIFGoodTillCanceled
	TIFImmediateOrCancel
	TIFFillOrKill
)

// OrderDraft is an order as entered.
type OrderDraft struct {
	Side     OrderSide
	Type     OrderType
	Quantity float64
	Price    float64
	Trigger  float64
	TIF      TimeInForce
}

// OrderEntryOptions configure an OrderEntry. Symbol is required.
type OrderEntryOptions struct {
	Symbol string
	// Last is the last price; zero means unknown.
	Last float64
	// TickSize and LotSize step prices and quantities; zero: 0.01, 1.
	TickSize, LotSize float64
	// MaxQuantity caps the quantity; zero means no cap.
	MaxQuantity float64
	// Initial is what Reset restores; a zero Side means buy.
	Initial OrderDraft
	// Validate adds the caller's rules, returning a message or "".
	Validate func(o OrderDraft) string
	// Check starts an async check, finished with OrderEntryState.Resolve.
	Check    func(o OrderDraft, version uint64)
	Disabled bool
}

// OrderEntryState holds the draft and its form state.
type OrderEntryState struct {
	Draft   OrderDraft
	Form    input.FormState
	started bool
}

// Resolve delivers Check's result; stale versions return false.
func (s *OrderEntryState) Resolve(version uint64, err string) bool {
	return s.Form.Resolve("quantity", version, err)
}

// EndSubmit finishes a submit; a formError shows, keeping the draft.
func (s *OrderEntryState) EndSubmit(formError string) { s.Form.EndSubmit(formError) }

// OrderEntryResult is a built OrderEntry and its submitted order.
type OrderEntryResult struct {
	Element   *ui.Element
	order     OrderDraft
	submitted bool
}

// Submitted returns the checked order; call EndSubmit once placed.
func (r OrderEntryResult) Submitted() (OrderDraft, bool) { return r.order, r.submitted }

// OrderEntry edits an order and submits it once checked.
func OrderEntry(c *ui.Context, s *OrderEntryState, opts OrderEntryOptions) OrderEntryResult {
	if opts.Symbol == "" {
		panic("mujica: OrderEntry needs a Symbol")
	}
	if opts.TickSize < 0 || opts.LotSize < 0 || opts.MaxQuantity < 0 || opts.Last < 0 {
		panic("mujica: OrderEntry TickSize, LotSize, MaxQuantity and Last cannot be negative")
	}
	tick, lot := opts.TickSize, opts.LotSize
	if tick == 0 {
		tick = 0.01
	}
	if lot == 0 {
		lot = 1
	}
	initial := opts.Initial
	if initial.Side == 0 {
		initial.Side = SideBuy
	}
	if !s.started {
		s.started, s.Draft = true, initial
	}
	priceSpec := input.NumberInputResolve(input.NumberInputOptions{Step: tick})
	k := core.Tokens(c)
	u := core.Density(c).Unit()
	var r OrderEntryResult
	r.Element = ui.Column(c).Gap(u * 4).MaxWidth(360).Disabled(opts.Disabled).Role(ui.RoleGroup).Label(core.Msg(c, "order_entry.label"))
	disabled := r.Element.IsDisabled()
	r.Element.Children(func() {
		side := 0
		if s.Draft.Side == SideSell {
			side = 1
		}
		sides := []navigation.Segment{{Label: core.Msg(c, "order_entry.buy")}, {Label: core.Msg(c, "order_entry.sell")}}
		navigation.SegmentedControl(c, &side, sides, navigation.SegmentedControlOptions{Equal: true, Label: core.Msg(c, "order_entry.side"), Disabled: disabled})
		s.Draft.Side = SideBuy + OrderSide(side)
		var types []input.SelectOption[OrderType]
		for t := OrderLimit; t <= OrderStopLimit; t++ {
			types = append(types, input.SelectOption[OrderType]{Value: t, Label: core.Msg(c, fmt.Sprintf("order_entry.type%d", t))})
		}
		ui.Row(c).Gap(u * 3).AlignItems(ui.Center).Children(func() {
			ui.Text(c, opts.Symbol).Font(theme.SerifFont).FontSize(core.FontSize(c, theme.H3Size)).FontWeight(600).Grow(1)
			input.Select(c, &s.Draft.Type, types, input.SelectOptions{Label: core.Msg(c, "order_entry.type"), Disabled: disabled})
		})
		submit := fmt.Sprintf(core.Msg(c, fmt.Sprintf("order_entry.submit%d", s.Draft.Side)), opts.Symbol)
		res := input.Form(c, &s.Form, input.FormOptions{SubmitLabel: submit}, func(f *input.FormScope) {
			orderEntryFields(c, f, s, opts, initial, tick, lot, disabled)
		})
		ui.Row(c).Gap(u * 2).AlignItems(ui.Center).Children(func() {
			ui.Text(c, core.Msg(c, "order_entry.estimate")).TextColor(k.TextMuted)
			est := "—"
			if p := orderEntryPrice(s.Draft, opts.Last); p > 0 && s.Draft.Quantity > 0 {
				est = "≈ " + data.StatisticFormat(p*s.Draft.Quantity, priceSpec.Decimals)
			}
			ui.Text(c, est).FontFeatures("tnum").FontWeight(600)
		})
		if res.Submitted() {
			r.order, r.submitted = s.Draft, true
			slog.Debug("mujica: order_entry submitted", "side", s.Draft.Side, "type", s.Draft.Type)
		}
	})
	return r
}

// orderEntryPrice is the expected fill price, or 0 if unknown.
func orderEntryPrice(o OrderDraft, last float64) float64 {
	switch o.Type {
	case OrderLimit, OrderStopLimit:
		return o.Price
	case OrderStop:
		return o.Trigger
	}
	return last
}

// orderEntryFields registers the fields the order type needs.
func orderEntryFields(c *ui.Context, f *input.FormScope, s *OrderEntryState, opts OrderEntryOptions, initial OrderDraft, tick, lot float64, disabled bool) {
	maxQ := opts.MaxQuantity
	input.FormItem(c, f, "quantity", &s.Draft, input.FormItemOptions[OrderDraft]{Label: core.Msg(c, "order_entry.quantity"), Required: true, Initial: initial,
		Validate: func(o OrderDraft) string { return orderEntryCheck(c, o, opts, maxQ) },
		Async:    opts.Check,
	}, func() *ui.Element {
		return input.NumberInput(c, &s.Draft.Quantity, input.NumberInputOptions{Step: lot, Disabled: disabled, Error: s.Form.Error("quantity")}).Input
	})
	if s.Draft.Type == OrderStop || s.Draft.Type == OrderStopLimit {
		input.FormItem(c, f, "trigger", &s.Draft.Trigger, input.FormItemOptions[float64]{Label: core.Msg(c, "order_entry.trigger"), Required: true, Initial: initial.Trigger,
			Validate: func(v float64) string { return orderEntryTrigger(c, s.Draft, opts.Last) },
		}, func() *ui.Element {
			return input.NumberInput(c, &s.Draft.Trigger, input.NumberInputOptions{Step: tick, Disabled: disabled, Error: s.Form.Error("trigger")}).Input
		})
	}
	if s.Draft.Type == OrderLimit || s.Draft.Type == OrderStopLimit {
		input.FormItem(c, f, "price", &s.Draft.Price, input.FormItemOptions[float64]{Label: core.Msg(c, "order_entry.price"), Required: true, Initial: initial.Price,
			Validate: func(v float64) string {
				if v <= 0 {
					return core.Msg(c, "order_entry.price_positive")
				}
				return ""
			},
		}, func() *ui.Element {
			return input.NumberInput(c, &s.Draft.Price, input.NumberInputOptions{Step: tick, Disabled: disabled, Error: s.Form.Error("price")}).Input
		})
	}
	var tifs []input.SelectOption[TimeInForce]
	for t := TIFDay; t <= TIFFillOrKill; t++ {
		tifs = append(tifs, input.SelectOption[TimeInForce]{Value: t, Label: core.Msg(c, fmt.Sprintf("order_entry.tif%d", t))})
	}
	input.FormItem(c, f, "tif", &s.Draft.TIF, input.FormItemOptions[TimeInForce]{Label: core.Msg(c, "order_entry.tif"), Initial: initial.TIF}, func() *ui.Element {
		return input.Select(c, &s.Draft.TIF, tifs, input.SelectOptions{Label: core.Msg(c, "order_entry.tif"), Disabled: disabled}).Element
	})
}

// orderEntryCheck checks the quantity, then the caller's rules.
func orderEntryCheck(c *ui.Context, o OrderDraft, opts OrderEntryOptions, maxQ float64) string {
	switch {
	case o.Quantity <= 0:
		return core.Msg(c, "order_entry.quantity_positive")
	case maxQ > 0 && o.Quantity > maxQ:
		return fmt.Sprintf(core.Msg(c, "order_entry.quantity_max"), data.StatisticFormat(maxQ, 0))
	}
	if opts.Validate != nil {
		return opts.Validate(o)
	}
	return ""
}

// orderEntryTrigger checks a stop's side of the last price.
func orderEntryTrigger(c *ui.Context, o OrderDraft, last float64) string {
	switch {
	case o.Trigger <= 0:
		return core.Msg(c, "order_entry.trigger_positive")
	case last > 0 && o.Side == SideBuy && o.Trigger <= last:
		return core.Msg(c, "order_entry.trigger_above")
	case last > 0 && o.Side == SideSell && o.Trigger >= last:
		return core.Msg(c, "order_entry.trigger_below")
	}
	return ""
}
