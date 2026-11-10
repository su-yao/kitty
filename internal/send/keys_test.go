package send

import "testing"

func TestCanTypeASCII(t *testing.T) {
	ok := []string{
		"cash 5000",
		"add_ship NAME_Battleship",
		"effect spawn_megastructure = { type = dyson_sphere_5 }",
		"resource sr_entropy 5000",
	}
	for _, s := range ok {
		if !canTypeASCII(s) {
			t.Fatalf("want typeable %q", s)
		}
	}
	bad := []string{"中文", "add_ship 女王", "§"}
	for _, s := range bad {
		if canTypeASCII(s) {
			t.Fatalf("want clipboard %q", s)
		}
	}
}

func TestAltNumpadVKs(t *testing.T) {
	got := altNumpadVKs()
	if len(got) != 3 || got[0] != VK_MENU || got[1] != VK_NUMPAD2 || got[2] != VK_NUMPAD1 {
		t.Fatalf("alt numpad VKs: %v", got)
	}
	if VK_NUMPAD1 == 0x4F || VK_NUMPAD2 == 0x50 {
		t.Fatal("numpad VKs must not be Down/End scancodes")
	}
}
