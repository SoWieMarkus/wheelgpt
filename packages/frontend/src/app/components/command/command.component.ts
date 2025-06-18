import { Component, input } from "@angular/core";
import type { Command } from "../../services/commands.service";
import { TwitchMessagePipe } from "../../pipes/twitch-message.pipe";
import { TranslatePipe } from "@ngx-translate/core";

@Component({
  selector: "app-command",
  imports: [TwitchMessagePipe, TranslatePipe],
  templateUrl: "./command.component.html",
  styleUrl: "./command.component.scss",
})
export class CommandComponent {
  public readonly command = input.required<Command>();

  protected get showUserBadge() {
    return this.command().accessLevel !== "User";
  }

  protected get userBadge() {
    switch (this.command().accessLevel) {
      case "User":
        return "";
      case "Mod":
        return "badges/mod.png";
      case "Streamer":
        return "badges/streamer.png";
      case "VIP":
        return "badges/vip.png";
      case "Subscriber":
        return "badges/subscriber.png";
      default:
        return "";
    }
  }
}
