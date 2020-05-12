package chat_test

func TestServerConn(t *testing.T) {
	addr := ":8080"
	errExpected := false

	conn, err := chat.NewConn(addr)
	defer conn.Close()

	if errExpected != (err != nil) {
		t.Errorf("NewConn(%s): error expected %t, got %v", addr, errExpected, err)
	}
}


