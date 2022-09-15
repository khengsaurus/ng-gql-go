import { Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';
import { TodosService } from 'src/app/services';
import { ITodo, Nullable } from 'src/types';

@Component({
  selector: 'todo-cards-page',
  templateUrl: './todos-page.component.html',
  styleUrls: ['./todos-page.component.scss'],
})
export class TodosPage implements OnInit, OnDestroy {
  sidenavOpen: boolean = true;
  selectedTodo: Nullable<ITodo> = null;
  private userTodosSub: Nullable<Subscription> = null;

  constructor(public todosService: TodosService) {}

  ngOnInit(): void {
    this.userTodosSub = this.todosService.currentUserTodos$.subscribe();
  }

  ngOnDestroy(): void {
    this.userTodosSub?.unsubscribe();
  }

  toggleSidenav() {
    this.sidenavOpen = !this.sidenavOpen;
  }

  selectTodo(todo: ITodo) {
    this.selectedTodo = todo;
  }

  trackById(_: number, item: Nullable<ITodo>) {
    return item?.id;
  }
}
