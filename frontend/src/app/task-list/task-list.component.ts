import { Component, inject } from '@angular/core';
import { TaskService } from '../task.service';
import { Task, TaskRequest } from '../types';
import { Subscription, catchError, of } from 'rxjs';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { TaskComponent } from '../task/task.component';
import { Router } from '@angular/router';

const filters = [
  "All",
  "Pending",
  "Completed",
] as const;
type Filter = typeof filters[number];


@Component({
  selector: 'app-task-list',
  standalone: true,
  imports: [
    CommonModule,
    TaskComponent,
    ReactiveFormsModule
  ],
  templateUrl: './task-list.component.html',
  styleUrl: './task-list.component.css'
})
export class TaskListComponent {
  private taskService = inject(TaskService);
  private router = inject(Router);

  private tasks: Task[] = [];

  private subscriptions: Subscription[] = [];

  filteredTasks: Task[] = [];

  filters = filters;

  hasError = false;

  currentFilter: Filter = 'All';

  taskForm = new FormGroup({
    title: new FormControl('', [Validators.required]),
  });

  ngOnInit() {
    const tasksSub = this.taskService
      .getTasks()
      .pipe(catchError(e => {
        if (e.status === 401) {
          this.router.navigate(['login']);
        }
        this.hasError = true;
        return of(null);
      }))
      .subscribe(todos => {
        if (!todos) {
          return;
        }

        this.tasks = todos;
        this.filteredTasks = this.tasks;
      });

    this.subscriptions.push(tasksSub);
  }

  ngOnDestroy() {
    this.subscriptions.forEach(s => s.unsubscribe());
  }

  createTask(): void {
    const payload = {
      title: this.taskForm.get('title')?.value || '',
    } as TaskRequest;

    const subscription = this.taskService
      .saveTask(payload)
      .pipe(catchError(() => {
        this.hasError = true;
        return of(null);
      }))
      .subscribe(newTask => {
        if (newTask) {
          this.tasks.push(newTask);
        }
        this.taskForm.reset({
          title: '',
        });
      });

    this.subscriptions.push(subscription);
  }

  onTaskDeleted(id: any) {
    console.log(id);
    this.tasks = this.tasks.filter(todo => todo.id !== id);
    this.setFilter(this.currentFilter);
  }

  onTaskEdited(newTask: any) {
    console.log(newTask);
    this.tasks = this.tasks.map(task => {
      if (task.id === newTask.id) {
        return newTask;
      }

      return task;
    });

    this.setFilter(this.currentFilter);
  }

  setFilter(filter: Filter) {
    this.currentFilter = filter;
    this.filteredTasks = this.tasks.filter(task => {
      if (this.currentFilter === 'All') {
        return true;
      }

      return this.currentFilter === 'Completed'
        ? task.isCompleted
        : !task.isCompleted;
    });
  }

  onError() {
    this.hasError = true;
  }
}
