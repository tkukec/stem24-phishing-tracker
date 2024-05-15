import useAuth from "@/hooks/useAuth";
import useAxiosPrivate from "@/hooks/useAxiosPrivate";
import { User } from "@/types/User";
import { useQuery } from "@tanstack/react-query";

export const UsersService = {
  useGetCurrentUser() {
    const { auth } = useAuth();
    const axios = useAxiosPrivate();
    return useQuery<User, Error>({
      queryKey: ["UsersService.getCurrentUser", auth?.accessToken],
      queryFn: () => axios.get("/users/me").then((response) => response.data),
      staleTime: Infinity,
    });
  },
}
