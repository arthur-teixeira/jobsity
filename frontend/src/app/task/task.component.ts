import { Component, EventEmitter, Input, OnDestroy, OnInit, Output, inject } from '@angular/core';
import { Task } from '../types';
import { TaskService } from '../task.service';
import { CommonModule } from '@angular/common';
import { FormControl, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { faFloppyDisk, faPenToSquare } from '@fortawesome/free-regular-svg-icons';
import { faTrash } from '@fortawesome/free-solid-svg-icons';
import { Subscription, catchError, of } from 'rxjs';

@Component({
  selector: 'app-task',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    FontAwesomeModule,
  ],
  templateUrl: './task.component.html',
  styleUrl: './task.component.css'
})

export class TaskComponent implements OnInit, OnDestroy {
  @Input() task!: Task;
  @Output() deleteEvent: EventEmitter<number> = new EventEmitter();
  @Output() editEvent: EventEmitter<Task> = new EventEmitter();
  @Output() errorEvent: EventEmitter<void> = new EventEmitter();

  todoService = inject(TaskService);

  isEditing = false;
  subscriptions: Subscription[] = [];

  faPenToSquare = faPenToSquare;
  faTrash = faTrash;
  faSave = faFloppyDisk;

  taskForm = new FormGroup({
    id: new FormControl(0),
    title: new FormControl('', [Validators.required]),
    isCompleted: new FormControl(false),
  });

  ngOnInit() {
    this.taskForm.setValue({
      id: this.task.id,
      title: this.task.title,
      isCompleted: this.task.isCompleted,
    });
  }

  ngOnDestroy() {
    this.subscriptions.forEach(s => s.unsubscribe());
  }

  onDelete() {
    const subscription = this.todoService.deleteTask(this.task.id)
      .pipe(catchError(() => {
        this.errorEvent.emit();
        return of(true);
      }))
      .subscribe((hasError) => {
        if (!hasError) {
          this.deleteEvent.emit(this.task.id);
        }
      });

    this.subscriptions.push(subscription);
  }

  private save(task: Task) {
    const subscription = this.todoService.editTask(task)
      .pipe(catchError(() => {
        this.errorEvent.emit();
        return of(null);
      }))
      .subscribe(updatedTask => {
        if (!updatedTask) {
          return;
        }

        this.task = updatedTask;
        this.editEvent.emit(this.task);
      });

    this.subscriptions.push(subscription);
  }

  onComplete(event: any) {
    this.save({
      id: this.task.id,
      title: this.task.title,
      isCompleted: event.target.checked,
    });
  }

  onEdit() {
    this.save(this.taskForm.value as Task)
    this.isEditing = false;
  }
}
