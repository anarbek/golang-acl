import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { User } from '@/models/user';

@Injectable({
  providedIn: 'root'
})
export class UsersService {

  private apiUrl = 'http://localhost:8081/api/v1/users/';

  constructor(private http: HttpClient,) { }

  getItems(): Observable<any> {
    return this.http.get<any>(`${this.apiUrl}`)
  }

  getItem(id: number): Observable<User> {
    return this.http.get<User>(`${this.apiUrl}${id}`)
  }

  updateItem(id: number, item: any): Observable<any> {
    return this.http.put(`${this.apiUrl}update`, item);
  }

  insertItem(item: any): Observable<any> {
    return this.http.post(`${this.apiUrl}insert`, item);
  }
}
