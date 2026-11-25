package bot

import "context"

func resetGuessesCommand() *Command {
	return &Command{
		Name: "resetguesses", Aliases: []string{"rg"}, Access: Moderator,
		Usage:       "!resetguesses",
		Description: "Delete all open guesses.",
		Run:         runResetGuesses,
	}
}

func runResetGuesses(ctx context.Context, b *Bot, inv Invocation) (string, error) {
	if err := b.Guesses.Reset(ctx, inv.Channel.ID); err != nil {
		return "", err
	}
	return inv.mention() + " All guesses have been reset.", nil
}
