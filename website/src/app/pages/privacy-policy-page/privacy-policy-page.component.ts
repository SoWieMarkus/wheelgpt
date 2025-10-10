import { Component } from "@angular/core";
import { TranslatePipe } from "@ngx-translate/core";

@Component({
	selector: "app-privacy-policy-page",
	imports: [TranslatePipe],
	templateUrl: "./privacy-policy-page.component.html",
	styleUrl: "./privacy-policy-page.component.scss",
})
export class PrivacyPolicyPage {}
