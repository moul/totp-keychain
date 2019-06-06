package main // import "moul.io/totp-keychain"

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	keychain "github.com/keybase/go-keychain"
	"github.com/pquerna/otp/totp"
	"github.com/urfave/cli"
)

const ServiceName = "totp-keychain"

func main() {
	app := cli.NewApp()
	app.Name = "totp-keychain"
	app.Commands = []cli.Command{
		{
			Name: "add",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Usage: "Name of the account",
				},
				cli.StringFlag{
					Name:  "secret, s",
					Usage: "TOTP Secret",
				},
			},
			Action: func(c *cli.Context) error {
				if c.NArg() < 2 {
					return errors.New("usage: add '<name>' '<secret>'")
				}
				name := c.Args()[0]
				secret := c.Args()[1]
				item := keychain.NewItem()
				item.SetSecClass(keychain.SecClassGenericPassword)
				item.SetService(ServiceName)
				item.SetAccount(name)
				item.SetData([]byte(secret))
				//item.SetSynchronizable(keychain.SynchronizableYes)
				item.SetAccessible(keychain.AccessibleWhenUnlocked)
				err := keychain.AddItem(item)
				if err == keychain.ErrorDuplicateItem {
					log.Println("Item already exists. To update it, remove it first and add it back.")
					return nil
				}
				log.Println("OK")
				return err
			},
		},
		{
			Name: "rm",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					return errors.New("usage: rm <name>")
				}
				account := c.Args().First()
				item := keychain.NewItem()
				item.SetSecClass(keychain.SecClassGenericPassword)
				item.SetService(ServiceName)
				item.SetAccount(account)
				if err := keychain.DeleteItem(item); err != nil {
					return err
				}
				log.Println("OK")
				return nil
			},
		},
		{
			Name: "ls",
			Action: func(c *cli.Context) error {
				item := keychain.NewItem()
				item.SetSecClass(keychain.SecClassGenericPassword)
				item.SetService(ServiceName)
				item.SetMatchLimit(keychain.MatchLimitAll)
				item.SetReturnAttributes(true)
				results, err := keychain.QueryItem(item)
				if err != nil {
					return err
				}
				for _, r := range results {
					fmt.Printf("- %s\n", r.Account)
				}
				return nil
			},
		},
		{
			Name: "get",
			Action: func(c *cli.Context) error {
				if c.NArg() < 1 {
					return errors.New("usage: get <name>")
				}
				account := c.Args().First()
				query := keychain.NewItem()
				query.SetSecClass(keychain.SecClassGenericPassword)
				query.SetService(ServiceName)
				query.SetAccount(account)
				query.SetMatchLimit(keychain.MatchLimitOne)
				query.SetReturnData(true)
				results, err := keychain.QueryItem(query)
				if err != nil {
					return err
				}
				if len(results) < 1 {
					return errors.New("No such account with that name")
				}
				code, err := totp.GenerateCode(string(results[0].Data), time.Now())
				if err != nil {
					return err
				}
				fmt.Println(code)
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
