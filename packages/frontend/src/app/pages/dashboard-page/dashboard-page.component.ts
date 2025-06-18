import { Component, inject } from "@angular/core";
import { ProfileService } from "../../services/profile.service";

@Component({
	selector: "app-dashboard-page",
	imports: [],
	templateUrl: "./dashboard-page.component.html",
	styleUrl: "./dashboard-page.component.scss",
})
export class DashboardPage {
	protected readonly profileService = inject(ProfileService);
}
