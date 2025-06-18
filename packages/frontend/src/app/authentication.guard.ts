import { Injectable, inject } from "@angular/core";
import { type ActivatedRouteSnapshot, type CanActivateFn, Router, type RouterStateSnapshot } from "@angular/router";
import { AuthenticationService } from "./services/authentication.service";

@Injectable({
	providedIn: "root",
})
class AuthenticationPermissionService {
	private readonly router = inject(Router);
	private readonly authenticationService = inject(AuthenticationService);

	public canActivate(_: ActivatedRouteSnapshot, __: RouterStateSnapshot): boolean {
		const isTokenExpired = this.authenticationService.isTokenExpired();
		const doesTokenExist = this.authenticationService.getToken() !== null;

		if (isTokenExpired && !doesTokenExist) {
			this.router.navigate(["/landing"]);
			return false;
		}

		if (isTokenExpired && doesTokenExist) {
			this.router.navigate(["/login"]);
			return false;
		}

		return true;
	}
}

export const AuthenticationGuard: CanActivateFn = (
	next: ActivatedRouteSnapshot,
	state: RouterStateSnapshot,
): boolean => {
	const service = inject(AuthenticationPermissionService);
	return service.canActivate(next, state);
};
