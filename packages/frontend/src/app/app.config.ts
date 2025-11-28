import { provideHttpClient, withInterceptors } from "@angular/common/http";
import { type ApplicationConfig, provideZonelessChangeDetection } from "@angular/core";
import { provideRouter } from "@angular/router";
import { provideTranslateService } from "@ngx-translate/core";
import { provideTranslateHttpLoader } from "@ngx-translate/http-loader";
import { routes } from "./app.routes";
import { authenticationInterceptor } from "./http/authentication.interceptor";
import { unauthorizedInterceptor } from "./http/unauthorized.interceptor";

export const appConfig: ApplicationConfig = {
	providers: [
		provideZonelessChangeDetection(),
		provideRouter(routes),
		provideHttpClient(withInterceptors([authenticationInterceptor, unauthorizedInterceptor])),
		provideTranslateService({
			loader: provideTranslateHttpLoader({
				prefix: "./locale/",
				suffix: ".json",
			}),
		}),
	],
};
