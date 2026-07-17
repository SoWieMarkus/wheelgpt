package emote

type Emote string

const (
	HeyGuys          Emote = "HeyGuys"
	YEK              Emote = "YEK"
	ChampionMedal    Emote = "champion_medal"
	AuthorMedal      Emote = "AuthorTime"
	GoldMedal        Emote = "GLOD"
	SilverMedal      Emote = "SLIVER"
	BronzeMedal      Emote = "BORNZE"
	GigaChad         Emote = "GIGACHAD"
	BWOAH            Emote = "BWOAH"
	Ok               Emote = "ok"
	DinkDonk         Emote = "dinkDonk"
	Omegalul         Emote = "OMEGALUL"
	ICant            Emote = "ICANT"
	PepePoint        Emote = "PepePoint"
	ReallyFuckingMad Emote = "ReallyFuckingMad"
)

func GetAllEmotes() []Emote {
	return []Emote{
		HeyGuys,
		YEK,
		ChampionMedal,
		AuthorMedal,
		GoldMedal,
		SilverMedal,
		BronzeMedal,
		GigaChad,
		BWOAH,
		Ok,
		DinkDonk,
		Omegalul,
		ICant,
		PepePoint,
		ReallyFuckingMad,
	}
}
