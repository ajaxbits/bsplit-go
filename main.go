package main

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	_ "embed"
	"log"

	// "net/http"
	"time"

	"ajaxbits.com/bsplit/internal/dbc"
	"github.com/google/uuid"

	// "ajaxbits.com/bsplit/internal/handlers"
)

//go:embed schema.sql
var ddl string

func main() {
	database, err := sql.Open("sqlite3", "./expenses.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	ctx := context.Background()
	queries := dbc.New(database)

	if _, err := database.ExecContext(ctx, ddl); err != nil {
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

	alice, err := queries.CreateUser(ctx, dbc.CreateUserParams{
		Uuid:    aliceUuid.String(),
		Name:    "Alice",
		VenmoID: nil,
	})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(alice)
	}

	bob, err := queries.CreateUser(ctx, dbc.CreateUserParams{
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

	group, err := queries.CreateGroup(ctx, dbc.CreateGroupParams{
		Uuid:        groupUuid.String(),
		Name:        "Lovers",
		Description: &groupDescription,
	})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println(group)
	}

	txUuid1, err := uuid.NewV7()
	if err != nil {
		log.Fatal(err)
	}
	txUuid2, err := uuid.NewV7()
	if err != nil {
		log.Fatal(err)
	}
    tpUuid1, err := uuid.NewV7()
	if err != nil {
		log.Fatal(err)
	}
	tpUuid2, err := uuid.NewV7()
	if err != nil {
		log.Fatal(err)
	}
	tpUuid3, err := uuid.NewV7()
	if err != nil {
		log.Fatal(err)
	}
	tpUuid4, err := uuid.NewV7()
	if err != nil {
		log.Fatal(err)
	}

	tx1, err := queries.CreateTransaction(ctx, dbc.CreateTransactionParams{
		Uuid:        txUuid1.String(),
		Type:        "expense",
		Description: "dinner",
		Amount:      1000,
		Date:        time.Now().UTC(),
		PaidBy:      alice.Uuid,
		GroupUuid:   nil,
	})
	if err != nil {
		log.Fatal(err)
	}


    if _, err = queries.CreateTransactionParticipants(ctx, dbc.CreateTransactionParticipantsParams{
        Uuid: tpUuid1.String(),
        TxnUuid: tx1.Uuid,
        UserUuid: alice.Uuid,
        Share: 500,
    }); err != nil {
        log.Fatal(err)
    }
    if _, err = queries.CreateTransactionParticipants(ctx, dbc.CreateTransactionParticipantsParams{
        Uuid: tpUuid2.String(),
        TxnUuid: tx1.Uuid,
        UserUuid: bob.Uuid,
        Share: 500,
    }); err != nil {
        log.Fatal(err)
    }

	tx2, err := queries.CreateTransaction(ctx, dbc.CreateTransactionParams{
		Uuid:        txUuid2.String(),
		Type:        "expense",
		Description: "new car",
		Amount:      5000,
		Date:        time.Now().UTC(),
		PaidBy:      bob.Uuid,
		GroupUuid:   nil,
	})
	if err != nil {
		log.Fatal(err)
	}

    if _, err = queries.CreateTransactionParticipants(ctx, dbc.CreateTransactionParticipantsParams{
        Uuid: tpUuid3.String(),
        TxnUuid: tx2.Uuid,
        UserUuid: alice.Uuid,
        Share: 2000,
    }); err != nil {
        log.Fatal(err)
    }
    if _, err = queries.CreateTransactionParticipants(ctx, dbc.CreateTransactionParticipantsParams{
        Uuid: tpUuid4.String(),
        TxnUuid: tx2.Uuid,
        UserUuid: bob.Uuid,
        Share: 2000,
    }); err != nil {
        log.Fatal(err)
    }


	debts, err := queries.GetDebts(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, debt := range debts {
		debtor, err := queries.GetUser(ctx, debt.Debtor)
		if err != nil {
			log.Fatal(err)
		}
		creditor, err := queries.GetUser(ctx, debt.Creditor)
		if err != nil {
			log.Fatal(err)
		}

		log.Default().Printf("%s owes %s: %v", debtor.Name, creditor.Name, debt.NetAmount)
	}

	// http.HandleFunc("/", handlers.RootHandler)
	// http.HandleFunc("/split", handlers.SplitHandler)
	// log.Println("Starting bsplit server on :8080")
	// log.Fatal(http.ListenAndServe(":8080", nil))
}
