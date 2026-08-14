package controllertest

import (
	"os"
	"restaurant-order-management/config"
	"testing"
)

func TestMain(m *testing.M) {
	config.ConnectDB("../../.env")
	os.Exit(m.Run())
}
