package model

import (
	"net/http"
	"testing"
)

func TestIpAddrFromRemoteAddr(t *testing.T) {
	// Test case 1: Verify that the ipAddrFromRemoteAddr() function returns the expected IP address.
	remoteAddr := "192.168.1.1:8080"
	expectedIP := "192.168.1.1"
	result := ipAddrFromRemoteAddr(remoteAddr)
	if result != expectedIP {
		t.Errorf("Test case 1: expected ipAddrFromRemoteAddr(%q) to return %q, but got %q", remoteAddr, expectedIP, result)
	}

}

func TestRequestGetRemoteAddress(t *testing.T) {
	// Test case 1: Verify that the requestGetRemoteAddress() function returns the expected remote address with X-Forwarded-For header.
	req1 := &http.Request{
		RemoteAddr: "192.168.1.1:8080",
		Header: http.Header{
			"X-Real-Ip":       []string{"192.168.1.2"},
			"X-Forwarded-For": []string{"192.168.1.3, 192.168.1.4"},
		},
	}
	expectedIP1 := "192.168.1.3"
	result1 := requestGetRemoteAddress(req1)
	if result1 != expectedIP1 {
		t.Errorf("Test case 1: expected requestGetRemoteAddress() to return %q, but got %q", expectedIP1, result1)
	}

	// Test case 2: Verify that the requestGetRemoteAddress() function returns the expected remote address with X-Real-Ip header.
	req2 := &http.Request{
		RemoteAddr: "192.168.1.1:8080",
		Header: http.Header{
			"X-Real-Ip": []string{"192.168.1.2"},
		},
	}
	expectedIP2 := "192.168.1.2"
	result2 := requestGetRemoteAddress(req2)
	if result2 != expectedIP2 {
		t.Errorf("Test case 2: expected requestGetRemoteAddress() to return %q, but got %q", expectedIP2, result2)
	}

	// Test case 3: Verify that the requestGetRemoteAddress() function returns the expected remote address with RemoteAddr.
	req3 := &http.Request{
		RemoteAddr: "192.168.1.1:8080",
		Header:     http.Header{},
	}
	expectedIP3 := "192.168.1.1"
	result3 := requestGetRemoteAddress(req3)
	if result3 != expectedIP3 {
		t.Errorf("Test case 3: expected requestGetRemoteAddress() to return %q, but got %q", expectedIP3, result3)
	}

	// Test case 4: Verify that the requestGetRemoteAddress() function returns an empty string with an empty request.
	req4 := &http.Request{}
	expectedIP4 := ""
	result4 := requestGetRemoteAddress(req4)
	if result4 != expectedIP4 {
		t.Errorf("Test case 4: expected requestGetRemoteAddress() to return %q, but got %q", expectedIP4, result4)
	}

}
