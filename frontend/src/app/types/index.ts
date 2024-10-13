export type Task = {
  id: number,
  title: string,
  isCompleted: boolean
};

export type TaskRequest = {
  title: string
};

export type AuthResponse = {
  token: string
}
