package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// DisableWiFi : turns WiFi on or off
func DisableWiFi(disable bool) {

	var buttonQuerySelector string
	if disable {
		buttonQuerySelector = ".enable-disable-on#tf1_wifiEnable2point4Gh"
	} else {
		buttonQuerySelector = ".enable-disable-off#tf1_wifiEnable2point4Gh"
	}

	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)

	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 3500*time.Second)
	defer cancel()

	res := ""

	err := chromedp.Run(ctx,

		chromedp.ActionFunc(func(ctx context.Context) error {

			chromedp.ListenTarget(ctx, func(ev interface{}) {

				if _, ok := ev.(*page.EventJavascriptDialogOpening); ok { // page loaded

					t := page.HandleJavaScriptDialog(true)
					go func() {
						if err := chromedp.Run(ctx, t); err != nil {
							fmt.Println(err)
						}

						fmt.Printf(ev.(*page.EventJavascriptDialogOpening).Message) // holds msg!

					}()
				}
			})
			return nil
		}),

		chromedp.Navigate(os.Getenv("DWR_URL")),

		chromedp.SendKeys(
			"#tf1_usrName",
			os.Getenv("DWR_USERNAME"),
		),
		chromedp.SendKeys(
			"#tf1_password",
			os.Getenv("DWR_PASSWORD"),
		),

		chromedp.Click("#btSave"), // login button

		chromedp.Click(".wifiIconn"),

		chromedp.Click(buttonQuerySelector),

		chromedp.EvaluateAsDevTools(`document.querySelector('.btnSubmit').click()`, &res),
	)

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[+] Done! completed the script")
}
