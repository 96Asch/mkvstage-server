export interface User {
  id: number;
  email: string;
  password: string;
}

export const emptyUser = {
  id: 0,
  email: "",
  password: "",
};
