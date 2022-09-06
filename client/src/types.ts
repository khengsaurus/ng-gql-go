import { Route } from './enums';

export type Nullable<T> = T | null;

export interface IResponse<T> {
  data: {
    [key: string]: T;
  };
}

export interface IUser {
  id: string;
  username: string;
  email?: string;
}

export interface ITodo {
  id: string;
  userId: string;
  text: string;
  priority: number;
  tag: string;
  done: boolean;
}

export interface ILink {
  route: Route | string;
  label: string;
}

export interface ITodos {
  userId: string;
  todos: ITodo[];
}