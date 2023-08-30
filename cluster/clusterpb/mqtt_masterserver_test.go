package clusterpb

import "testing"

func Test_NewMqttMasterServer(t *testing.T) {

	var (
		advertiseAddr = ":8081"
	)

	server := NewMqttMasterServer()

	err := server.Listener(advertiseAddr)
	if err != nil {
		t.Fatal(err)
	}

	select {}
}
