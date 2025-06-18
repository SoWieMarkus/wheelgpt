import type { HttpInterceptorFn } from "@angular/common/http";
import { inject } from "@angular/core";
import { Router } from "@angular/router";
import { catchError, throwError } from "rxjs";
import { AuthenticationService } from "../services/authentication.service";

export const unauthorizedInterceptor: HttpInterceptorFn = (request, next) => {
	const authenticationService = inject(AuthenticationService);
	const router = inject(Router);

	return next(request).pipe(
		catchError((error) => {
			if (error.status === 401) {
				authenticationService.removeToken();
				router.navigate(["/landing"]);
			}
			return throwError(() => error);
		}),
	);
};
