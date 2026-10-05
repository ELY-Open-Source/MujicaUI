package mujica

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestReadReceiptStates(t *testing.T) {
	readers := []AvatarItem{{Name: "Uika"}, {Name: "Mutsumi"}, {Name: "Umiri"}, {Name: "Nyamu"}}
	tt := ui.NewTester(func(c *ui.Context) {
		Use(c, Settings{})
		ui.Column(c).Gap(8).AlignItems(ui.Start).Children(func() {
			ReadReceipt(c, ReceiptSent, ReadReceiptOptions{})
			ReadReceipt(c, ReceiptDelivered, ReadReceiptOptions{})
			ReadReceipt(c, ReceiptRead, ReadReceiptOptions{Readers: readers, Max: 2})
		})
	}, 400, 160)
	for _, s := range []string{"Sent", "Delivered", "+2"} {
		if !tt.HasText(s) {
			t.Fatalf("no %q in %v", s, tt.Texts())
		}
	}
	if _, ok := tt.Find("Read by Uika, Mutsumi, Umiri, Nyamu"); !ok {
		t.Fatal("readers not named in full")
	}
	mediaPanics(t, "bad state", func(c *ui.Context) { ReadReceipt(c, ReceiptState(5), ReadReceiptOptions{}) })
}

func TestReadReceiptSnapshot(t *testing.T) {
	snapshot(t, "read_receipt", 360, 130, func(c *ui.Context) {
		ui.Column(c).Padding(20).Gap(10).AlignItems(ui.Start).Children(func() {
			ReadReceipt(c, ReceiptSent, ReadReceiptOptions{})
			ReadReceipt(c, ReceiptDelivered, ReadReceiptOptions{})
			ReadReceipt(c, ReceiptRead, ReadReceiptOptions{Readers: []AvatarItem{{Name: "Uika Misumi"}, {Name: "Mutsumi Wakaba"}, {Name: "Umiri Yahata"}, {Name: "Nyamu"}}})
		})
	})
}
