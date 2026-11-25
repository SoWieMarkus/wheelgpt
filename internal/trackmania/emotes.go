package trackmania

// Emote is a 7TV emote the bot uses in its messages.
type Emote struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// The emotes used in chat messages. Streamers should add them to their channel.
var (
	EmoteHeyGuys          = Emote{"HeyGuys", "https://7tv.app/emotes/01F6MZGCNG000255K4X1K7NTHR"}
	EmoteYek              = Emote{"YEK", "https://7tv.app/emotes/01GH7JSW200009M1346QW1QZT6"}
	EmoteChampionMedal    = Emote{"champion_medal", "https://7tv.app/emotes/01GMECGJ5G000DYM1YMXQS1V39"}
	EmoteAuthorMedal      = Emote{"AuthorTime", "https://7tv.app/emotes/01FNS50W1G000EJT2EVEY3F7Y2"}
	EmoteGigachad         = Emote{"GIGACHAD", "https://7tv.app/emotes/01F6MZGCNG000255K4X1K7NTHR"}
	EmoteBwoah            = Emote{"BWOAH", "https://7tv.app/emotes/01FJY0YY4R000E100ZQ63RFZVZ"}
	EmoteOk               = Emote{"ok", "https://7tv.app/emotes/01G1M77D4R0004YN3NKDRR9YKJ"}
	EmoteDinkDonk         = Emote{"dinkDonk", "https://7tv.app/emotes/01EZPJCXQ8000C438200A44F38"}
	EmoteOmegalul         = Emote{"OMEGALUL", "https://7tv.app/emotes/01F00Z3A9G0007E4VV006YKSK9"}
	EmoteICANT            = Emote{"ICANT", "https://7tv.app/emotes/01FA0PBXM8000AD7Y8T9ZMTFZ5"}
	EmotePepePoint        = Emote{"PepePoint", "https://7tv.app/emotes/01FNET3EQ0000FZADBM40VPFGR"}
	EmoteGoldMedal        = Emote{"GLOD", "https://7tv.app/emotes/01G5VMBZT80002C2Y8CGQ69T53"}
	EmoteSilverMedal      = Emote{"SLIVER", "https://7tv.app/emotes/01FZBK3RJ8000DPFPNZDKV8ZG5"}
	EmoteBronzeMedal      = Emote{"BORNZE", "https://7tv.app/emotes/01G0QABGS8000005X345PAXJBK"}
	EmoteReallyFuckingMad = Emote{"ReallyFuckingMad", "https://7tv.app/emotes/01FQ1GTQJ80005G17242N6NC7E"}
)

// Emotes lists all emotes in a stable order.
var Emotes = []Emote{
	EmoteHeyGuys, EmoteYek, EmoteChampionMedal, EmoteAuthorMedal, EmoteGigachad,
	EmoteBwoah, EmoteOk, EmoteDinkDonk, EmoteOmegalul, EmoteICANT, EmotePepePoint,
	EmoteGoldMedal, EmoteSilverMedal, EmoteBronzeMedal, EmoteReallyFuckingMad,
}
