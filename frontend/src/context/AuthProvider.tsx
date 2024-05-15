import { AuthContextData, AuthContextProps, AuthContextType } from "@/types/Auth";
import { createContext, useState } from "react";

const AuthContext = createContext<AuthContextType>({} as AuthContextType);

export const AuthProvider = ({ children }: AuthContextProps) => {
    const [auth, setAuth] = useState<AuthContextData | null>(null);
    return <AuthContext.Provider value={{ auth, setAuth }}>{children}</AuthContext.Provider>;
};

export default AuthContext;
