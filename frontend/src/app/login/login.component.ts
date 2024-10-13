import { Component, inject } from '@angular/core';
import { catchError, of, take } from 'rxjs';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { TaskComponent } from '../task/task.component';
import { AuthService } from '../auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [
    CommonModule,
    TaskComponent,
    ReactiveFormsModule
  ],
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})
export class LoginComponent {
  private authService = inject(AuthService);
  private router = inject(Router);

  hasError = false;
  errorMessage = "";

  public loginForm = new FormGroup({
    email: new FormControl('', [Validators.required]),
    password: new FormControl('', [Validators.required]),
  });

  public logIn() {
    const email = this.loginForm.controls.email.value!;
    const password = this.loginForm.controls.password.value!;

    this.authService.logIn(email, password)
      .pipe(take(1), catchError(e => {
        this.hasError = true;
        this.errorMessage = e.message || "An error has occurred";
        return of();
      })).subscribe();
  }

  goToSignup() {
    this.router.navigate(['signup']);
  }
}
