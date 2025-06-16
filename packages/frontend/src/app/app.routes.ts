import type { Routes } from "@angular/router";
import { AuthenticationGuard } from "./authentication.guard";
import { DashboardPage } from "./pages/dashboard-page/dashboard-page.component";
import { LandingPage } from "./pages/landing-page/landing-page.component";
import { PrivacyPolicyPage } from "./pages/privacy-policy-page/privacy-policy-page.component";
import { TwitchRedirectPage } from "./pages/twitch-redirect-page/twitch-redirect-page.component";

export const routes: Routes = [
	{
		path: "twitch/callback",
		component: TwitchRedirectPage,
		title: "WheelGPT: Redirecting ...",
	},
	{
		path: "",
		component: LandingPage,
		title: "WheelGPT: Welcome!",
	},
	{
		path: "privacy",
		component: PrivacyPolicyPage,
		title: "WheelGPT: Privacy Policy",
	},
	{
		path: "",
		title: "WheelGPT: Dashboard",
		component: DashboardPage,
		canActivate: [AuthenticationGuard],
	},
	{
		path: "**",
		redirectTo: "/",
		pathMatch: "full",
	},
];
