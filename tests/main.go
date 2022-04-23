package main

import (
	"fmt"
	"os"
	"time"

	alerthooks "github.com/aldanasjuan/alerthooks-go"
	"github.com/aldanasjuan/errs"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("no env file found")
	}

	alerthooks.SetKey(os.Getenv("KEY"))                   // set your api key
	err = alerthooks.SetSignature(os.Getenv("SIGNATURE")) // your api key has a signing secret, set it here

	record, err := alerthooks.NewRecord(&alerthooks.NewRecordParams{ // create a new record with the url of your own server or function
		Method:   "POST",                     //method you will receive
		Endpoint: os.Getenv("WEBHOOK_URL"),   // your own url
		Type:     alerthooks.RECORD_ONE_TIME, // type of record: one time or recurring
		DueDate:  time.Now().Unix() + 120,    // when it should be set, this example sets it 1:30 minutes after starting server.
		Data: map[string]interface{}{ // data you will receive in the body of the hook
			"testing": "webhook",
			"date":    time.Now().Add(time.Second * 120).Format(time.RFC822),
		},
	})
	errs.Log(err)
	fmt.Printf("%+v\n", record)

	app := fiber.New()

	app.Post("/webhooks", func(ctx *fiber.Ctx) error { // setting an example handler for the hook
		signature := ctx.Get("signature") // the hook comes with a header called 'signature'
		errs.Log(err)
		valid := alerthooks.ValidateSignature(signature)                                                   // use this method to validate it comes from us
		fmt.Printf("[alerthook] signature: %v valid: %v data: %v\n", signature, valid, string(ctx.Body())) // the body will have the info you set up
		return nil
	})

	err = app.Listen(":2121")
	errs.Log(err)
}
