package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

type SheetClient struct {
	srv           *sheets.Service
	spreadsheetID string
}

func NewSheetClient(ctx context.Context, spreadsheetID string) (*SheetClient, error) {
	b, err := ioutil.ReadFile("secret.json")
	if err != nil {
		return nil, err
	}
	// read & write permission
	jwt, err := google.JWTConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}
	srv, err := sheets.New(jwt.Client(ctx))
	if err != nil {
		return nil, err
	}
	return &SheetClient{
		srv:           srv,
		spreadsheetID: spreadsheetID,
	}, nil
}

func (s *SheetClient) BatchUpdate() error {
	sheetID, err := strconv.Atoi(os.Getenv("SHEET_ID"))
	if err != nil {
		return err
	}
	req := &sheets.Request{
		SetDataValidation: &sheets.SetDataValidationRequest{
			Rule: &sheets.DataValidationRule{
				Condition: &sheets.BooleanCondition{
					Type: "BOOLEAN",
				},
			},
			Range: &sheets.GridRange{
				StartRowIndex:    int64(1),
				EndRowIndex:      int64(3),
				StartColumnIndex: int64(1),
				EndColumnIndex:   int64(3),
				SheetId:          int64(sheetID),
			},
		},
	}
	bus := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{req},
	}
	_, err = s.srv.Spreadsheets.BatchUpdate(s.spreadsheetID, bus).Do()
	return err
}

func main() {
	client, err := NewSheetClient(context.Background(), os.Getenv("SPREAD_SHEET_ID"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.BatchUpdate()
	if err != nil {
		log.Fatalf("client.BatchUpdate: %s", err)
	}
}
