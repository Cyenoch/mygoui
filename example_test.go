package mygoui_test

import (
	"fmt"

	"github.com/egoist/mygo/ui"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
)

func ExampleButton() {
	count := 0
	tester := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		mygoui.Button(c, mygoui.ButtonOptions{
			Label: "Increment", OnClick: func() { count++ },
		})
	}, 240, 80)
	if err := tester.Click("Increment"); err != nil {
		panic(err)
	}
	fmt.Println(count)
	// Output: 1
}
