import type { HttpInterceptorFn } from "@angular/common/http";
import { inject } from "@angular/core";
import { AuthenticationService } from "../services/authentication.service";

export const authenticationInterceptor: HttpInterceptorFn = (request, next) => {
	const authenticationService = inject(AuthenticationService);

	const token = authenticationService.getToken();
	const modifiedRequest = token
		? request.clone({
				setHeaders: {
					Authorization: token,
				},
			})
		: request;

	return next(modifiedRequest);
};
