package main

import (
	"context"
	_ "embed"
	"log"
	log2 "github.com/labstack/gommon/log"
	// "net/http"
	// "time"

	"ajaxbits.com/bsplit/internal/db"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed schema.sql
var ddl string

var ctx = context.Background()
var readDb, writeDb = db.Init()

func main() {
	defer readDb.Close()
	defer writeDb.Close()

	if _, err := writeDb.ExecContext(ctx, ddl); err != nil {
		log.Fatal(err)
	}

	aliceUuid, err := uuid.NewV7()
	if err != nil {
		log.Fatal(err)
	}
	bobUuid, err := uuid.NewV7()
	if err != nil {
		log.Fatal(err)
	}

	alice, err := writeQueries.CreateUser(ctx, db.CreateUserParams{
		Uuid:    aliceUuid.String(),
		Name:    "Alice",
		VenmoID: nil,
	})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(alice)
	}

	bob, err := writeQueries.CreateUser(ctx, db.CreateUserParams{
		Uuid:    bobUuid.String(),
		Name:    "Bob",
		VenmoID: nil,
	})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(bob)
	}

	groupUuid, err := uuid.NewV7()
	if err != nil {
		log.Fatal(err)
	}
	groupDescription := "Alice and Bob are in love!"

	group, err := writeQueries.CreateGroup(ctx, db.CreateGroupParams{
		Uuid:        groupUuid.String(),
		Name:        "Lovers",
		Description: &groupDescription,
	})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(group)
	}

	// txUuid1, err := uuid.NewV7()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// txUuid2, err := uuid.NewV7()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// tpUuid1, err := uuid.NewV7()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// tpUuid2, err := uuid.NewV7()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// tpUuid3, err := uuid.NewV7()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// tpUuid4, err := uuid.NewV7()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// tx1, err := writeQueries.CreateTransaction(ctx, db.CreateTransactionParams{
	// 	Uuid:        txUuid1.String(),
	// 	Type:        "expense",
	// 	Description: "dinner",
	// 	Date:        time.Now().UTC().Unix(),
	// 	PaidBy:      alice.Uuid,
	// 	GroupUuid:   nil,
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if _, err = writeQueries.CreateTransactionParticipants(ctx, db.CreateTransactionParticipantsParams{
	// 	Uuid:     tpUuid1.String(),
	// 	TxnUuid:  tx1.Uuid,
	// 	UserUuid: alice.Uuid,
	// 	Share:    500,
	// }); err != nil {
	// 	log.Fatal(err)
	// }
	// if _, err = writeQueries.CreateTransactionParticipants(ctx, db.CreateTransactionParticipantsParams{
	// 	Uuid:     tpUuid2.String(),
	// 	TxnUuid:  tx1.Uuid,
	// 	UserUuid: bob.Uuid,
	// 	Share:    500,
	// }); err != nil {
	// 	log.Fatal(err)
	// }

	// tx2, err := writeQueries.CreateTransaction(ctx, db.CreateTransactionParams{
	// 	Uuid:        txUuid2.String(),
	// 	Type:        "expense",
	// 	Description: "new car",
	// 	Amount:      5000,
	// 	Date:        time.Now().UTC().Unix(),
	// 	PaidBy:      bob.Uuid,
	// 	GroupUuid:   nil,
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if _, err = writeQueries.CreateTransactionParticipants(ctx, db.CreateTransactionParticipantsParams{
	// 	Uuid:     tpUuid3.String(),
	// 	TxnUuid:  tx2.Uuid,
	// 	UserUuid: alice.Uuid,
	// 	Share:    2000,
	// }); err != nil {
	// 	log.Fatal(err)
	// }
	// if _, err = writeQueries.CreateTransactionParticipants(ctx, db.CreateTransactionParticipantsParams{
	// 	Uuid:     tpUuid4.String(),
	// 	TxnUuid:  tx2.Uuid,
	// 	UserUuid: bob.Uuid,
	// 	Share:    2000,
	// }); err != nil {
	// 	log.Fatal(err)
	// }

	// debts, err := readQueries.GetDebts(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// for _, debt := range debts {
	// 	debtor, err := readQueries.GetUser(ctx, debt.Debtor)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	creditor, err := readQueries.GetUser(ctx, debt.Creditor)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	log.Default().Printf("%s owes %s: %v", debtor.Name, creditor.Name, debt.NetAmount)
	// }
	
	e := echo.New()
	e.Logger.SetLevel(log2.INFO)
	e.Use(middleware.Logger())
	e.GET("/user", UserHandler)
	e.Logger.Fatal(e.Start(":8080"))

	// http.HandleFunc("/", handlers.RootHandler)
	// http.HandleFunc("/split", handlers.SplitHandler)
	// http.HandleFunc("/user", UserHandler)
	// http.HandleFunc("/group", GroupHandler)
	// http.HandleFunc("/txn", TransactionHandler)
	// log.Println("Starting bsplit server on :8080")
	// log.Fatal(http.ListenAndServe(":8080", nil))
}
