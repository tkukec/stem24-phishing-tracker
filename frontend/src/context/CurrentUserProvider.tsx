import { UsersService } from "@/api/users";
import { createContext } from "react";
import type { User } from "@/types/User";
export const CurrentUserContext = createContext<User | null>(null);

const CurrentUserContextProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const currentUserQuery = UsersService.useGetCurrentUser();
  if (currentUserQuery.isLoading) {
    return null;
  }

  if (currentUserQuery.isError) {
    return <pre>{JSON.stringify(currentUserQuery.error)}</pre>;
  }

  return (
    <CurrentUserContext.Provider value={currentUserQuery.data ?? null}>
      {children}
    </CurrentUserContext.Provider>
  );
};

export default CurrentUserContextProvider;
