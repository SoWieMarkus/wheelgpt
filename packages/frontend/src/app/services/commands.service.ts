import { Injectable, signal } from "@angular/core";

export type CommandExampleMessage = {
  isBot: boolean;
  text: string;
};

export type Command = {
  name: string;
  aliases: string[];
  description: string;
  example: CommandExampleMessage[];
  accessLevel: "User" | "Subscriber" | "VIP" | "Mod" | "Streamer";
};

@Injectable({
  providedIn: "root",
})
export class CommandsService {
  public readonly commands = signal<Command[]>([
    {
      name: "!map",
      aliases: [],
      description: "commands.map.description",
      example: [
        { isBot: false, text: "!map" },
        { isBot: true, text: '"MyCoolTrack" by MyCoolMapper | $$AUTHOR_MEDAL$$ 1:23.458' },
      ],
      accessLevel: "User",
    },
    {
      name: "!room",
      aliases: [],
      description: "commands.room.description",
      example: [
        { isBot: false, text: "!room" },
        { isBot: true, text: "Your Mamas Room [5/100]" },
      ],
      accessLevel: "User",
    },
    {
      name: "!guess",
      aliases: ["!g"],
      description: "commands.guess.description",
      example: [
        { isBot: false, text: "!guess 23.456" },
        { isBot: false, text: "!g 1:23.456" },
      ],
      accessLevel: "User",
    },
    {
      name: "!myguess",
      aliases: ["!mg"],
      description: "commands.myguess.description",
      example: [
        { isBot: false, text: "!myguess" },
        { isBot: true, text: "@sowiemarkus 1:23.456" },
        { isBot: false, text: "!mg" },
        { isBot: true, text: "@sowiemarkus 1:23.456" },
      ],
      accessLevel: "User",
    },
    {
      name: "!result",
      aliases: [],
      description: "commands.result.description",
      example: [
        { isBot: false, text: "!result 1:23.457" },
        {
          isBot: true,
          text: "NEW PERSONAL BEST 1:23.457 $$DINK_DONK$$ That's $$AUTHOR_MEDAL$$ $$OK$$ Nobody guessed it correctly but @sowiemarkus guessed 1:23.456 (- 0.001) $$OK$$",
        },
      ],
      accessLevel: "Mod",
    },
    {
      name: "!reset-guesses",
      aliases: ["!rg"],
      description: "commands.reset-guesses.description",
      example: [
        { isBot: false, text: "!reset-guesses" },
        { isBot: true, text: "@sowiemarkus All guesses have been reset." },
        { isBot: false, text: "!rg" },
        { isBot: true, text: "@sowiemarkus All guesses have been reset." },
      ],
      accessLevel: "Mod",
    },
    {
      name: "!wgpt-emotes",
      aliases: [],
      description: "commands.wgpt-emotes.description",
      example: [
        { isBot: false, text: "!wgpt-emotes" },
        { isBot: true, text: "$$YEK$$ $$BWOAH$$ $$GIGACHAD$$ ..." },
      ],
      accessLevel: "Mod",
    },
    {
      name: "!format",
      aliases: [],
      description: "commands.format.description",
      example: [
        { isBot: false, text: "!format viewer" },
        { isBot: true, text: "@viewer hh:mm:ss.xxx (hours and minutes are optional)" },
      ],
      accessLevel: "User",
    },
  ]);
}
