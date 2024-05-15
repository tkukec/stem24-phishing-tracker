interface StateContextType<T> {
  auth: T;
  setAuth: React.Dispatch<React.SetStateAction<T>>;
}
export interface AuthContextData {
  user: Record<string, unknown>;
  accessToken: string;
  active: boolean;
}
export type AuthContextType = StateContextType<AuthContextData | null>;

export interface AuthContextProps {
  children: React.ReactNode;
}

export interface TokenResponse {
  token: string;
  active: boolean;
}

export interface RefreshTokenResponse extends TokenResponse {
  refreshToken: string;
}

export interface RegisterData {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  gender: string;
  age: number;
}

export interface LoginData {
  email: string;
  password: string;
}
