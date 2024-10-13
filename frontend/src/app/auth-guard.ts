import { inject } from "@angular/core";
import { AuthService } from "./auth.service";
import { CanActivateFn, Router } from "@angular/router";

export function authenticationGuard(): CanActivateFn {
  return () => {
    const authService: AuthService = inject(AuthService);
    const router: Router = inject(Router);

    const isAuth = authService.isAuthenticated();
    if (!isAuth) {
      router.navigate(['/login'])
    }

    return isAuth;
  };
}
