import { User } from '@/models/user';
import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private readonly TOKEN_KEY = 'token';
  private readonly USER_KEY = 'user';
  private readonly PERMISSIONS_KEY = 'permissions';

  constructor() { }

  setToken(token: string): void {
    localStorage.setItem(this.TOKEN_KEY, token);
  }

  getToken(): string | null {
    return localStorage.getItem(this.TOKEN_KEY);
  }

  setUser(user: User): void {
    localStorage.setItem(this.USER_KEY, JSON.stringify(user));
  }

  getUser(): User | null {
    const user = localStorage.getItem(this.USER_KEY);
    return user ? JSON.parse(user) as User : null;
  }

  setPermissions(permissions: string[]): void {
    localStorage.setItem(this.PERMISSIONS_KEY, JSON.stringify(permissions));
  }

  getPermissions(): string[] | null {
    const permissions = localStorage.getItem(this.PERMISSIONS_KEY);
    return permissions ? JSON.parse(permissions) as string[] : null;
  }

  decodeToken(token: string): any {
    return jwt_decode(token);
  }
}

function jwt_decode(token: string): any {
  try {
    // Split the JWT into header, payload, and signature
    const parts = token.split('.');
    if (parts.length !== 3) {
      throw new Error('The token is not a valid JWT');
    }

    // Decode the payload
    const decoded = atob(parts[1]);

    // Parse the JSON payload
    const payload = JSON.parse(decoded);

    return payload;
  } catch (error) {
    throw new Error('Error decoding the token: ' + error.message);
  }
}

