package bot

// Emote holds the chat name of a 7TV/Twitch emote.
type Emote struct{ Name string }

var (
	EmoteHeyGuys           = Emote{"HeyGuys"}
	EmoteYek               = Emote{"YEK"}
	EmoteChampionMedal     = Emote{"champion_medal"}
	EmoteAuthorMedal       = Emote{"AuthorTime"}
	EmoteGigachad          = Emote{"GIGACHAD"}
	EmoteBwoah             = Emote{"BWOAH"}
	EmoteOK                = Emote{"ok"}
	EmoteDinkDonk          = Emote{"dinkDonk"}
	EmoteOmegalul          = Emote{"OMEGALUL"}
	EmoteICant             = Emote{"ICANT"}
	EmotePepePoint         = Emote{"PepePoint"}
	EmoteGoldMedal         = Emote{"GLOD"}
	EmoteSilverMedal       = Emote{"SLIVER"}
	EmoteBronzeMedal       = Emote{"BORNZE"}
	EmoteReallyFuckingMad  = Emote{"ReallyFuckingMad"}
)

func AllEmotes() []Emote {
	return []Emote{
		EmoteHeyGuys, EmoteYek, EmoteChampionMedal, EmoteAuthorMedal,
		EmoteGigachad, EmoteBwoah, EmoteOK, EmoteDinkDonk,
		EmoteOmegalul, EmoteICant, EmotePepePoint,
		EmoteGoldMedal, EmoteSilverMedal, EmoteBronzeMedal, EmoteReallyFuckingMad,
	}
}
