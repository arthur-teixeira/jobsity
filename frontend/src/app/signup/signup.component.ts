import { Component, inject } from '@angular/core';
import { take } from 'rxjs';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { TaskComponent } from '../task/task.component';
import { AuthService } from '../auth.service';
import { Router } from '@angular/router';

@Component({
  selector: 'app-signup',
  standalone: true,
  imports: [
    CommonModule,
    TaskComponent,
    ReactiveFormsModule
  ],
  templateUrl: './signup.component.html',
  styleUrl: './signup.component.css'
})
export class SignupComponent {
  private authService = inject(AuthService);
  private router = inject(Router);

  public signupForm = new FormGroup({
    email: new FormControl('', [Validators.required]),
    password: new FormControl('', [Validators.required]),
  });

  public signUp() {
    const email = this.signupForm.controls.email.value!;
    const password = this.signupForm.controls.password.value!;

    this.authService.signUp(email, password)
      .pipe(take(1))
      .subscribe();
  }

  public goToLogin() {
    this.router.navigate(['/login']);
  }
}
