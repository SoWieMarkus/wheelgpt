import { HttpClient, provideHttpClient, withInterceptors } from "@angular/common/http";
import { type ApplicationConfig, importProvidersFrom, provideZoneChangeDetection } from "@angular/core";
import { provideRouter } from "@angular/router";
import { TranslateLoader, TranslateModule } from "@ngx-translate/core";
import { TranslateHttpLoader } from "@ngx-translate/http-loader";
import { routes } from "./app.routes";
import { authenticationInterceptor } from "./http/authentication.interceptor";
import { unauthorizedInterceptor } from "./http/unauthorized.interceptor";

export function HttpLoaderFactory(http: HttpClient) {
	return new TranslateHttpLoader(http, "./locale/", ".json");
}

export const appConfig: ApplicationConfig = {
	providers: [
		provideZoneChangeDetection({ eventCoalescing: true }),
		provideRouter(routes),
		provideHttpClient(withInterceptors([authenticationInterceptor, unauthorizedInterceptor])),
		importProvidersFrom(
			TranslateModule.forRoot({
				loader: {
					provide: TranslateLoader,
					useFactory: HttpLoaderFactory,
					deps: [HttpClient],
				},
			}),
		),
	],
};
