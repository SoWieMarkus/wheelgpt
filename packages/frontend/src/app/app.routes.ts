import type { Routes } from "@angular/router";
import { AuthenticationGuard } from "./authentication.guard";
import { DashboardPage } from "./pages/dashboard-page/dashboard-page.component";
import { ImprintPage } from "./pages/imprint-page/imprint-page.component";
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
		path: "landing",
		component: LandingPage,
		title: "WheelGPT: Welcome!",
	},
	{
		path: "privacy-policy",
		component: PrivacyPolicyPage,
		title: "WheelGPT: Privacy Policy",
	},
	{
		path: "imprint",
		component: ImprintPage,
		title: "WheelGPT: Imprint",
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
