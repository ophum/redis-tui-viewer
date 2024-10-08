package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rivo/tview"
)

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer rdb.Close()
	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		panic(err)
		return
	}

	f, _ := os.Create("./log")
	defer f.Close()
	log.SetOutput(f)
	app := tview.NewApplication()

	list := tview.NewList()
	list.SetBorder(true)
	list.SetTitle("keys")

	textView := tview.NewTextView()
	textView.SetBorder(true)
	textView.SetTitle("value")

	table := tview.NewTable().SetBorders(true)
	table.SetBorder(true)
	table.SetTitle("keys")
	for i, key := range keys {
		table.SetCell(i, 0, tview.NewTableCell(key))
		list = list.AddItem(key, "", 'a', nil)
	}
	table.SetSelectable(true, true)
	table.Select(0, 0).SetFixed(1, 1).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
		if key == tcell.KeyEnter {
			log.Println("enter")
		}
	}).SetSelectedFunc(func(row, column int) {
		log.Println("selected key:", keys[row])
		//table.SetSelectable(false, false)

		key := keys[row]
		res, err := rdb.Get(ctx, key).Result()
		if err != nil {
			panic(err)
		}

		var m map[string]any
		if err := json.Unmarshal([]byte(res), &m); err != nil {
			textView.SetText(res)
		} else {
			text, err := json.MarshalIndent(m, "", "  ")
			if err != nil {
				panic(err)
			}

			textView.SetText(string(text))
		}

	})

	flex := tview.NewFlex().
		AddItem(table, 0, 1, false).
		AddItem(textView, 0, 1, false)

	page := tview.NewPages()
	page.AddPage("redis tui viewer", flex, true, true)
	page.SetBorder(true).SetTitle("redis tui viewer")

	if err := app.SetRoot(page, true).SetFocus(table).Run(); err != nil {
		panic(err)
	}

}
