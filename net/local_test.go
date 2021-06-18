package net

import "testing"

func TestGetLocalIp(t *testing.T) {
	localIp := "192.168.20.98"

	ip, err := GetLocalIp()
	if err != nil {
		t.Error(err)
	}

	t.Log(localIp == ip)
}
