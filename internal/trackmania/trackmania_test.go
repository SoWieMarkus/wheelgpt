package trackmania

import (
	"strings"
	"testing"
)

func TestParseTime(t *testing.T) {
	tests := []struct {
		in   string
		want Time
		ok   bool
	}{
		{"12.345", 12_345, true},
		{"1.5", 1_500, true},
		{"1.05", 1_050, true},
		{"0.001", 1, true},
		{"1:02.345", 62_345, true},
		{"01:23.456", 83_456, true},
		{"1:01:02.345", 3_662_345, true},
		{"12", 0, false}, // milliseconds are required
		{"12.3456", 0, false},
		{"abc", 0, false},
		{"1:2:3:4.5", 0, false},
		{"-1.000", 0, false},
		{"", 0, false},
	}
	for _, tt := range tests {
		got, ok := ParseTime(tt.in)
		if ok != tt.ok || got != tt.want {
			t.Errorf("ParseTime(%q) = %d, %v; want %d, %v", tt.in, got, ok, tt.want, tt.ok)
		}
	}
}

func TestTimeString(t *testing.T) {
	tests := map[Time]string{
		0:         "0.000",
		1:         "0.001",
		12_345:    "12.345",
		62_345:    "1:02.345",
		3_662_345: "1:01:02.345",
		-1_500:    "1.500",
	}
	for in, want := range tests {
		if got := in.String(); got != want {
			t.Errorf("Time(%d).String() = %q; want %q", in, got, want)
		}
	}
}

var testMap = Map{
	UID: "abc", Name: "$o$fffSummer $z05", Author: "Nadeo",
	ChampionTime: 40_000, AuthorTime: 41_000, GoldTime: 44_000, SilverTime: 50_000, BronzeTime: 60_000,
}

func TestMedal(t *testing.T) {
	tests := map[Time]Medal{
		39_000: MedalChampion, 40_000: MedalChampion, 40_001: MedalAuthor,
		44_000: MedalGold, 49_999: MedalSilver, 60_000: MedalBronze, 60_001: MedalNone,
	}
	for in, want := range tests {
		if got := testMap.Medal(in); got != want {
			t.Errorf("Medal(%d) = %v; want %v", in, got, want)
		}
	}
	noChampion := testMap
	noChampion.ChampionTime = 0
	if got := noChampion.Medal(1); got != MedalAuthor {
		t.Errorf("without champion time: got %v", got)
	}
}

func TestStripStyling(t *testing.T) {
	if got := StripStyling(testMap.Name); got != "Summer 05" {
		t.Errorf("StripStyling = %q", got)
	}
	if got := StripStyling("$$money"); got != "money" {
		t.Errorf("StripStyling($$) = %q", got)
	}
}

func TestMapString(t *testing.T) {
	m := testMap
	m.TMXID = 42
	want := `"Summer 05" by Nadeo | champion_medal  40.000 | AuthorTime  41.000 | https://trackmania.exchange/maps/42 | https://trackmania.io/#/leaderboard/abc`
	if got := m.String(); got != want {
		t.Errorf("String() =\n%s\nwant\n%s", got, want)
	}
	m.TMXID, m.ChampionTime = 0, 0
	want = `"Summer 05" by Nadeo | AuthorTime  41.000 | https://trackmania.io/#/leaderboard/abc`
	if got := m.String(); got != want {
		t.Errorf("String() =\n%s\nwant\n%s", got, want)
	}
}

func TestEvaluate(t *testing.T) {
	guesses := []Guess{
		{UserID: "1", DisplayName: "a", Time: 10_000},
		{UserID: "2", DisplayName: "b", Time: 10_200},
		{UserID: "3", DisplayName: "c", Time: 9_800},
	}
	winners := Evaluate(guesses, 10_100)
	if len(winners) != 2 || winners[0].UserID != "1" || winners[1].UserID != "2" {
		t.Fatalf("tie: got %+v", winners)
	}
	if p, perfect := PointsFor(winners, 10_100); p != PointsClosestGuess || perfect {
		t.Errorf("points = %d, %v", p, perfect)
	}
	winners = Evaluate(guesses, 9_800)
	if len(winners) != 1 || winners[0].UserID != "3" {
		t.Fatalf("perfect: got %+v", winners)
	}
	if p, perfect := PointsFor(winners, 9_800); p != PointsPerfectGuess || !perfect {
		t.Errorf("points = %d, %v", p, perfect)
	}
	if len(Evaluate(nil, 1)) != 0 {
		t.Error("no guesses must have no winners")
	}
}

func TestMessages(t *testing.T) {
	m := testMap
	tests := []struct {
		pb   Time
		want string
	}{
		{39_500, "That's Champion Medal champion_medal BWOAH"},
		{40_005, "HAHAHAHAH HE MISSED CHAMPION BY 0.005s"},
		{40_500, "That's AuthorTime ok"},
		{41_003, "HAHAHAHAH HE MISSED AUTHOR BY 0.003s"},
		{43_000, "That's only GLOD"},
		{70_000, "Not even BORNZE PepePoint"},
	}
	for _, tt := range tests {
		if got := PBMessage(&m, tt.pb); !strings.Contains(got, tt.want) {
			t.Errorf("PBMessage(%d) = %q; want it to contain %q", tt.pb, got, tt.want)
		}
	}
	if got := PBMessage(nil, 1); !strings.Contains(got, "didn't know you are on a map") {
		t.Errorf("no map: %q", got)
	}

	if got := WinnerMessage(nil, 1); got != "But no chatter participated ReallyFuckingMad" {
		t.Errorf("no winners: %q", got)
	}
	one := []Guess{{DisplayName: "Viewer", Time: 10_000}}
	if got := WinnerMessage(one, 10_000); got != "@Viewer the GIGACHAD guessed it correctly! BWOAH (+5 points)" {
		t.Errorf("perfect: %q", got)
	}
	if got := WinnerMessage(one, 10_250); got != "Nobody guessed it correctly but @Viewer guessed 10.000 (-0.250) ok (+1 point)" {
		t.Errorf("closest: %q", got)
	}
	two := []Guess{{DisplayName: "A", Time: 10_000}, {DisplayName: "B", Time: 10_000}}
	if got := WinnerMessage(two, 9_900); got != "Nobody guessed it correctly but @A, @B guessed 10.000 (+0.100) ok (+1 point)" {
		t.Errorf("tie: %q", got)
	}
}
