import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TaskListComponent } from './task-list/task-list.component';
import { TaskComponent } from './task/task.component';
import { ReactiveFormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { AppComponent } from './app.component';

@NgModule({
  declarations: [],
  imports: [
    ReactiveFormsModule,
    HttpClientModule,
    TaskListComponent,
    TaskListComponent,
    CommonModule,
    AppComponent,
    TaskComponent
  ],
  exports: [
    AppComponent,
  ]
})
export class AppModule { }
