import { HttpErrorResponse, HttpEvent, HttpHandler, HttpInterceptor, HttpRequest } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { AppService } from '@services/app.service';
import { AuthService } from '@services/auth.service';
import { Observable, catchError, throwError } from 'rxjs';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {

  constructor(private appService: AppService, private authService: AuthService) { }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    // Get the auth token from the service.
    const authToken = this.authService.getToken();

    let authReq: HttpRequest<any>;
    if (authToken) {
      // Clone the request and replace the original headers with
      // cloned headers, updated with the authorization.
      authReq = req.clone({
        headers: req.headers.set('Authorization', `Bearer ${authToken}`)
      });
    } else {
      // If authToken is null, clone the request without modifying the headers.
      authReq = req.clone();
    }

    // send cloned request with header to the next handler.
    return next.handle(authReq).pipe(
      catchError((error: HttpErrorResponse) => {
        // If the error status is 401 (Unauthorized), navigate to the login page
        if (error.status === 401 || error.status === 403) {
          this.appService.logout();
        }

        return throwError(() => error);
      })
    );

    // send cloned request with header to the next handler.
    return next.handle(authReq);
  }
}
