import { Injectable } from '@angular/core';
import { Task, TaskRequest } from './types';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class TaskService {
  private baseUrl = environment.baseUrl;

  constructor(
    private httpClient: HttpClient,
  ) { }

  private headers() {
    console.log("fetching token", localStorage.getItem('token'))
    return {
      'Content-Type': 'application/json',
      'Authorization': localStorage.getItem('token') || '',
    }
  }

  public getTasks(): Observable<Task[]> {
    return this.httpClient.get<Task[]>(`${this.baseUrl}/tasks`, {
      headers: this.headers(),
    });
  }

  public saveTask(todo: TaskRequest): Observable<Task> {
    return this.httpClient.post<Task>(`${this.baseUrl}/task`, JSON.stringify(todo), {
      headers: this.headers(),
    });
  }

  public deleteTask(id: number): Observable<void> {
    return this.httpClient.delete<void>(`${this.baseUrl}/task?id=${id}`, {
      headers: this.headers()
    });
  }

  public editTask(task: Task): Observable<Task> {
    return this.httpClient.put<Task>(`${this.baseUrl}/task`, JSON.stringify(task), {
      headers: this.headers(),
    });
  }
}
