import { Routes } from '@angular/router';
import { LoginComponent } from './login/login.component';
import { TaskListComponent } from './task-list/task-list.component';
import { authenticationGuard } from './auth-guard';
import { SignupComponent } from './signup/signup.component';

export const routes: Routes = [
  { path: 'login', component: LoginComponent },
  { path: 'signup', component: SignupComponent },
  {
    path: '',
    component: TaskListComponent,
    canActivate: [authenticationGuard()],
  },
  { path: '**', redirectTo: '' }
];
