package handlertest

import (
	"os"
	"restaurant-order-management/config"
	"testing"
)

// func TestServer()

func TestMain(m *testing.M) {
	config.ConnectDB("../../.env")
	os.Exit(m.Run())
}
