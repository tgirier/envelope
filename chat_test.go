package chat_test

import "testing"

func TestServerConn(t *testing.T) {
	addr := ":8080"
	errExpected := false

	conn, err := chat.NewConn(addr)
	defer conn.Close()

	if errExpected != (err != nil) {
		t.Errorf("NewConn(%s): error expected %t, got %v", addr, errExpected, err)
	}
}

func TestWelcomeMessage(t *testing.T) {
	addr := ":8080"
	want := "Welcome to GoChat !"
	errExpected := false

	conn, _ := chat.NewConnection(addr)
	defer conn.Close()

	got, err := conn.Read()

	if errExpected != (err != nil) {
		t.Errorf("conn.Read(): error expected %t, got %v", errExpected, err)
	}
	if got != want {
		t.Errorf("conn.Read() welcome message: got %s, want %s", got, want)
	}
}
