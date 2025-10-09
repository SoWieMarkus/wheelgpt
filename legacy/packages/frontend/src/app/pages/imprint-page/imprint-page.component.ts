import { Component } from "@angular/core";
import { TranslatePipe } from "@ngx-translate/core";

@Component({
	selector: "app-imprint-page",
	imports: [TranslatePipe],
	templateUrl: "./imprint-page.component.html",
	styleUrl: "./imprint-page.component.scss",
})
export class ImprintPage {}
