import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, map } from 'rxjs';
import { AuthService } from './auth.service';

@Injectable({
  providedIn: 'root'
})
export class LoginService {

  private apiUrl = 'http://localhost:8081/loginUser';

  constructor(private http: HttpClient, private authService: AuthService) { }

  getUser(username: string, password: string): Observable<any> {
    return this.http.post<any>(`${this.apiUrl}`, {
      username: username,
      password: password
    }).pipe(
      map(data => {
        // The server should return a bearer token on successful login
        const token = data.token;
        this.authService.setToken(token);
  
        // Decode the token to get the user data
        const user = this.authService.decodeToken(token);
        this.authService.setUser(user);
  
        return data;
      })
    );
  }
}
