import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, map } from 'rxjs';
import { environment } from '../environments/environment';
import { JwtHelperService } from '@auth0/angular-jwt';
import { AuthResponse } from './types';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private baseUrl = environment.baseUrl;

  constructor(
    private httpClient: HttpClient,
    private jwtHelper: JwtHelperService,
    private router: Router,
  ) { }

  isAuthenticated() {
    const token = localStorage.getItem('token');
    return !this.jwtHelper.isTokenExpired(token);
  }

  private saveSessionAndRedirect(res: AuthResponse) {
    localStorage.setItem('token', res.token);
    this.router.navigate(['']);
  }

  logIn(email: string, password: string): Observable<void> {
    return this.httpClient.post<AuthResponse>(`${this.baseUrl}/auth/login`, {
      email,
      password
    }).pipe(map(this.saveSessionAndRedirect.bind(this)));
  }

  signUp(email: string, password: string): Observable<void> {
    return this.httpClient.post<AuthResponse>(`${this.baseUrl}/auth/signup`, {
      email,
      password
    }).pipe(map(this.saveSessionAndRedirect.bind(this)));
  }
}
