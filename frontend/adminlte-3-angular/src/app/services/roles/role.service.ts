import { Role } from '@/models/role';
import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class RoleService {
  private apiUrl = 'http://localhost:8081/api/v1/roles/';

  constructor(private http: HttpClient,) { }

  getRoles(): Observable<any> {
    return this.http.get<any>(`${this.apiUrl}`)
  }

  getRole(id: number): Observable<Role> {
    return this.http.get<Role>(`${this.apiUrl}${id}`)
  }

  updateRole(id: number, role: any): Observable<any> {
    return this.http.put(`${this.apiUrl}update`, role);
  }

  insertRole(role: any): Observable<any> {
    return this.http.post(`${this.apiUrl}insert`, role);
  }
}
