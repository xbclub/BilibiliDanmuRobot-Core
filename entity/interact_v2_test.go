package entity

import "testing"

func TestParseInteractWordV2(t *testing.T) {
	pb := "CJ3N6wwSCW1pa3VmaWxjayICAwEoAzDwxdkOOLL/484GQKHWmMPXM0oxCOeVgKuwraYGEBUaCeWBmueMq+eahCDLqGkoy6hpMJK7ygI4y6hpQAFg8MXZDmjsE2IAeKbf5aPU38DSGJoBALIBrgIInc3rDBK9AQoJbWlrdWZpbGNrEkpodHRwczovL2kyLmhkc2xiLmNvbS9iZnMvZmFjZS82YTBkZDM2YmE3ZmExOGU4NGM2NTI3Yzg3YmViYTAyYTVkMmRlM2Y5LmpwZzJXCgltaWt1ZmlsY2sSSmh0dHBzOi8vaTIuaGRzbGIuY29tL2Jmcy9mYWNlLzZhMGRkMzZiYTdmYTE4ZTg0YzY1MjdjODdiZWJhMDJhNWQyZGUzZjkuanBnOgsg////////////ARplCgnlgZrnjKvnmoQQFRjLqGkgkrvKAijLqGkwy6hpSAFQ55WAq7CtpgZg7BN6CSMzRkI0RjY5OYIBCSMzRkI0RjY5OYoBCSMzRkI0RjY5OZIBByNGRkZGRkaaAQkjM0ZCNEY2RTYyALoBAA=="
	got, err := ParseInteractWordV2(pb)
	if err != nil {
		t.Fatalf("ParseInteractWordV2() error = %v", err)
	}
	if got.UID != 26928797 || got.Uname != "mikufilck" || got.MsgType != 3 {
		t.Fatalf("unexpected identity: %+v", got)
	}
	if got.FansMedal == nil {
		t.Fatal("FansMedal is nil")
	}
	m := got.FansMedal
	if m.TargetID != 3546384651258599 || m.MedalLevel != 21 || m.MedalName != "做猫的" || m.IsLighted != 1 || m.AnchorRoomid != 30827248 || m.GuardLevel != 0 {
		t.Fatalf("unexpected medal: %+v", m)
	}
}

func TestParseInteractWordV2Invalid(t *testing.T) {
	if _, err := ParseInteractWordV2("!!!"); err == nil {
		t.Fatal("expected an error for invalid base64")
	}
}
