package node

import "testing"

type testHostIp struct {
	begin string
	no int
	excepted string
}

func TestGenerateHostIp(t *testing.T) {
	generateNodeIp("127.0.0.1", 0)
	tests := []testHostIp{
		{
			"127.0.0.1",
			0,
			"127.0.0.1",
		},
		{
			"127.0.0.1",
			1,
			"127.0.0.2",
		},
		{
			"127.0.0.1",
			256,
			"127.0.1.1",
		},
		{
			"127.0.0.1",
			257,
			"127.0.1.2",
		},
		{
			"127.0.0.1",
			256*256,
			"127.1.0.1",
		},
		{
			"a7.0.0.1",
			0,
			"",
		},
	}

	for _,test := range tests {
		ip := generateNodeIp(test.begin, test.no)
		if ip != test.excepted {
			t.Fatalf("test:%+v,real:%s", test, ip)
		}
	}
}