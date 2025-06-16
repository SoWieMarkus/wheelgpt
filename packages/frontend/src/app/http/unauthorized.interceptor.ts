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
				console.warn("Unauthorized request, redirecting to login");
				authenticationService.removeToken();
				router.navigate(["/login"]);
			}
			return throwError(() => error);
		}),
	);
};
